package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/codeforge-ide/codeforgeai.go/config"
	"github.com/codeforge-ide/codeforgeai.go/engine"
	"github.com/codeforge-ide/codeforgeai.go/linter"
)

// Server is the API server.
type Server struct {
	engine *engine.Engine
	linter *linter.Linter
}

// NewServer creates a new API server.
func NewServer() *Server {
	cfg, _ := config.EnsureConfigPrompts("")
	return &Server{
		engine: engine.NewEngine(&cfg),
		linter: linter.NewLinter(),
	}
}

// Start starts the API server.
func (s *Server) Start(port int) {
	http.HandleFunc("/analyze", s.handleAnalyze)
	http.HandleFunc("/prompt", s.handlePrompt)
	http.HandleFunc("/lint", s.handleLint)
	fmt.Printf("Starting server on port %d...\n", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

func (s *Server) handleAnalyze(w http.ResponseWriter, r *http.Request) {
	s.engine.RunAnalysis()
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Analysis complete."))
}

type PromptRequest struct {
	Prompt string `json:"prompt"`
}

type PromptResponse struct {
	Response string `json:"response"`
}

func (s *Server) handlePrompt(w http.ResponseWriter, r *http.Request) {
	var req PromptRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := s.engine.ProcessPrompt(req.Prompt)
	json.NewEncoder(w).Encode(PromptResponse{Response: resp})
}

type LintRequest struct {
	Path string `json:"path"`
}

func (s *Server) handleLint(w http.ResponseWriter, r *http.Request) {
	var req LintRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	issues, err := s.linter.LintFile(req.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(issues)
}
