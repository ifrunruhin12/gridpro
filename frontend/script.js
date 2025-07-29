// Connect 4 Frontend Game Logic (API version)
const ROWS = 6;
const COLS = 7;
const PLAYER1 = 1;
const PLAYER2 = -1;
const API_URL = 'http://localhost:8080/api';

let board = [];
let currentPlayer = PLAYER1;
let gameActive = true;
let gameId = null;

const boardDiv = document.getElementById('board');
const statusDiv = document.getElementById('status');
const restartBtn = document.getElementById('restart');

async function newGame() {
    const res = await fetch(`${API_URL}/new`, { method: 'POST' });
    const data = await res.json();
    gameId = data.gameId;
    updateFromState(data);
    setStatus(`Player 1's turn (🔴)`);
    gameActive = true;
}

function updateFromState(data) {
    board = data.state.Board;
    currentPlayer = data.state.CurrentPlayer;
    renderBoard();
    if (data.checkWin && data.checkWin !== 0) {
        setStatus(`🎉 Player ${data.checkWin === 1 ? '1 (🔴)' : '2 (🟡)'} wins!`);
        gameActive = false;
    } else if (data.isDraw) {
        setStatus(`🤝 It's a draw!`);
        gameActive = false;
    }
}

function renderBoard() {
    boardDiv.innerHTML = '';
    for (let i = 0; i < ROWS; i++) {
        for (let j = 0; j < COLS; j++) {
            const cell = document.createElement('div');
            cell.classList.add('cell');
            if (board[i][j] === PLAYER1) cell.classList.add('player1');
            if (board[i][j] === PLAYER2) cell.classList.add('player2');
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
    // Only allow move if top cell is empty
    if (board[0][col] !== 0) {
        setStatus('That column is full. Try another.');
        return;
    }
    // Send move to backend
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
    if (data.checkWin && data.checkWin !== 0) {
        setStatus(`🎉 Player ${data.checkWin === 1 ? '1 (🔴)' : '2 (🟡)'} wins!`);
        gameActive = false;
    } else if (data.isDraw) {
        setStatus(`🤝 It's a draw!`);
        gameActive = false;
    } else {
        setStatus(`Player ${data.state.CurrentPlayer === 1 ? "1's turn (🔴)" : "2's turn (🟡)"}`);
    }
}

function setStatus(msg) {
    statusDiv.textContent = msg;
}

restartBtn.addEventListener('click', newGame);

newGame(); 