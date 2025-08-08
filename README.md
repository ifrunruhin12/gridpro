# GridPro - The Smartest Connect-4 AI

> An AI-powered Connect-4 game built in Go — born from a real-life defeat and rebuilt for total domination. The mission? Never lose again.

## Project Status

 **Live & Learning** — Features a competitive Connect-4 AI with Minimax and Alpha-Beta pruning. Fully playable via a modern web UI.

## AI Implementation

The AI uses a combination of advanced techniques to provide a challenging opponent:

### Minimax Algorithm with Alpha-Beta Pruning
- **Iterative Deepening + Time Budget**: Searches progressively deeper up to `MaxDepth` (currently 10) within ~`TimeLimitMs` (900ms) per AI move
- **Alpha-Beta Pruning**: Optimizes the search by eliminating branches that won't affect the final decision
- **Heuristic Evaluation**: Evaluates board positions based on:
  - Center control (prioritizes center columns)
  - Potential winning moves (threat detection)
  - Defensive play (blocks opponent's winning moves)

### Evaluation Function
- **Center Control**: Bonus points for pieces in central columns
- **Potential Wins**: Detects and prioritizes creating multiple winning opportunities
- **Defensive Play**: Recognizes and blocks opponent's threats
- **Win/Loss Detection**: Immediate evaluation of terminal states

### Performance Optimizations
- **Opening Book**: Hard-coded optimal early moves to accelerate strong starts
- **Preferred Move Ordering**: Center-out ordering improves pruning
- **Suicidal Move Filter**: Avoids moves that allow the opponent an immediate win
- **Transposition Table**: Caches board evaluations to avoid recomputation
- **Early Termination**: Respects a strict time limit per move

## Project Structure

```
gridpro/
├── backend/
│   ├── ai.go        # AI logic (Minimax, Alpha-Beta, iterative deepening, opening book)
│   ├── board.go     # Game logic (Connect 4 rules, state, etc.)
│   └── server.go    # HTTP API handlers and server setup
├── docs/
│   ├── index.html   # Web UI for Connect 4
│   ├── style.css    # UI styles
│   └── script.js    # Frontend logic (talks to backend API)
├── main.go          # Entry point for backend server
├── go.mod           # Go module definition
└── README.md        # This file
```

## How to Run

### 1. Backend (Go API)
- Make sure you have Go installed (1.18+ recommended).
- Start the backend server:
  ```sh
  go run main.go
  ```
- The backend will listen on `http://localhost:8080`.

#### API Endpoints
- `POST   /api/new`   — Start a new game. Returns a game ID and initial state.
- `POST   /api/move`  — Make a move. Body: `{ "gameId": string, "col": int }`. Returns updated state.
- `GET    /api/state` — Get current state. Query: `?gameId=...`.
- `GET    /api/info`  — Get backend info (version, opening preferences).

### 2. Frontend (Web UI)
- Serve the `docs/` folder with any static file server. For example:
  ```sh
  npx serve docs -l 3000
  # or
  python3 -m http.server 3000 --directory docs
  ```
- Open [http://localhost:3000](http://localhost:3000) in your browser.
- The frontend will talk to the backend API at `localhost:8080`.

## Deployment

- Frontend: Deployed via GitHub Pages.
- Backend: Hosted on Render.
- Make sure to update the API base URL in `script.js` if deploying somewhere else.

## Customizing the AI

You can adjust the AI's difficulty by modifying these constants in `backend/ai.go`:
- `MaxDepth`: Maximum search depth (default 10). Higher = stronger but slower.
- `TimeLimitMs`: Per-move search budget in milliseconds (default 900ms).
- Evaluation function weights in `evaluateBoard()`

## Future Improvements

- [ ] Add difficulty levels (Easy, Medium, Hard)
- [ ] Add move history and undo functionality
- [ ] Optimize evaluation function for better performance
- [ ] Persistent transposition table across moves within a game with better keying

## Tech Stack
- **Language**: Go (Golang)
- **Frontend**: HTML/CSS/JavaScript
- **Backend**: Go net/http
- **AI**: Minimax with Alpha-Beta pruning

## Why GridPro?

Connect-4 is a **solved game**. With perfect play, the first player can always win. GridPro implements advanced AI techniques to provide a challenging opponent that adapts to your skill level.

## License

MIT 2025 popcycle
