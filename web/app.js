const canvas = document.getElementById('canvas');
const ctx = canvas.getContext('2d');
const fpsEl = document.getElementById('fps');
const clientsEl = document.getElementById('clients');

let width, height;
let frameCount = 0;
let lastFpsTime = performance.now();
let fps = 0;
let lines = [];

function resize() {
    width = window.innerWidth;
    height = window.innerHeight;
    canvas.width = width;
    canvas.height = height;
}

function draw() {
    // Black background
    ctx.fillStyle = '#000';
    ctx.fillRect(0, 0, width, height);

    // Green wireframe lines
    ctx.strokeStyle = '#0f0';
    ctx.lineWidth = 1.5;
    ctx.lineCap = 'round';

    ctx.beginPath();
    for (const [x1, y1, x2, y2] of lines) {
        ctx.moveTo(x1, y1);
        ctx.lineTo(x2, y2);
    }
    ctx.stroke();

    frameCount++;
}

function updateFPS() {
    const now = performance.now();
    const elapsed = now - lastFpsTime;
    if (elapsed >= 1000) {
        fps = Math.round((frameCount / elapsed) * 1000);
        frameCount = 0;
        lastFpsTime = now;
        fpsEl.textContent = `${fps} FPS`;
    }
}

function connect() {
    const protocol = location.protocol === 'https:' ? 'wss' : 'ws';
    const ws = new WebSocket(`${protocol}://${location.host}/ws`);

    ws.onopen = () => {
        console.log('WebSocket connected');
    };

    ws.onmessage = (event) => {
        const frame = JSON.parse(event.data);
        lines = frame.lines;
        clientsEl.textContent = `${frame.clients || 1} client(s)`;
    };

    ws.onclose = () => {
        console.log('WebSocket disconnected, reconnecting...');
        setTimeout(connect, 1000);
    };

    ws.onerror = (err) => {
        console.error('WebSocket error:', err);
    };
}

function loop() {
    draw();
    updateFPS();
    requestAnimationFrame(loop);
}

window.addEventListener('resize', resize);
resize();
connect();
loop();