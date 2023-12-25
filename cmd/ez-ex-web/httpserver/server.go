package httpserver

import (
	"context"
	"net/http"
	"time"
)

const (
	defaultAddr            = ":80"
	defaultReadTimeout     = 4 * time.Second
	defaultWriteTimeout    = 4 * time.Second
	defaultShutdownTimeout = 4 * time.Second
)

// ServerConfig represent an HTTP server
type ServerConfig struct {
	server          *http.Server
	notify          chan error
	shutdownTimeout time.Duration
}

// Start create a new HTTP server and runs it
func Start(handler http.Handler, opts ...ServerOptionFactory) *ServerConfig {
	httpServer := &http.Server{
		Handler:      handler,
		ReadTimeout:  defaultReadTimeout,
		WriteTimeout: defaultWriteTimeout,
		Addr:         defaultAddr,
	}

	s := &ServerConfig{
		server:          httpServer,
		notify:          make(chan error, 1),
		shutdownTimeout: defaultShutdownTimeout,
	}

	// Custom options
	for _, opt := range opts {
		opt(s)
	}

	s.start()

	return s
}

func (s *ServerConfig) start() {
	go func() {
		s.notify <- s.server.ListenAndServe()
		close(s.notify)
	}()
}

// Notify errors while listening
func (s *ServerConfig) Notify() <-chan error {
	return s.notify
}

// Shutdown attempts to gracefully shut down the HTTP server
func (s *ServerConfig) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	return s.server.Shutdown(ctx)
}
