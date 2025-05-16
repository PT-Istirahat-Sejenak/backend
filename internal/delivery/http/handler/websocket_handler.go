package handler

import (
	"backend/internal/entity"
	"backend/internal/usecase"
	"context"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type WebSocketHandler struct {
	clients        map[*websocket.Conn]uint // Maps connection to userID
	connections    map[uint]*websocket.Conn // Maps userID to connection
	broadcast      chan entity.MessageRequest
	upgrader       websocket.Upgrader
	messageUseCase usecase.MessageUseCase
	mutex          sync.Mutex // For thread safety when manipulating maps
}

func NewWebSockerHandler(messageUseCase usecase.MessageUseCase) *WebSocketHandler {
	return &WebSocketHandler{
		clients:     make(map[*websocket.Conn]uint),
		connections: make(map[uint]*websocket.Conn),
		broadcast:   make(chan entity.MessageRequest),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		messageUseCase: messageUseCase,
		mutex:          sync.Mutex{},
	}
}

func (h *WebSocketHandler) HandleConnection(w http.ResponseWriter, r *http.Request) {
	// Upgrade the HTTP connection to a WebSocket connection
	ws, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}

	// Set up a close handler to clean up resources
	ws.SetCloseHandler(func(code int, text string) error {
		log.Printf("Connection closed with code %d: %s", code, text)
		h.cleanupConnection(ws)
		return nil
	})

	// Set read/write deadlines
	ws.SetReadDeadline(time.Now().Add(60 * time.Second))
	ws.SetPongHandler(func(string) error {
		ws.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})
	
	// Wait for authentication message containing user ID
	var authMsg struct {
		UserID uint `json:"user_id"`
	}
	
	if err := ws.ReadJSON(&authMsg); err != nil {
		log.Println("Authentication error:", err)
		ws.Close()
		return
	}

	userID := authMsg.UserID
	log.Printf("User %d connected via WebSocket", userID)
	
	// Register this connection
	h.registerConnection(ws, userID)
	
	// Send confirmation message
	confirmMsg := map[string]interface{}{
		"type":    "connected",
		"user_id": userID,
		"time":    time.Now(),
	}
	
	if err := ws.WriteJSON(confirmMsg); err != nil {
		log.Println("Failed to send confirmation:", err)
		h.cleanupConnection(ws)
		ws.Close()
		return
	}

	// Start ping-pong for connection health check
	go h.pingConnection(ws)
	
	// Main message handling loop
	for {
		var msg entity.MessageRequest
		err := ws.ReadJSON(&msg)
		
		if err != nil {
			log.Printf("Read error for user %d: %v", userID, err)
			h.cleanupConnection(ws)
			ws.Close()
			break
		}
		
		// Validate the message has required fields
		if msg.SenderID == 0 || msg.ReceiverID == 0 || msg.Content == "" {
			log.Println("Invalid message format:", msg)
			continue
		}
		
		// Set created time if not provided
		if msg.CreatedAt.IsZero() {
			msg.CreatedAt = time.Now()
		}
		
		// Send to broadcast channel for processing
		h.broadcast <- msg
		log.Printf("Message from %d to %d queued for delivery", msg.SenderID, msg.ReceiverID)
	}
}

// Register a new WebSocket connection
func (h *WebSocketHandler) registerConnection(ws *websocket.Conn, userID uint) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	
	// If user already has a connection, close it
	if oldConn, exists := h.connections[userID]; exists {
		log.Printf("User %d already has a connection, closing old one", userID)
		
		// Remove from clients map
		delete(h.clients, oldConn)
		
		// Close the connection
		oldConn.Close()
	}
	
	// Register new connection
	h.clients[ws] = userID
	h.connections[userID] = ws
	
	log.Printf("User %d registered with new connection", userID)
}

// Clean up a closed connection
func (h *WebSocketHandler) cleanupConnection(ws *websocket.Conn) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	
	// Get the user ID for this connection
	userID, exists := h.clients[ws]
	if !exists {
		log.Println("Connection not found in clients map")
		return
	}
	
	// Remove from maps
	delete(h.clients, ws)
	
	// Only delete from connections if this is still the active connection for the user
	if conn, ok := h.connections[userID]; ok && conn == ws {
		delete(h.connections, userID)
		log.Printf("User %d connection cleaned up", userID)
	}
}

// Send periodic pings to keep connection alive
func (h *WebSocketHandler) pingConnection(ws *websocket.Conn) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	
	for {
		select {
		case <-ticker.C:
			if err := ws.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(10*time.Second)); err != nil {
				log.Println("Ping error:", err)
				return
			}
		}
	}
}

func (h *WebSocketHandler) HandleMessages() {
	log.Println("Message handler started and waiting for messages...")
	
	for {
		// Wait for a message from the broadcast channel
		msg := <-h.broadcast
		
		// Save to database
		ctx := context.Background()
		err := h.messageUseCase.SaveMessage(ctx, msg.SenderID, msg.ReceiverID, msg.Content)
		
		if err != nil {
			log.Printf("Failed to save message from %d to %d: %v", msg.SenderID, msg.ReceiverID, err)
			continue
		}
		
		log.Printf("Message from %d to %d saved to database", msg.SenderID, msg.ReceiverID)
		
		// Try to deliver message to recipient if online
		h.mutex.Lock()
		conn, recipientOnline := h.connections[msg.ReceiverID]
		h.mutex.Unlock()
		
		if recipientOnline {
			deliveryErr := conn.WriteJSON(msg)
			
			if deliveryErr != nil {
				log.Printf("Failed to deliver message to user %d: %v", msg.ReceiverID, deliveryErr)
				
				// Clean up bad connection
				h.mutex.Lock()
				delete(h.clients, conn)
				delete(h.connections, msg.ReceiverID)
				h.mutex.Unlock()
				
				conn.Close()
			} else {
				log.Printf("Message delivered to user %d", msg.ReceiverID)
				
				// Send delivery confirmation to sender
				h.mutex.Lock()
				senderConn, senderOnline := h.connections[msg.SenderID]
				h.mutex.Unlock()
				
				if senderOnline {
					confirmation := map[string]interface{}{
						"type":        "delivery_confirmation",
						"message_time": msg.CreatedAt,
						"recipient_id": msg.ReceiverID,
						"delivered_at": time.Now(),
					}
					
					if err := senderConn.WriteJSON(confirmation); err != nil {
						log.Printf("Failed to send delivery confirmation to user %d: %v", msg.SenderID, err)
					}
				}
			}
		} else {
			log.Printf("Recipient %d is offline, message saved to database only", msg.ReceiverID)
		}
	}
}