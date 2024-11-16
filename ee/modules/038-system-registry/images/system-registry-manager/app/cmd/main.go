/*
Copyright 2024 Flant JSC
Licensed under the Deckhouse Platform Enterprise Edition (EE) license. See https://github.com/deckhouse/deckhouse/blob/main/ee/LICENSE
*/

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"

	"embeded-registry-manager/internal/controllers/registry_controller"
	staticpodmanager "embeded-registry-manager/internal/static-pod"
	httpclient "embeded-registry-manager/internal/utils/http_client"
)

type managerStatus struct {
	isLeader                bool
	staticPodManagerRunning bool
}

const (
	metricsBindAddressPort = "127.0.0.1:8081"
	healthPort             = ":8097"
)

func main() {
	ctrl.SetLogger(zap.New(zap.UseDevMode(true)))
	status := &managerStatus{}
	ctrl.Log.Info("Starting embedded registry manager", "component", "main")

	// Load Kubernetes configuration
	cfg, err := loadKubeConfig()
	if err != nil {
		ctrl.Log.Error(err, "Unable to get kubeconfig", "component", "main")
		os.Exit(1)
	}

	// Create Kubernetes client
	kubeClient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		ctrl.Log.Error(err, "Unable to create Kubernetes client", "component", "main")
		os.Exit(1)
	}

	// Create custom HTTP client
	HttpClient, err := httpclient.NewDefaultHttpClient()
	if err != nil {
		ctrl.Log.Error(err, "Unable to create HTTP client", "component", "main")
		os.Exit(1)
	}

	ctx := ctrl.SetupSignalHandler()
	context.AfterFunc(ctx, func() {
		ctrl.Log.Info("Received shutdown signal")
	})

	// Start static pod manager
	go startStaticPodManager(ctx, status)

	// Start health server for readiness and liveness probes
	go startHealthServer(status)

	// Set up and start manager
	if err := setupAndStartManager(ctx, cfg, kubeClient, HttpClient, status); err != nil {
		ctrl.Log.Error(err, "Failed to start the embedded registry manager", "component", "main")
		os.Exit(1) // #TODO
	}
}

// setupAndStartManager sets up the manager, adds components, and starts the manager
func setupAndStartManager(ctx context.Context, cfg *rest.Config, kubeClient *kubernetes.Clientset, httpClient *httpclient.Client, status *managerStatus) error {
	// Set up the manager with leader election and other options
	mgr, err := ctrl.NewManager(cfg, manager.Options{
		Metrics: metricsserver.Options{
			BindAddress: metricsBindAddressPort,
		},
		LeaderElection:          true,
		LeaderElectionID:        "embedded-registry-manager-leader",
		LeaderElectionNamespace: "d8-system",
		GracefulShutdownTimeout: &[]time.Duration{10 * time.Second}[0],
		Cache: cache.Options{
			DefaultNamespaces: map[string]cache.Config{
				"d8-system": {},
			},
		},
	})
	if err != nil {
		return fmt.Errorf("unable to set up manager: %w", err)
	}

	// Create a new instance of RegistryReconciler
	reconciler := registry_controller.RegistryReconciler{
		Client:     mgr.GetClient(),
		APIReader:  mgr.GetAPIReader(),
		Scheme:     mgr.GetScheme(),
		KubeClient: kubeClient,
		Recorder:   mgr.GetEventRecorderFor("embedded-registry-controller"),
		HttpClient: httpClient,
	}

	// Set up the controller with the manager
	if err := reconciler.SetupWithManager(mgr, ctx); err != nil {
		return fmt.Errorf("unable to create controller: %w", err)
	}

	// Add leader status update runnable
	err = mgr.Add(leaderRunnableFunc(func(ctx context.Context) error {
		status.isLeader = true
		defer func() {
			status.isLeader = false
		}()

		// Call SecretsStartupCheckCreate with the existing reconciler instance
		if err := reconciler.SecretsStartupCheckCreate(ctx); err != nil {
			return fmt.Errorf("failed to initialize secrets: %w", err)
		}

		// Wait until the context is done to handle graceful shutdown
		<-ctx.Done()
		return nil
	}))

	if err != nil {
		return fmt.Errorf("unable to add leader runnable: %w", err)
	}

	// Start the manager
	ctrl.Log.Info("Starting manager", "component", "main")

	if err := mgr.Start(ctx); err != nil {
		return fmt.Errorf("unable to start manager: %w", err)
	}

	return nil
}

// loadKubeConfig tries to load the in-cluster config. If not available, it loads kubeconfig from home directory
func loadKubeConfig() (*rest.Config, error) {
	// Try to load in-cluster configuration
	cfg, err := rest.InClusterConfig()
	return cfg, err
}

// startStaticPodManager starts the static pod manager and monitors its status
func startStaticPodManager(ctx context.Context, status *managerStatus) {
	status.staticPodManagerRunning = true
	defer func() {
		status.staticPodManagerRunning = false
	}()

	if err := staticpodmanager.Run(ctx); err != nil {
		ctrl.Log.Error(err, "Failed to run static pod manager", "component", "main")
		os.Exit(1) // #TODO
	}
}

// startHealthServer starts a health server that provides readiness and liveness probes
func startHealthServer(status *managerStatus) {
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		healthHandler(w, status)
	})
	http.HandleFunc("/readyz", func(w http.ResponseWriter, r *http.Request) {
		healthHandler(w, status)
	})

	if err := http.ListenAndServe(healthPort, nil); err != nil {
		ctrl.Log.Error(err, "Failed to start health server", "component", "main")
		os.Exit(1)
	}
}

// healthHandler handles health and readiness probe requests
func healthHandler(w http.ResponseWriter, status *managerStatus) {
	response := struct {
		IsLeader                bool `json:"isLeader"`
		StaticPodManagerRunning bool `json:"staticPodManagerRunning"`
	}{
		IsLeader:                status.isLeader,
		StaticPodManagerRunning: status.staticPodManagerRunning,
	}

	w.Header().Set("Content-Type", "application/json")
	if status.staticPodManagerRunning {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}

	_ = json.NewEncoder(w).Encode(response)
}
