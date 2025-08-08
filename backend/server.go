package backend

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

type GameStore struct {
	games map[string]*Board
	mu    sync.Mutex
}

type GameResponse struct {
	GameId   string `json:"gameId,omitempty"`
	State    *Board `json:"state"`
	CheckWin int    `json:"checkWin"`
	IsDraw   bool   `json:"isDraw"`
}

type MoveRequest struct {
	GameId string `json:"gameId"`
	Col    int    `json:"col"`
}

type Server struct {
	store *GameStore
	mux   *http.ServeMux
}

const BackendVersion = "v1.1-center-opening"

func NewHandler() http.Handler {
	store := &GameStore{games: make(map[string]*Board)}
	s := &Server{store: store, mux: http.NewServeMux()}
	s.mux.HandleFunc("/api/new", s.newGameHandler)
	s.mux.HandleFunc("/api/move", s.moveHandler)
	s.mux.HandleFunc("/api/state", s.getStateHandler)
	// Info endpoint for debugging deployments
	s.mux.HandleFunc("/api/info", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"version":       BackendVersion,
			"centerOpening": true,
			"preferredCols": preferredCols,
		})
	})
	return enableCORS(s.mux)
}

func (s *Server) newGameHandler(w http.ResponseWriter, r *http.Request) {
	s.store.mu.Lock()
	defer s.store.mu.Unlock()

	id := fmt.Sprintf("game-%d", len(s.store.games)+1)
	board := &Board{
		CurrentTurn: AI, // AI starts first
		LastMoveRow: -1,
		LastMoveCol: -1,
	}

	// Reset transposition table for a fresh game
	transposition = map[string]int{}

	// AI makes the first move (force center column)
	if countPieces(board) == 0 && board.IsValidMove(3) {
		board.Drop(3)
	} else {
		aiMove := GetAIMove(board)
		if aiMove != -1 {
			board.Drop(aiMove)
		}
	}

	s.store.games[id] = board

	resp := GameResponse{
		GameId:   id,
		State:    board,
		CheckWin: board.CheckWin(),
		IsDraw:   board.IsDraw(),
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
	board, ok := s.store.games[req.GameId]
	s.store.mu.Unlock()

	if !ok {
		http.Error(w, "Game not found", http.StatusNotFound)
		return
	}

	if board.CurrentTurn != Player {
		http.Error(w, "Not your turn", http.StatusBadRequest)
		return
	}

	// Human move
	if !board.Drop(req.Col) {
		http.Error(w, "Invalid move", http.StatusBadRequest)
		return
	}

	// Check after human move
	if win := board.CheckWin(); win != 0 || board.IsDraw() {
		resp := GameResponse{
			State:    board,
			CheckWin: win,
			IsDraw:   board.IsDraw(),
		}
		json.NewEncoder(w).Encode(resp)
		return
	}

	// AI move
	aiMove := GetAIMove(board)
	if aiMove != -1 {
		board.Drop(aiMove)
	}

	resp := GameResponse{
		State:    board,
		CheckWin: board.CheckWin(),
		IsDraw:   board.IsDraw(),
	}
	json.NewEncoder(w).Encode(resp)
}

func (s *Server) getStateHandler(w http.ResponseWriter, r *http.Request) {
	gameId := r.URL.Query().Get("gameId")
	s.store.mu.Lock()
	board, ok := s.store.games[gameId]
	s.store.mu.Unlock()
	if !ok {
		http.Error(w, "Game not found", http.StatusNotFound)
		return
	}
	resp := GameResponse{
		State:    board,
		CheckWin: board.CheckWin(),
		IsDraw:   board.IsDraw(),
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
