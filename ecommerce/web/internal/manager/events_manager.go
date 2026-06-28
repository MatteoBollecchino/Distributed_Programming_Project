package manager

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

// EventsManager manages synchronization between server and browser
type EventsManager struct {
	mu      sync.Mutex
	clients map[chan struct{}]bool
}

// NewEventsManager creates a new manager
func NewEventsManager() *EventsManager {
	return &EventsManager{
		clients: make(map[chan struct{}]bool),
	}
}

// HandleEvents is the HTTP Handler for the events, that the browser will call
func (em *EventsManager) HandleEvents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming not supported", http.StatusInternalServerError)
		return
	}

	// Creation of the channel for a specific client
	messageChan := make(chan struct{}, 1)

	// Register the client in mutual exclusion
	em.mu.Lock()
	em.clients[messageChan] = true
	em.mu.Unlock()

	// Cleaning when client disconnects
	defer func() {
		em.mu.Lock()
		delete(em.clients, messageChan)
		em.mu.Unlock()
		close(messageChan)
	}()

	// Invia un commento di keep-alive iniziale per Firefox
	fmt.Fprintf(w, ": ok\n\n")
	flusher.Flush()

	for {
		select {
		case <-messageChan:
			fmt.Fprintf(w, "data: reload\n\n")
			flusher.Flush()
		case <-r.Context().Done():
			return
		}
	}
}

// NotifyCatalogUpdate tells the connected clients to refresh
func (em *EventsManager) NotifyCatalogUpdate() {
	em.mu.Lock()
	defer em.mu.Unlock()

	log.Printf("Notification of refreshing sent to %d connected clients", len(em.clients))
	for messageChan := range em.clients {
		select {
		case messageChan <- struct{}{}:
		default: // The channel is full -> do nothing and skip to next client
		}
	}
}
