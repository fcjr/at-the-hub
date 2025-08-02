package server

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/fcjr/at-the-hub/internal/middleware"
	"github.com/fcjr/at-the-hub/internal/printer"
	"github.com/fcjr/at-the-hub/internal/recurse"
)

const defaultReadTimeout = 5 * time.Second
const defaultWriteTimeout = 10 * time.Second
const defaultShutdownGracePeriod = 30 * time.Second

type Server struct {
	logger          *slog.Logger
	shouldServeDocs bool
	printer         *printer.Printer
	recurseClient   *recurse.Client
}

type NewParams struct {
	Logger          *slog.Logger
	ShouldServeDocs bool
	Printer         *printer.Printer
	RecurseClient   *recurse.Client
}

func New(params NewParams, opts ...func(*Server) error) (*Server, error) {
	s := &Server{
		logger:          params.Logger,
		shouldServeDocs: params.ShouldServeDocs,
		printer:         params.Printer,
		recurseClient:   params.RecurseClient,
	}
	for _, opt := range opts {
		if err := opt(s); err != nil {
			return nil, err
		}
	}
	return s, nil
}

func (s *Server) ListenAndServe(ctx context.Context, addr string) error {
	mux := http.NewServeMux()
	if s.shouldServeDocs {
		mux.HandleFunc("GET /openapi.json", s.handleSchema)
		mux.HandleFunc("GET /docs", s.handleDocs)
		mux.HandleFunc("GET /", s.redirectTo("/docs"))
	}

	mux.HandleFunc("POST /api/v1/print_checkins", s.handlePrintCheckins)

	// apply middlewares
	handler := middleware.Chain(mux,
		middleware.WithPanicRecovery(s.logger),
		middleware.WithRequestResponseLogging(s.logger),
	)

	httpServer := &http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  defaultReadTimeout,
		WriteTimeout: defaultWriteTimeout,
		BaseContext: func(listener net.Listener) context.Context {
			return ctx
		},
	}

	// start server
	s.logger.Info("server listening",
		"addr", addr,
	)

	errCh := make(chan error, 1)
	go func() {
		errCh <- httpServer.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		s.logger.Info("server shutting down")
		ctx, cancel := context.WithTimeout(context.Background(), defaultShutdownGracePeriod)
		defer cancel()
		return httpServer.Shutdown(ctx)
	case err := <-errCh:
		s.logger.Error("server listen error, shuting down",
			"err", err,
		)
		return err
	}
}

func (s *Server) redirectTo(path string) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		http.Redirect(res, req, path, http.StatusTemporaryRedirect)
	}
}
