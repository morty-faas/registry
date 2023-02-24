package registry

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/polyxia-org/morty-registry/internal/storage"
	"github.com/polyxia-org/morty-registry/internal/storage/s3"
	"github.com/polyxia-org/morty-registry/pkg/config"
)

const (
	healthEndpoint = "/healthz"
)

type Server struct {
	cfg     *config.Config
	storage storage.Storage
}

func NewServer() (*Server, error) {
	cfg, err := config.Parse()
	if err != nil {
		return nil, err
	}

	s3, err := s3.NewStorage(cfg.Storage.S3)
	if err != nil {
		return nil, err
	}

	return &Server{cfg: cfg, storage: s3}, nil
}

func (s *Server) Serve() {
	p := s.cfg.Port

	ctx, stop := context.WithCancel(context.Background())
	server := &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%d", p),
		Handler: s.router(),
	}

	// Listen for syscall signals for process to interrupt/quit
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sigs

		// Shutdown signal with grace period of 30 seconds
		shutdownCtx, cancel := context.WithTimeout(ctx, 30*time.Second)

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				cancel()
				log.Fatal("graceful shutdown timed out... forcing exit")
			}
		}()

		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Fatal(err)
		}

		stop()
	}()

	log.Printf("Server is listening on : %d\n", p)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}

	<-ctx.Done()
}

func (s *Server) router() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)

	r.Get("/v1/functions/{id}/upload-link", s.UploadHandler)
	r.Get(healthEndpoint, s.HealthcheckHandler)

	return r
}
