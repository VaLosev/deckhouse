package server

import (
	"context"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"time"

	shapp "github.com/flant/shell-operator/pkg/app"
	"github.com/flant/shell-operator/pkg/kube"
	"github.com/flant/shell-operator/pkg/metric_storage"
	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"

	"d8.io/upmeter/pkg/crd"
	dbcontext "d8.io/upmeter/pkg/db/context"
	"d8.io/upmeter/pkg/db/dao"
	"d8.io/upmeter/pkg/db/migrations"
	"d8.io/upmeter/pkg/server/api"
	"d8.io/upmeter/pkg/server/remotewrite"
)

// server initializes all dependencies:
// - kubernetes client
// - crd monitor
// - database connection
// - metrics storage
// If everything is ok, it starts http server.

type server struct {
	config *Config

	logger *log.Logger

	server                *http.Server
	downtimeMonitor       *crd.DowntimeMonitor
	remoteWriteController *remotewrite.Controller
}

type Config struct {
	ListenHost string
	ListenPort string

	DatabasePath           string
	DatabaseMigrationsPath string

	OriginsCount int
}

func New(config *Config, logger *log.Logger) *server {
	return &server{
		config: config,
		logger: logger,
	}
}

func (s *server) Start(ctx context.Context) error {
	var err error

	kubeClient, err := initKubeClient()
	if err != nil {
		return fmt.Errorf("init kubernetes client: %v", err)
	}

	// Database connection with pool
	dbCtx, err := migrations.GetMigratedDatabase(ctx, s.config.DatabasePath, s.config.DatabaseMigrationsPath)
	if err != nil {
		return fmt.Errorf("database not connected: %v", err)
	}

	// Downtime CR monitor
	s.downtimeMonitor, err = initDowntimeMonitor(ctx, kubeClient)
	if err != nil {
		return fmt.Errorf("cannot start downtimes.deckhouse.io monitor: %v", err)
	}

	// Metrics controller
	s.remoteWriteController, err = initRemoteWriteController(ctx, dbCtx, kubeClient, s.config.OriginsCount, s.logger)
	if err != nil {
		s.logger.Debugf("starting controller... did't happen: %v", err)
		return fmt.Errorf("cannot start remote_write controller: %v", err)
	}

	go cleanOld30sEpisodes(ctx, dbCtx)

	// Start http server. It blocks, that's why it is the last here.
	s.logger.Debugf("starting HTTP server")
	listenAddr := s.config.ListenHost + ":" + s.config.ListenPort
	s.server = initHttpServer(dbCtx, s.downtimeMonitor, s.remoteWriteController, listenAddr)

	err = s.server.ListenAndServe()
	if err == http.ErrServerClosed {
		return nil
	}

	return err
}

func (s *server) Stop() error {
	err := s.server.Shutdown(context.Background())
	if err != nil && err != http.ErrServerClosed {
		return err
	}
	s.remoteWriteController.Stop()
	s.downtimeMonitor.Stop()

	return nil
}

func cleanOld30sEpisodes(ctx context.Context, dbCtx *dbcontext.DbContext) {
	dayBack := -24 * time.Hour
	period := 30 * time.Second

	conn := dbCtx.Start()
	defer conn.Stop()

	storage := dao.NewEpisodeDao30s(conn)

	ticker := time.NewTicker(period)

	for {
		select {
		case <-ticker.C:
			deadline := time.Now().Truncate(period).Add(dayBack)
			err := storage.DeleteUpTo(deadline)
			if err != nil {
				log.Errorf("cannot clean old episodes: %v", err)
			}
		case <-ctx.Done():
			ticker.Stop()
			return
		}
	}
}

func initHttpServer(dbCtx *dbcontext.DbContext, downtimeMonitor *crd.DowntimeMonitor, controller *remotewrite.Controller, addr string) *http.Server {
	mux := http.NewServeMux()

	// Setup API handlers
	mux.Handle("/api/probe", &api.ProbeListHandler{DbCtx: dbCtx})
	mux.Handle("/api/status/range", &api.StatusRangeHandler{DbCtx: dbCtx, DowntimeMonitor: downtimeMonitor})
	mux.Handle("/public/api/status", &api.PublicStatusHandler{DbCtx: dbCtx, DowntimeMonitor: downtimeMonitor})
	mux.Handle("/downtime", &api.AddEpisodesHandler{DbCtx: dbCtx, RemoteWrite: controller})
	mux.Handle("/stats", &api.StatsHandler{DbCtx: dbCtx})
	// Kubernetes probes
	mux.HandleFunc("/healthz", writeOk)
	mux.HandleFunc("/ready", writeOk)

	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	return server
}

func writeOk(w http.ResponseWriter, _ *http.Request) {
	_, _ = w.Write([]byte("OK"))
}

func initRemoteWriteController(ctx context.Context, dbCtx *dbcontext.DbContext, kubeClient kube.KubernetesClient, originsCount int, logger *log.Logger) (*remotewrite.Controller, error) {
	config := &remotewrite.ControllerConfig{
		// collecting/exporting episodes as metrics
		Period: 2 * time.Second,
		// monitor configs in kubernetes
		Kubernetes: kubeClient,
		// read metrics and track exporter state in the DB
		DbCtx:        dbCtx,
		OriginsCount: originsCount,
		Logger:       logger,
	}
	controller := config.Controller()
	return controller, controller.Start(ctx)
}

func initDowntimeMonitor(ctx context.Context, kubeClient kube.KubernetesClient) (*crd.DowntimeMonitor, error) {
	m := crd.NewMonitor(ctx)
	m.Monitor.WithKubeClient(kubeClient)
	return m, m.Start()
}

func initKubeClient() (kube.KubernetesClient, error) {
	client := kube.NewKubernetesClient()

	client.WithContextName(shapp.KubeContext)
	client.WithConfigPath(shapp.KubeConfig)
	client.WithRateLimiterSettings(shapp.KubeClientQps, shapp.KubeClientBurst)
	client.WithMetricStorage(metric_storage.NewMetricStorage())

	// FIXME: Kubernetes client is configured successfully with 'out-of-cluster' config
	//      operator.component=KubernetesAPIClient
	return client, client.Init()
}
