/*
Copyright 2024 Flant JSC
Licensed under the Deckhouse Platform Enterprise Edition (EE) license. See https://github.com/deckhouse/deckhouse/blob/main/ee/LICENSE
*/

package registry_controller

import (
	"context"
	"fmt"
	"regexp"

	"embeded-registry-manager/internal/state"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var (
	nodePKISecretRegex        = regexp.MustCompile(`^registry-node-(.*)-pki$`)
	masterNodesMatchingLabels = client.MatchingLabels{
		state.LabelNodeIsMasterKey: "",
	}
)

type NodeController = nodeController

var _ reconcile.Reconciler = &nodeController{}

type nodeController struct {
	Client    client.Client
	Namespace string

	reprocessCh chan event.TypedGenericEvent[reconcile.Request]
}

var nodeReprocessAllRequest = reconcile.Request{
	NamespacedName: types.NamespacedName{
		Namespace: "--reprocess-all-nodes--",
		Name:      "--reprocess-all-nodes--",
	},
}

func (r *nodeController) SetupWithManager(ctx context.Context, mgr ctrl.Manager) error {
	r.reprocessCh = make(chan event.TypedGenericEvent[reconcile.Request])

	controllerName := "node-controller"

	// Set up the field indexer to index Pods by spec.nodeName
	err := mgr.GetFieldIndexer().
		IndexField(ctx, &corev1.Pod{}, "spec.nodeName", func(obj client.Object) []string {
			pod := obj.(*corev1.Pod)
			return []string{pod.Spec.NodeName}
		})

	if err != nil {
		return fmt.Errorf("failed to set up index on spec.nodeName: %w", err)
	}

	nodePredicate := predicate.Funcs{
		CreateFunc: func(e event.TypedCreateEvent[client.Object]) bool {
			// Only process master nodes
			return nodeObjectIsMaster(e.Object)
		},
		DeleteFunc: func(e event.TypedDeleteEvent[client.Object]) bool {
			// Only process master nodes
			return nodeObjectIsMaster(e.Object)
		},
		UpdateFunc: func(e event.TypedUpdateEvent[client.Object]) bool {
			// Only on master status change
			return nodeObjectIsMaster(e.ObjectOld) != nodeObjectIsMaster(e.ObjectNew)
		},
	}

	secretsPredicate := predicate.NewPredicateFuncs(secretObjectIsNodePKI)

	err = ctrl.NewControllerManagedBy(mgr).
		Named(controllerName).
		For(
			&corev1.Node{},
			builder.WithPredicates(nodePredicate),
		).
		Watches(
			&corev1.Secret{},
			handler.EnqueueRequestsFromMapFunc(func(_ context.Context, obj client.Object) []reconcile.Request {
				name := obj.GetName()
				sub := nodePKISecretRegex.FindStringSubmatch(name)

				if len(sub) < 2 {
					return nil
				}

				var ret reconcile.Request
				ret.Name = sub[1]

				log := ctrl.LoggerFrom(ctx)

				log.Info(
					"Node PKI secret was changed, will trigger reconcile for node",
					"secret", obj.GetName(),
					"namespace", obj.GetNamespace(),
					"node", ret.Name,
					"controller", controllerName,
				)

				return []reconcile.Request{ret}
			}),
			builder.WithPredicates(secretsPredicate),
		).
		WatchesRawSource(r.reprocessChannelSource()).
		WithOptions(controller.Options{
			MaxConcurrentReconciles: 10,
		}).
		Complete(r)

	if err != nil {
		return fmt.Errorf("cannot build controller: %w", err)
	}

	return nil
}

func (nc *nodeController) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := ctrl.LoggerFrom(ctx)

	if req == nodeReprocessAllRequest {
		return nc.handleReprocessAll(ctx)
	}

	if req.Namespace != "" {
		log.Info("Fired by supplementary object", "namespace", req.Namespace)
		req.Namespace = ""
	}

	node := &corev1.Node{}
	err := nc.Client.Get(ctx, req.NamespacedName, node)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return nc.handleNodeDelete(ctx, req.Name)
		}

		return ctrl.Result{}, fmt.Errorf("cannot get node: %w", err)
	}

	if hasMasterLabel(node) {
		return nc.handleMasterNode(ctx, node)
	} else {
		return nc.handleNodeNotMaster(ctx, node)
	}
}

func (nc *nodeController) handleReprocessAll(ctx context.Context) (ctrl.Result, error) {
	log := ctrl.LoggerFrom(ctx)

	log.Info("All nodes will be reprocessed")

	// Will trigger reprocess for all master nodes
	nodes := &corev1.NodeList{}
	if err := nc.Client.List(ctx, nodes, &masterNodesMatchingLabels); err != nil {
		return ctrl.Result{}, fmt.Errorf("cannot list nodes: %w", err)
	}

	for _, node := range nodes.Items {
		req := reconcile.Request{}
		req.Name = node.Name
		req.Namespace = node.Namespace

		if err := nc.triggerReconcileForNode(ctx, req); err != nil {
			// It currently only when ctx done
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

func (nc *nodeController) triggerReconcileForNode(ctx context.Context, req reconcile.Request) error {
	evt := event.TypedGenericEvent[reconcile.Request]{Object: req}

	select {
	case nc.reprocessCh <- evt:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (nc *nodeController) ReprocessAllNodes(ctx context.Context) error {
	return nc.triggerReconcileForNode(ctx, nodeReprocessAllRequest)
}

func (nc *nodeController) handleMasterNode(ctx context.Context, node *corev1.Node) (ctrl.Result, error) {
	log := ctrl.LoggerFrom(ctx)

	log.Info("Processing master node", "node", node.Name)

	return ctrl.Result{}, nil
}

func (nc *nodeController) isFirstMasterNode(ctx context.Context, node *corev1.Node) (bool, error) {
	// Will trigger reprocess for all master nodes
	nodes := &corev1.NodeList{}
	if err := nc.Client.List(ctx, nodes, &masterNodesMatchingLabels); err != nil {
		return false, fmt.Errorf("cannot list nodes: %w", err)
	}

	for _, item := range nodes.Items {
		if item.Name == node.Name {
			continue
		}

		if item.CreationTimestamp.Before(&node.CreationTimestamp) {
			return false, nil
		}
	}

	return true, nil
}

func (nc *nodeController) handleNodeNotMaster(ctx context.Context, node *corev1.Node) (ctrl.Result, error) {
	// Delete node secret if exists
	if err := nc.deleteNodePKI(ctx, node.Name); err != nil {
		return ctrl.Result{}, err
	}

	/*
		Here we may stop static pods, but it will be a race with k8s scheduler
		So NOOP
	*/

	return ctrl.Result{}, nil
}

func (nc *nodeController) handleNodeDelete(ctx context.Context, name string) (ctrl.Result, error) {
	// Delete node secret if exists
	if err := nc.deleteNodePKI(ctx, name); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (nc *nodeController) deleteNodePKI(ctx context.Context, nodeName string) error {
	log := ctrl.LoggerFrom(ctx)

	secretName := fmt.Sprintf("registry-node-%s-pki", nodeName)
	secret := &corev1.Secret{}

	err := nc.Client.Get(ctx, types.NamespacedName{Name: secretName, Namespace: nc.Namespace}, secret)

	if err != nil {
		if apierrors.IsNotFound(err) {
			// Already absent
			return nil
		}

		return fmt.Errorf("get node PKI secret error: %w", err)
	}

	err = nc.Client.Delete(ctx, secret)
	if client.IgnoreNotFound(err) != nil {
		return fmt.Errorf("delete node PKI secret error: %w", err)
	} else {
		log.Info("Deleted node PKI", "node", nodeName, "name", secret.Name, "namespace", secret.Namespace)
	}

	return nil
}

func (nc *nodeController) reprocessChannelSource() source.Source {
	return source.Channel(nc.reprocessCh, handler.TypedEnqueueRequestsFromMapFunc(
		func(_ context.Context, req reconcile.Request) []reconcile.Request {
			return []reconcile.Request{req}
		},
	))
}

func secretObjectIsNodePKI(obj client.Object) bool {
	labels := obj.GetLabels()

	if labels[state.LabelTypeKey] != state.LabelNodeSecretTypeValue {
		return false
	}

	if labels[state.LabelHeritageKey] != state.LabelHeritageValue {
		return false
	}

	return nodePKISecretRegex.MatchString(obj.GetName())
}

func nodeObjectIsMaster(obj client.Object) bool {
	if obj == nil {
		return false
	}

	labels := obj.GetLabels()
	if labels == nil {
		return false
	}

	_, hasMasterLabel := labels["node-role.kubernetes.io/master"]

	return hasMasterLabel
}

func hasMasterLabel(node *corev1.Node) bool {
	_, isMaster := node.Labels["node-role.kubernetes.io/master"]
	return isMaster
}
