package server

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all connections
	},
}

// Server manages WebSocket clients and broadcasts messages.
type Server struct {
	clients    map[*websocket.Conn]bool
	broadcast  chan []byte
	register   chan *websocket.Conn
	unregister chan *websocket.Conn
	mutex      sync.Mutex
}

// NewServer creates a new WebSocket server instance.
func NewServer() *Server {
	return &Server{
		clients:    make(map[*websocket.Conn]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *websocket.Conn),
		unregister: make(chan *websocket.Conn),
	}
}

// Run starts the WebSocket server's message handling loops.
func (s *Server) Run() {
	for {
		select {
		case conn := <-s.register:
			s.mutex.Lock()
			s.clients[conn] = true
			s.mutex.Unlock()
			log.Println("Client registered")
		case conn := <-s.unregister:
			s.mutex.Lock()
			if _, ok := s.clients[conn]; ok {
				delete(s.clients, conn)
				conn.Close()
				log.Println("Client unregistered")
			}
			s.mutex.Unlock()
		case message := <-s.broadcast:
			s.mutex.Lock()
			for conn := range s.clients {
				if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
					log.Printf("Error broadcasting message: %v", err)
					s.unregister <- conn
				}
			}
			s.mutex.Unlock()
		}
	}
}

// BroadcastMessage sends a message to all connected clients.
func (s *Server) BroadcastMessage(message []byte) {
	s.broadcast <- message
}

// HandleConnections upgrades HTTP requests to WebSocket connections.
func (s *Server) HandleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading connection: %v", err)
		return
	}
	s.register <- conn
}
