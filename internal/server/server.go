package server

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/YaroslavGaponov/tigger/internal/api"
	"github.com/YaroslavGaponov/tigger/internal/oracle"
)

type Server struct {
	oracle *oracle.Oracle
}

func NewServer(oracle *oracle.Oracle) (*Server, error) {
	return &Server{oracle: oracle}, nil
}

func (s *Server) Run(ctx context.Context) error {
	return s.oracle.Run(ctx)
}

// Ask Ask a question
// (POST /ask)
func (s *Server) Ask(w http.ResponseWriter, r *http.Request) {
	var request api.Question
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	answer, err := s.oracle.Ask(r.Context(), request.Kind, request.Question)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := api.Answer{
		Answer: answer,
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

// Health Health check
// (GET /health)
func (s *Server) Health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}

// GetChat Get chat interface
// (GET /)
func (s *Server) GetChat(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(ChatHTML))
}
