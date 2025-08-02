
# 🧠 GridGod - The Unbeatable Connect-4 AI

> A blazing-fast, AI-powered Connect-4 game built with Go — born from a crushing loss and reborn as an unrelenting grid overlord. You play yellow. AI plays red. Good luck. You'll need it.

---

## 📦 Project Status

✅ **Playable & Deployed** — RandomBot & GreedyBot complete.  
🧠 Next up: Lookahead AI (Minimax with Alpha-Beta Pruning).

**Play it now**: https://ifrunruhin12.github.io/gridgod  
**Backend API**: https://gridgod.onrender.com/api

---

## 🗂 Project Structure

```
gridgod/
├── backend/              # Backend codebase
│   ├── board.go          # Game logic (board state, win/draw, moves, etc.)
│   ├── server.go         # HTTP handlers & game state management
│   └── ai.go             # AI decision logic (Random, Greedy, etc.)
├── frontend/             # Web frontend
│   ├── index.html        # UI layout
│   ├── style.css         # Visual styles
│   └── script.js         # Game loop, rendering, API calls
├── main.go               # Entry point for backend server
├── go.mod                # Go module definition
└── README.md             # You're here.
```

---

## 🚀 How it Works

#### API Routes

| Method | Endpoint         | Description                        |
|--------|------------------|------------------------------------|
| POST   | `/api/new`       | Starts a new game (AI moves first) |
| POST   | `/api/move`      | Makes a move (by human player)     |
| GET    | `/api/state`     | Gets current game state            |

---

#### Frontend

If using the deployed version:
- Make sure frontend API points to `https://gridgod.onrender.com/api`

---

## 🌐 Deployment

| Part      | Platform         | Link                                        |
|-----------|------------------|---------------------------------------------|
| Backend   | Render           | https://gridgod.onrender.com                |
| Frontend  | GitHub Pages     | https://ifrunruhin12.github.io/gridgod      |

To update frontend:
1. Push to `gh-pages` branch (or deploy `/frontend` from `main`)
2. GitHub Pages will auto-publish the static site

---

## 🔮 Vision

- ✅ RandomBot ✅ GreedyBot  
- 🧠 Up next: Lookahead-1 AI  
- ⚔️ Future: Full Minimax + Alpha-Beta Pruning
- 🎮 Goal: Become **unbeatable** (first-move win guaranteed)

---

## 🛠️ Tech Stack

- 💻 **Language**: Go (Golang)
- 🌐 **Frontend**: HTML, CSS, JavaScript
- 🔗 **Backend**: net/http, REST API
- 🧠 **AI**: Step-by-step evolution (Random → Greedy → Minimax)

---

## 🧠 Why?

Connect-4 is a solved game. The first player can always win with perfect play.  
This project is my journey to building that perfect AI — and learning a ton about game engines, algorithms, and web systems in the process.

---

## 📄 License

MIT © 2025 [popcycle](https://github.com/ifrunruhin12)
