package gapi

import "log/slog"

type HttpService struct {
	Addr string
}

func (s *HttpService) Listen() {
	slog.Info("HTTP service is listening", "addr", s.Addr)
}
