// Copyright 2024 Flant JSC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package hooks

import (
	"math"
	"strconv"

	"github.com/flant/addon-operator/pkg/module_manager/go_hook"
	"github.com/flant/addon-operator/sdk"
	"github.com/flant/shell-operator/pkg/kube/object_patch"
	"github.com/flant/shell-operator/pkg/kube_events_manager/types"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

const (
	masterNodeRole         = "node-role.kubernetes.io/master"
	virtualizationLevelKey = "node.deckouse.io/dvp-nesting-level"
)

var _ = sdk.RegisterFunc(&go_hook.HookConfig{
	Kubernetes: []go_hook.KubernetesConfig{
		{
			Name:       "virtualization_level_secret",
			ApiVersion: "v1",
			Kind:       "Secret",
			NameSelector: &types.NameSelector{
				MatchNames: []string{"d8-virtualization-level"},
			},
			NamespaceSelector: &types.NamespaceSelector{
				NameSelector: &types.NameSelector{
					MatchNames: []string{"d8-system"},
				},
			},
			FilterFunc: applyVirtualizationLevelFilter,
		},
		{
			Name:       "master_nodes",
			ApiVersion: "v1",
			Kind:       "Node",
			LabelSelector: &metav1.LabelSelector{MatchExpressions: []metav1.LabelSelectorRequirement{{
				Key:      masterNodeRole,
				Operator: metav1.LabelSelectorOpExists,
			}}},
			FilterFunc: applyMasterNodesFilter,
		},
	},
}, setGlobalVirtualizationLevel)

func applyVirtualizationLevelFilter(obj *unstructured.Unstructured) (go_hook.FilterResult, error) {
	return obj.GetName(), nil
}

func applyMasterNodesFilter(obj *unstructured.Unstructured) (go_hook.FilterResult, error) {
	var node corev1.Node

	err := sdk.FromUnstructured(obj, &node)
	if err != nil {
		return masterNodeInfo{}, err
	}

	virtualizationLevel := 0
	if value, exists := node.GetLabels()[virtualizationLevelKey]; exists {
		virtualizationLevel, _ = strconv.Atoi(value)
	}
	return masterNodeInfo{name: node.GetName(), virtualizationLevel: virtualizationLevel}, nil
}

func setGlobalVirtualizationLevel(input *go_hook.HookInput) error {
	virtLevelSecretSnap := input.Snapshots["virtualization_level_secret"]

	if len(virtLevelSecretSnap) == 0 {
		input.LogEntry.Info("secret d8-virtualization-level not found, will be created automatically")
		minimalVirtualizationLevel := math.MaxInt
		for _, masterNodeInfoSnap := range input.Snapshots["master_nodes"] {
			masterNodeInfo := masterNodeInfoSnap.(masterNodeInfo)
			if masterNodeInfo.virtualizationLevel >= 0 && masterNodeInfo.virtualizationLevel < minimalVirtualizationLevel {
				minimalVirtualizationLevel = masterNodeInfo.virtualizationLevel
			}
		}
		if minimalVirtualizationLevel == math.MaxInt {
			minimalVirtualizationLevel = 0
		}
		secret := &corev1.Secret{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Secret",
				APIVersion: "v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "d8-virtualization-level",
				Namespace: "d8-system",
				Labels: map[string]string{
					"app":      "deckhouse",
					"module":   "deckhouse",
					"heritage": "deckhouse",
				},
			},
			Data: map[string][]byte{"level": []byte(strconv.Itoa(minimalVirtualizationLevel))},
		}

		input.PatchCollector.Create(secret, object_patch.UpdateIfExists())

		input.Values.Set("global.discovery.dvpNestingLevel", minimalVirtualizationLevel)
		input.LogEntry.Infof("set DVP nesting level to: %d", minimalVirtualizationLevel)
	}

	return nil
}

type masterNodeInfo struct {
	name                string
	virtualizationLevel int
}
