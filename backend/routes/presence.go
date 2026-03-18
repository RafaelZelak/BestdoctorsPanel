package routes

import (
	"log"
	"net/http"
	"sync"
	"time"

	"bestdoctors_service/internal/session"
	"bestdoctors_service/middleware"

	"github.com/gorilla/websocket"
)

var wsUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type UserPresence struct {
	Username string    `json:"username"`
	FullName string    `json:"full_name"`
	Status   string    `json:"status"`
	LastSeen time.Time `json:"last_seen"`
}

type PresenceMessage struct {
	Scope  string `json:"scope"`
	Prompt string `json:"prompt"`
	Status string `json:"status"`
}

type presenceClient struct {
	username string
	send     chan map[string]map[string]UserPresence
}

var presenceHub = struct {
	sync.RWMutex
	// key: "scope:prompt", value: map[username]UserPresence
	state   map[string]map[string]UserPresence
	clients map[*presenceClient]struct{}
}{
	state:   make(map[string]map[string]UserPresence),
	clients: make(map[*presenceClient]struct{}),
}

func presenceSnapshot() map[string]map[string]UserPresence {
	presenceHub.RLock()
	defer presenceHub.RUnlock()

	snapshot := make(map[string]map[string]UserPresence, len(presenceHub.state))
	for promptKey, users := range presenceHub.state {
		usersCopy := make(map[string]UserPresence, len(users))
		for username, presence := range users {
			usersCopy[username] = presence
		}
		snapshot[promptKey] = usersCopy
	}
	return snapshot
}

func broadcastPresence() {
	snapshot := presenceSnapshot()

	presenceHub.RLock()
	defer presenceHub.RUnlock()

	for client := range presenceHub.clients {
		select {
		case client.send <- snapshot:
		default:
		}
	}
}

func removeUserPresence(username string) {
	presenceHub.Lock()
	for promptKey, users := range presenceHub.state {
		delete(users, username)
		if len(users) == 0 {
			delete(presenceHub.state, promptKey)
		}
	}
	presenceHub.Unlock()
}

func PresenceWSHandler(w http.ResponseWriter, r *http.Request) {
	sessionData, ok := r.Context().Value(middleware.SessionDataKey).(*session.SessionData)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	conn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("presence: websocket upgrade failed: %v", err)
		return
	}

	client := &presenceClient{
		username: sessionData.Username,
		send:     make(chan map[string]map[string]UserPresence, 8),
	}

	presenceHub.Lock()
	presenceHub.clients[client] = struct{}{}
	presenceHub.Unlock()

	// Send initial snapshot so the client is up-to-date immediately.
	initialSnapshot := presenceSnapshot()
	if writeErr := conn.WriteJSON(initialSnapshot); writeErr != nil {
		log.Printf("presence: initial write failed for %s: %v", sessionData.Username, writeErr)
	}

	var writeWg sync.WaitGroup
	writeWg.Add(1)
	go func() {
		defer writeWg.Done()
		for snapshot := range client.send {
			if writeErr := conn.WriteJSON(snapshot); writeErr != nil {
				log.Printf("presence: write error for %s: %v", sessionData.Username, writeErr)
				return
			}
		}
	}()

	conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		var msg PresenceMessage
		if readErr := conn.ReadJSON(&msg); readErr != nil {
			break
		}

		conn.SetReadDeadline(time.Now().Add(60 * time.Second))

		if msg.Status == "gone" || msg.Prompt == "" {
			removeUserPresence(sessionData.Username)
			broadcastPresence()
			continue
		}

		promptKey := msg.Scope + ":" + msg.Prompt

		presenceHub.Lock()
		if _, exists := presenceHub.state[promptKey]; !exists {
			presenceHub.state[promptKey] = make(map[string]UserPresence)
		}

		// Remove previous prompt entry for this user before setting the new one.
		for existingKey, users := range presenceHub.state {
			if existingKey == promptKey {
				continue
			}
			if _, hadEntry := users[sessionData.Username]; hadEntry {
				delete(users, sessionData.Username)
				if len(users) == 0 {
					delete(presenceHub.state, existingKey)
				}
			}
		}

		presenceHub.state[promptKey][sessionData.Username] = UserPresence{
			Username: sessionData.Username,
			FullName: sessionData.Username,
			Status:   msg.Status,
			LastSeen: time.Now(),
		}
		presenceHub.Unlock()

		broadcastPresence()
	}

	presenceHub.Lock()
	delete(presenceHub.clients, client)
	close(client.send)
	presenceHub.Unlock()

	writeWg.Wait()
	conn.Close()

	removeUserPresence(sessionData.Username)
	broadcastPresence()
}

func StartPresenceEviction() {
	evictionTicker := time.NewTicker(10 * time.Second)
	defer evictionTicker.Stop()

	for range evictionTicker.C {
		staleThreshold := time.Now().Add(-15 * time.Second)
		changed := false

		presenceHub.Lock()
		for promptKey, users := range presenceHub.state {
			for username, presence := range users {
				if presence.LastSeen.Before(staleThreshold) {
					delete(users, username)
					changed = true
				}
			}
			if len(users) == 0 {
				delete(presenceHub.state, promptKey)
			}
		}
		presenceHub.Unlock()

		if changed {
			broadcastPresence()
		}
	}
}
