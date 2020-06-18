package api

import (
	"log"

	"golang.org/x/net/context"
)

// Server represents the gRPC server
type Server struct {
}

// SayHello generates response to a Ping request
func (s *Server) IsReady(ctx context.Context, in *Empty) (*Status, error) {
	log.Printf("IsReady is called")
	return &Status{Ready: false, Done: false}, nil
}
