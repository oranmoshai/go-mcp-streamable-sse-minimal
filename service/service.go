package service

import (
	"net/http"
	"log"
	"github.com/oranmoshai/go-mcp-streamable-sse-minimal/usecase"
)

type Service struct {
}

func New() *Service {
	return &Service{}
}

func (s *Service) Run() (err error) {
	u := usecase.New()
	http.HandleFunc("/sse", u.EventStreamHandler)

	// Start the HTTP server on port 8080.
	port := "8080"
	log.Printf("Starting SSE server on http://localhost:%s", port)

	// ListenAndServe will block until the server is stopped or an error occurs.
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	return nil
}
