package handler

import (
	"context"
	"log"
	"net/http"
)

type Server struct {
	server *http.Server
}

func (s *Server) Run(ctx context.Context) error {
	
	log.Println("Server listening on", s.server.Addr)
	log.Fatal(s.server.ListenAndServe())

	return nil
}

func NewServer(ctx context.Context, h *Handler, port string) *Server {
	return &Server{server: &http.Server{
		Addr: ":"+port,
		Handler: h.Router(),
	}}
}