# 🧠 GridPro - The Smartest Connect-4 AI

> An AI-powered Connect-4 game built in Go — born from a real-life defeat and rebuilt for total domination. The mission? Never lose again.

## 📦 Project Status

🚧 **Under Construction** — Currently building a competitive Connect-4 AI with lookahead strategies. Fully playable via a modern web UI.

## 🔥 Project Structure

```
gridpro/
├── backend/
│   ├── board.go      # Game logic (Connect 4 rules, state, etc.)
│   └── server.go     # HTTP API handlers and server setup
├── frontend/
│   ├── index.html    # Web UI for Connect 4
│   ├── style.css     # UI styles
│   └── script.js     # Frontend logic (talks to backend API)
├── main.go           # Entry point for backend server
├── go.mod            # Go module definition
└── README.md         # This file
```

## 🚀 How to Run

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

### 2. Frontend (Web UI)
- Serve the `frontend/` folder with any static file server. For example:
  ```sh
  npx serve frontend -l 3000
  # or
  python3 -m http.server 3000 --directory frontend
  ```
- Open [http://localhost:3000](http://localhost:3000) in your browser.
- The frontend will talk to the backend API at `localhost:8080`.

## 🌐 Deployment

- Frontend: Deployed via GitHub Pages.
- Backend: Hosted on Render — https://gridgod.onrender.com/
- Make sure to update the API base URL in `script.js` if deploying somewhere else.

## 🔮 Vision

This project aims to:
- Build a fully functioning Connect-4 engine in Go.
- Implement a powerful AI that levels up: Random → Greedy → Lookahead → Minimax.
- Expose the AI through a RESTful API.
- Design a sleek web interface to challenge the AI.

## 🛠️ Tech Stack
- **Language**: Go (Golang)
- **Frontend**: HTML/CSS/JavaScript
- **Backend**: Go net/http

## 🧠 Why GridPro?

Connect-4 is a **solved game**. With perfect play, the first player can always win. GridPro is your shot at mastering it — one move at a time.

## 📌 License

MIT © 2025 popcycle

