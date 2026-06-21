package html

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// Frame represents a single animation frame sent to clients.
type Frame struct {
	Lines [][4]float64 `json:"lines"` // x1, y1, x2, y2 in screen coordinates
	Clear bool         `json:"clear"`
}

// Server hosts the web client and manages WebSocket connections.
type Server struct {
	Addr      string
	clients   map[*websocket.Conn]bool
	mu        sync.Mutex
	StaticDir string // path to web/ directory
}

// NewServer creates a new web server.
func NewServer(addr, staticDir string) *Server {
	return &Server{
		Addr:      addr,
		clients:   make(map[*websocket.Conn]bool),
		StaticDir: staticDir,
	}
}

// Start begins listening and serving.
func (s *Server) Start() error {
	// Serve static files
	fs := http.FileServer(http.Dir(s.StaticDir))
	http.Handle("/", fs)

	// WebSocket endpoint
	http.HandleFunc("/ws", s.handleWebSocket)

	fmt.Printf("Serving %s on %s\n", s.StaticDir, s.Addr)
	return http.ListenAndServe(s.Addr, nil)
}

func (s *Server) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("WebSocket upgrade error: %v\n", err)
		return
	}

	s.mu.Lock()
	s.clients[conn] = true
	s.mu.Unlock()

	fmt.Printf("Client connected (%d total)\n", len(s.clients))

	// Keep connection alive until client disconnects
	go func() {
		defer func() {
			s.mu.Lock()
			delete(s.clients, conn)
			s.mu.Unlock()
			conn.Close()
			fmt.Printf("Client disconnected (%d remaining)\n", len(s.clients))
		}()

		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				break
			}
		}
	}()
}

// Broadcast sends a frame to all connected clients.
func (s *Server) Broadcast(frame Frame) error {
	data, err := json.Marshal(frame)
	if err != nil {
		return fmt.Errorf("marshaling frame: %w", err)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	for conn := range s.clients {
		err := conn.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			fmt.Printf("Write error, removing client: %v\n", err)
			conn.Close()
			delete(s.clients, conn)
		}
	}

	return nil
}

// ClientCount returns the number of connected clients.
func (s *Server) ClientCount() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.clients)
}
