package html

import "net/http"

// Server hosts the web client and manages WebSocket connections to stream coordinates.
type Server struct {
	Addr string
}

func (s *Server) Start() error {
	return http.ListenAndServe(s.Addr, nil)
}
