/*
Copyright 2021 Flant JSC

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package probe

import (
	"math/rand"
	"sort"
	"time"

	"github.com/sirupsen/logrus"

	"d8.io/upmeter/pkg/check"
	"d8.io/upmeter/pkg/kubernetes"
	"d8.io/upmeter/pkg/monitor/node"
	"d8.io/upmeter/pkg/probe/checker"
	"d8.io/upmeter/pkg/set"
)

type Loader interface {
	Load() []*check.Runner
	Groups() []string
	Probes() []check.ProbeRef
}

func NewFakeLoader(logger *logrus.Logger) *FakeLoader {
	return &FakeLoader{logger: logger}
}

type FakeLoader struct {
	logger *logrus.Logger
}

func (l *FakeLoader) Load() []*check.Runner {
	runners := make([]*check.Runner, 0, len(l.Probes()))

	for _, ref := range l.Probes() {
		// random period
		// 1. 5s is pretty common
		// 2. 1m is 2 times slower than episode export
		// 3. 200ms is fast
		period := randomPeriod(5*time.Second, 1*time.Minute, 200*time.Millisecond)
		logger := l.logger.WithField("group", ref.Group).WithField("probe", ref.Probe).WithField("check", "fake")

		r := check.NewRunner(ref.Group, ref.Probe, "1", period, checker.Fake(), logger)

		runners = append(runners, r)
	}
	return runners
}

func (l *FakeLoader) Groups() []string {
	groups := set.New()
	for _, group := range l.Probes() {
		groups.Add(group.Group)
	}
	return groups.Slice()
}

func (l *FakeLoader) Probes() []check.ProbeRef {
	return []check.ProbeRef{
		{Group: "control-plane", Probe: "apiserver"},
		{Group: "control-plane", Probe: "basic-functionality"},
		{Group: "control-plane", Probe: "cert-manager"},
		{Group: "control-plane", Probe: "controller-manager"},
		{Group: "control-plane", Probe: "namespace"},
		{Group: "control-plane", Probe: "scheduler"},
		{Group: "deckhouse", Probe: "cluster-configuration"},
		{Group: "extensions", Probe: "cluster-autoscaler"},
		{Group: "extensions", Probe: "cluster-scaling"},
		{Group: "extensions", Probe: "dashboard"},
		{Group: "extensions", Probe: "dex"},
		{Group: "extensions", Probe: "grafana"},
		{Group: "extensions", Probe: "openvpn"},
		{Group: "extensions", Probe: "prometheus-longterm"},
		{Group: "load-balancing", Probe: "load-balancer-configuration"},
		{Group: "load-balancing", Probe: "metallb"},
		// {Group: "monitoring-and-autoscaling", Probe: "horizontal-pod-autoscaler"},  // calculated
		{Group: "monitoring-and-autoscaling", Probe: "key-metrics-present"},
		{Group: "monitoring-and-autoscaling", Probe: "metrics-sources"},
		{Group: "monitoring-and-autoscaling", Probe: "prometheus"},
		{Group: "monitoring-and-autoscaling", Probe: "prometheus-metrics-adapter"},
		{Group: "monitoring-and-autoscaling", Probe: "trickster"},
		{Group: "monitoring-and-autoscaling", Probe: "vertical-pod-autoscaler"},
		{Group: "synthetic", Probe: "access"},
		{Group: "synthetic", Probe: "dns"},
		{Group: "synthetic", Probe: "neighbor"},
		{Group: "synthetic", Probe: "neighbor-via-service"},
	}

}

func NewLoader(
	filter Filter,
	access kubernetes.Access,
	nodeLister node.Lister,
	dynamic DynamicConfig,
	preflight checker.Doer,
	logger *logrus.Logger,
) *ProbeLoader {
	return &ProbeLoader{
		filter:     filter,
		access:     access,
		dynamic:    dynamic,
		nodeLister: nodeLister,
		preflight:  preflight,
		logger:     logger,
	}
}

type ProbeLoader struct {
	filter     Filter
	access     kubernetes.Access
	logger     *logrus.Logger
	dynamic    DynamicConfig
	nodeLister node.Lister
	preflight  checker.Doer

	// inner state

	groups []string
	probes []check.ProbeRef

	configs []runnerConfig
}

type DynamicConfig struct {
	IngressNginxControllers []string
	NodeGroups              []string
	Zones                   []string
	ZonePrefix              string
}

func (l *ProbeLoader) Load() []*check.Runner {
	runners := make([]*check.Runner, 0)
	for _, rc := range l.collectConfigs() {
		if !l.filter.Enabled(rc.Ref()) {
			continue
		}

		runnerLogger := l.logger.WithFields(map[string]interface{}{
			"group": rc.group,
			"probe": rc.probe,
			"check": rc.check,
		})

		runner := check.NewRunner(rc.group, rc.probe, rc.check, rc.period, rc.config.Checker(), runnerLogger)

		runners = append(runners, runner)
		l.logger.Infof("Register probe %s", runner.ProbeRef().Id())
	}

	return runners
}

func (l *ProbeLoader) Groups() []string {
	if l.groups != nil {
		return l.groups
	}

	groups := set.New()
	for _, rc := range l.collectConfigs() {
		if !l.filter.Enabled(rc.Ref()) {
			continue
		}
		groups.Add(rc.group)

	}

	l.groups = groups.Slice()
	return l.groups
}

func (l *ProbeLoader) Probes() []check.ProbeRef {
	if l.probes != nil {
		return l.probes
	}

	seen := set.New()
	l.probes = make([]check.ProbeRef, 0)
	for _, rc := range l.collectConfigs() {
		ref := rc.Ref()
		if !l.filter.Enabled(ref) {
			continue
		}
		if seen.Has(ref.Id()) {
			continue
		}
		seen.Add(ref.Id())
		l.probes = append(l.probes, ref)

	}
	sort.Sort(check.ByProbeRef(l.probes))
	return l.probes
}

func (l *ProbeLoader) collectConfigs() []runnerConfig {
	if l.configs != nil {
		// Already inited
		return l.configs
	}

	l.configs = make([]runnerConfig, 0)
	l.configs = append(l.configs, initSynthetic(l.access, l.logger)...)
	l.configs = append(l.configs, initControlPlane(l.access, l.preflight)...)
	l.configs = append(l.configs, initMonitoringAndAutoscaling(l.access, l.nodeLister, l.preflight)...)
	l.configs = append(l.configs, initExtensions(l.access, l.preflight)...)
	l.configs = append(l.configs, initLoadBalancing(l.access, l.preflight)...)
	l.configs = append(l.configs, initDeckhouse(l.access, l.preflight, l.logger)...)
	l.configs = append(l.configs, initNginx(l.access, l.preflight, l.dynamic.IngressNginxControllers)...)
	l.configs = append(l.configs, initNodeGroups(l.access, l.nodeLister, l.preflight, l.dynamic.NodeGroups, l.dynamic.Zones, l.dynamic.ZonePrefix)...)

	return l.configs
}

type runnerConfig struct {
	group  string
	probe  string
	check  string
	period time.Duration
	config checker.Config
}

func (rc runnerConfig) Ref() check.ProbeRef {
	return check.ProbeRef{Group: rc.group, Probe: rc.probe}
}

func NewProbeFilter(disabled []string) Filter {
	return Filter{refs: set.New(disabled...)}
}

type Filter struct {
	refs set.StringSet
}

func (f Filter) Enabled(ref check.ProbeRef) bool {
	return !(f.refs.Has(ref.Id()) || f.refs.Has(ref.Group) || f.refs.Has(ref.Group+"/"))
}

// randomPeriod picking from edge cases
func randomPeriod(periods ...time.Duration) time.Duration {
	return periods[rand.Intn(len(periods))]
}
