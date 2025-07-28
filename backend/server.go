package backend

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

type GameStore struct {
	games map[string]*GameState
	mu    sync.Mutex
}

type GameResponse struct {
	GameId   string     `json:"gameId,omitempty"`
	State    *GameState `json:"state"`
	CheckWin int        `json:"checkWin"`
	IsDraw   bool       `json:"isDraw"`
}

type MoveRequest struct {
	GameId string `json:"gameId"`
	Col    int    `json:"col"`
}

type Server struct {
	store *GameStore
	mux   *http.ServeMux
}

func NewHandler() http.Handler {
	store := &GameStore{games: make(map[string]*GameState)}
	s := &Server{store: store, mux: http.NewServeMux()}
	s.mux.HandleFunc("/api/new", s.newGameHandler)
	s.mux.HandleFunc("/api/move", s.moveHandler)
	s.mux.HandleFunc("/api/state", s.getStateHandler)
	return enableCORS(s.mux)
}

func (s *Server) newGameHandler(w http.ResponseWriter, r *http.Request) {
	s.store.mu.Lock()
	defer s.store.mu.Unlock()
	id := fmt.Sprintf("game-%d", len(s.store.games)+1)
	game := &GameState{CurrentPlayer: 1}
	s.store.games[id] = game
	resp := GameResponse{
		GameId:   id,
		State:    game,
		CheckWin: game.CheckWin(),
		IsDraw:   game.IsDraw(),
	}
	json.NewEncoder(w).Encode(resp)
}

func (s *Server) moveHandler(w http.ResponseWriter, r *http.Request) {
	var req MoveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	s.store.mu.Lock()
	game, ok := s.store.games[req.GameId]
	s.store.mu.Unlock()
	if !ok {
		http.Error(w, "Game not found", http.StatusNotFound)
		return
	}
	if !game.Drop(req.Col) {
		http.Error(w, "Invalid move", http.StatusBadRequest)
		return
	}
	resp := GameResponse{
		State:    game,
		CheckWin: game.CheckWin(),
		IsDraw:   game.IsDraw(),
	}
	json.NewEncoder(w).Encode(resp)
}

func (s *Server) getStateHandler(w http.ResponseWriter, r *http.Request) {
	gameId := r.URL.Query().Get("gameId")
	s.store.mu.Lock()
	game, ok := s.store.games[gameId]
	s.store.mu.Unlock()
	if !ok {
		http.Error(w, "Game not found", http.StatusNotFound)
		return
	}
	resp := GameResponse{
		State:    game,
		CheckWin: game.CheckWin(),
		IsDraw:   game.IsDraw(),
	}
	json.NewEncoder(w).Encode(resp)
}

func enableCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		h.ServeHTTP(w, r)
	})
}
