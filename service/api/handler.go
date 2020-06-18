package api

import (
	"log"
	"sync"

	"golang.org/x/net/context"
)

// Server represents the gRPC server
type Server struct {
	C  Controlable
	WG *sync.WaitGroup
}

// GetStatus ....
func (s *Server) GetStatus(ctx context.Context, in *Empty) (*Status, error) {
	log.Printf("GetStatus is called")
	return &Status{Status: s.C.GetStatus()}, nil
}

// Start ...
func (s *Server) Start(ctx context.Context, in *Empty) (*Status, error) {
	log.Printf("Start is called")
	status := s.C.Start()
	return &Status{Status: status}, nil
}

// Stop ..
func (s *Server) Stop(ctx context.Context, in *Empty) (*Status, error) {
	log.Printf("GetStatus is called")
	return &Status{Status: s.C.Stop()}, nil
}
