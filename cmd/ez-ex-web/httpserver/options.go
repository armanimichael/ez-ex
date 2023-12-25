package httpserver

import (
	"net"
	"time"
)

// ServerOptionFactory provides a reference to the HTTP ServerConfig configuration
type ServerOptionFactory func(*ServerConfig)

func WithPort(port string) ServerOptionFactory {
	return func(s *ServerConfig) {
		s.server.Addr = net.JoinHostPort("", port)
	}
}

func WithReadTimeout(timeout time.Duration) ServerOptionFactory {
	return func(s *ServerConfig) {
		s.server.ReadTimeout = timeout
	}
}

func WriteTimeout(timeout time.Duration) ServerOptionFactory {
	return func(s *ServerConfig) {
		s.server.WriteTimeout = timeout
	}
}

func ShutdownTimeout(timeout time.Duration) ServerOptionFactory {
	return func(s *ServerConfig) {
		s.shutdownTimeout = timeout
	}
}
