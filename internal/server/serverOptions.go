package server

import "log/slog"

func WithLogger(logger *slog.Logger) func(*Server) error {
	return func(s *Server) error {
		s.logger = logger
		return nil
	}
}

func WithServeDocs() func(*Server) error {
	return func(s *Server) error {
		s.shouldServeDocs = true
		return nil
	}
}
