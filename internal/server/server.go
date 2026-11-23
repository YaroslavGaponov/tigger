package server

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/YaroslavGaponov/tigger/internal/api"
	"github.com/YaroslavGaponov/tigger/internal/config"
	"github.com/YaroslavGaponov/tigger/internal/core"
)

type Server struct {
	core *core.Core
}

func NewServer(config config.Config) (*Server, error) {
	return &Server{core: core.NewCore(config)}, nil
}

func (s *Server) Run(ctx context.Context) error {
	return s.core.Run(ctx)
}

// Ask Ask a question
// (POST /ask)
func (s *Server) Ask(w http.ResponseWriter, r *http.Request) {
	var request api.Question
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	answer, err := s.core.Ask(r.Context(), request.Question)
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
