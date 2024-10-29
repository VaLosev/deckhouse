package main

import (
	"context"
	staticpodmanager "embeded-registry-manager/internal/static-pod"
	"encoding/json"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"

	"os"
	"time"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"

	//_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	"embeded-registry-manager/internal/controllers"
	httpclient "embeded-registry-manager/internal/utils/http_client"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
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

	// Create context with cancel function for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle OS signals for graceful shutdown
	go handleShutdown(cancel)

	// Start static pod manager
	go startStaticPodManager(ctx, kubeClient, status)

	// Start health server for readiness and liveness probes
	go startHealthServer(status)

	// Set up and start manager
	if err := setupAndStartManager(ctx, cfg, kubeClient, HttpClient, status); err != nil {
		ctrl.Log.Error(err, "Failed to start the embedded registry manager", "component", "main")
		os.Exit(1)
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

	// Add leader status update runnable
	err = mgr.Add(manager.RunnableFunc(func(ctx context.Context) error {
		status.isLeader = true
		<-ctx.Done()
		status.isLeader = false
		return nil
	}))
	if err != nil {
		return fmt.Errorf("unable to add leader runnable: %w", err)
	}

	// Create registry controller
	reconciler := &controllers.RegistryReconciler{
		Client:     mgr.GetClient(),
		Scheme:     mgr.GetScheme(),
		KubeClient: kubeClient,
		Recorder:   mgr.GetEventRecorderFor("embedded-registry-controller"),
		HttpClient: httpClient,
	}
	if err := reconciler.SetupWithManager(mgr, ctx); err != nil {
		return fmt.Errorf("unable to create controller: %w", err)
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
	/*	if err != nil {
			// If not in cluster, try to load kubeconfig from home directory
			cfg, err = clientcmd.BuildConfigFromFlags("", filepath.Join(homeDir(), ".kube", "config"))
		}
	*/
	return cfg, err
}

/*
// homeDir returns the home directory of the current user
func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE")
}
*/

// handleShutdown listens for system termination signals and cancels the context for graceful shutdown
func handleShutdown(cancel context.CancelFunc) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
	ctrl.Log.Info("Received shutdown signal")
	cancel()
}

// startStaticPodManager starts the static pod manager and monitors its status
func startStaticPodManager(ctx context.Context, kubeClient *kubernetes.Clientset, status *managerStatus) {
	status.staticPodManagerRunning = true
	if err := staticpodmanager.Run(ctx, kubeClient); err != nil {
		ctrl.Log.Error(err, "Failed to run static pod manager", "component", "main")
		status.staticPodManagerRunning = false
		os.Exit(1)
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
