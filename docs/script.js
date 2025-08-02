const ROWS = 6;
const COLS = 7;
const PLAYER = 1;
const AI = 2;
const API_URL = 'https://gridgod.onrender.com/api';

let board = [];
let currentPlayer = PLAYER;
let gameActive = true;
let gameId = null;

const boardDiv = document.getElementById('board');
const statusDiv = document.getElementById('status');
const restartBtn = document.getElementById('restart');

async function newGame() {
    setStatus('Loading...');
    try {
        const res = await fetch(`${API_URL}/new`, { method: 'POST' });
        if (!res.ok) {
            setStatus('Failed to start game: Backend error.');
            gameActive = false;
            return;
        }
        const data = await res.json();
        gameId = data.gameId;
        updateFromState(data);
        gameActive = true;
    } catch (err) {
        setStatus('Failed to connect to backend. Is the server running?');
        gameActive = false;
    }
}

function updateFromState(data) {
    board = data.state.grid;
    currentPlayer = data.state.current_turn;
    renderBoard();
    if (data.checkWin && data.checkWin !== 0) {
        setStatus(data.checkWin === PLAYER ? `🎉 You win! (🟡)` : `🤖 AI wins! (🔴)`);
        gameActive = false;
    } else if (data.isDraw) {
        setStatus(`🤝 It's a draw!`);
        gameActive = false;
    } else {
        setStatus(currentPlayer === PLAYER ? `Your turn (🟡)` : `AI's turn... (🔴)`);
    }
}

function renderBoard() {
    boardDiv.innerHTML = '';
    for (let i = 0; i < ROWS; i++) {
        for (let j = 0; j < COLS; j++) {
            const cell = document.createElement('div');
            cell.classList.add('cell');
            if (board[i][j] === PLAYER) cell.classList.add('player1'); // yellow
            if (board[i][j] === AI) cell.classList.add('player2');     // red
            cell.dataset.row = i;
            cell.dataset.col = j;
            cell.addEventListener('click', handleCellClick);
            boardDiv.appendChild(cell);
        }
    }
}

async function handleCellClick(e) {
    if (!gameActive) return;
    const col = parseInt(e.target.dataset.col);
    if (board[0][col] !== 0) {
        setStatus('That column is full. Try another.');
        return;
    }

    const res = await fetch(`${API_URL}/move`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ gameId, col })
    });

    if (!res.ok) {
        setStatus('Invalid move. Try another.');
        return;
    }

    const data = await res.json();
    updateFromState(data);
}

function setStatus(msg) {
    statusDiv.textContent = msg;
}

restartBtn.addEventListener('click', newGame);

newGame();

