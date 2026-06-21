// Constants & Variables
const canvas = document.getElementById('renderCanvas');
const ctx = canvas.getContext('2d');
const overlay = document.getElementById('overlay');
const connectBtn = document.getElementById('connectBtn');
const statusDot = document.querySelector('.status-dot');
const statusText = document.querySelector('.status-text');
const projectionModeSelect = document.getElementById('projectionMode');
const rotationSpeedSlider = document.getElementById('rotationSpeed');
const lineColorPicker = document.getElementById('lineColor');
const colorValueSpan = document.querySelector('.color-value');

let rotationSpeed = parseFloat(rotationSpeedSlider.value) / 1000;
let lineColor = lineColorPicker.value;
let projectionMode = projectionModeSelect.value;
let angleX = 0;
let angleY = 0;
let isConnected = false;

// Resize canvas to match container resolution
function resizeCanvas() {
    const rect = canvas.getBoundingClientRect();
    canvas.width = rect.width * window.devicePixelRatio;
    canvas.height = rect.height * window.devicePixelRatio;
    ctx.scale(window.devicePixelRatio, window.devicePixelRatio);
}
window.addEventListener('resize', resizeCanvas);
resizeCanvas();

// Standard 3D cube geometry
const vertices = [
    {x: -1, y: -1, z: -1},
    {x: 1, y: -1, z: -1},
    {x: 1, y: 1, z: -1},
    {x: -1, y: 1, z: -1},
    {x: -1, y: -1, z: 1},
    {x: 1, y: -1, z: 1},
    {x: 1, y: 1, z: 1},
    {x: -1, y: 1, z: 1}
];

const edges = [
    [0, 1], [1, 2], [2, 3], [3, 0], // Back face
    [4, 5], [5, 6], [6, 7], [7, 4], // Front face
    [0, 4], [1, 5], [2, 6], [3, 7]  // Connecting edges
];

// Project 3D points to 2D screen coordinates
function project(vertex, width, height) {
    let x = vertex.x;
    let y = vertex.y;
    let z = vertex.z;

    // Rotate X
    let cosX = Math.cos(angleX);
    let sinX = Math.sin(angleX);
    let y1 = y * cosX - z * sinX;
    let z1 = y * sinX + z * cosX;

    // Rotate Y
    let cosY = Math.cos(angleY);
    let sinY = Math.sin(angleY);
    let x2 = x * cosY + z1 * sinY;
    let z2 = -x * sinY + z1 * cosY;

    let scale = 120;
    let projectedX = 0;
    let projectedY = 0;

    if (projectionMode === 'perspective') {
        // Perspective projection: divide by z
        let distance = 3.0;
        let zPos = z2 + distance;
        projectedX = (x2 / zPos) * width * 0.4 + width / 2;
        projectedY = (y1 / zPos) * height * 0.5 + height / 2;
    } else {
        // Orthographic projection: flat scale
        projectedX = x2 * scale + width / 2;
        projectedY = y1 * scale + height / 2;
    }

    return { x: projectedX, y: projectedY };
}

// Visualizer Main Loop
function tick() {
    const width = canvas.width / window.devicePixelRatio;
    const height = canvas.height / window.devicePixelRatio;

    // Update rotations
    angleX += rotationSpeed;
    angleY += rotationSpeed * 1.5;

    // Draw frame (if not streaming web socket data)
    if (!isConnected) {
        ctx.clearRect(0, 0, width, height);
        
        ctx.strokeStyle = lineColor;
        ctx.lineWidth = 2.5;
        ctx.shadowBlur = 15;
        ctx.shadowColor = lineColor;
        ctx.lineCap = 'round';
        ctx.lineJoin = 'round';

        // Project all points
        const projected = vertices.map(v => project(v, width, height));

        // Draw edges
        for (const edge of edges) {
            const p1 = projected[edge[0]];
            const p2 = projected[edge[1]];

            ctx.beginPath();
            ctx.moveTo(p1.x, p1.y);
            ctx.lineTo(p2.x, p2.y);
            ctx.stroke();
        }
    }

    requestAnimationFrame(tick);
}
requestAnimationFrame(tick);

// UI Events
rotationSpeedSlider.addEventListener('input', (e) => {
    rotationSpeed = parseFloat(e.target.value) / 1000;
});

lineColorPicker.addEventListener('input', (e) => {
    lineColor = e.target.value;
    colorValueSpan.textContent = lineColor;
});

projectionModeSelect.addEventListener('change', (e) => {
    projectionMode = e.target.value;
});

connectBtn.addEventListener('click', () => {
    if (isConnected) {
        disconnect();
    } else {
        connect();
    }
});

function connect() {
    statusDot.className = 'status-dot';
    statusDot.style.backgroundColor = '#fbbf24'; // orange
    statusDot.style.boxShadow = '0 0 10px #fbbf24';
    statusText.textContent = 'Connecting...';
    
    // Simulate connection for visual purposes (placeholder for real WebSocket handshake)
    setTimeout(() => {
        isConnected = true;
        statusDot.className = 'status-dot connected';
        statusText.textContent = 'Connected (WS Live)';
        connectBtn.textContent = 'Disconnect Stream';
        connectBtn.style.background = 'linear-gradient(135deg, #ef4444, #b91c1c)';
        connectBtn.style.boxShadow = '0 4px 15px rgba(239, 68, 68, 0.2)';
        overlay.innerHTML = `<h2>Live Stream Active</h2><p>Receiving wireframe paths from Go WebSocket server.</p>`;
    }, 1000);
}

function disconnect() {
    isConnected = false;
    statusDot.className = 'status-dot pulsing';
    statusDot.style.backgroundColor = '#ef4444'; // red
    statusDot.style.boxShadow = '0 0 8px #ef4444';
    statusText.textContent = 'Disconnected';
    connectBtn.textContent = 'Connect Stream';
    connectBtn.style.background = 'linear-gradient(135deg, var(--accent-teal), var(--accent-violet))';
    connectBtn.style.boxShadow = '0 4px 15px rgba(0, 255, 204, 0.2)';
    overlay.innerHTML = `<h2>Ready to Render</h2><p>Start the Go WebSocket server to stream 3D asset wireframes.</p>`;
}
