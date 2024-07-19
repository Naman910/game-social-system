package api

import (
	"game-social-system/pkg/store"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	connections   = make(map[string]*websocket.Conn)
	connectionsMu sync.Mutex
)

func WSHandler(c *gin.Context) {
	userID := c.Param("user_id")
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	connectionsMu.Lock()
	connections[userID] = conn
	connectionsMu.Unlock()

	notifyFriendsStatus(userID, "online")

	defer func() {
		connectionsMu.Lock()
		delete(connections, userID)
		connectionsMu.Unlock()
		notifyFriendsStatus(userID, "offline")
		conn.Close()
	}()

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}
}

func notifyFriendsStatus(userID, status string) {
	store.Mu.Lock()
	user, exists := store.Users[userID]
	store.Mu.Unlock()

	if !exists {
		return
	}

	for _, friendID := range user.Friends {
		connectionsMu.Lock()
		friendConn, friendOnline := connections[friendID]
		connectionsMu.Unlock()

		if friendOnline {
			message := map[string]string{
				"user_id": userID,
				"status":  status,
			}
			friendConn.WriteJSON(message)
		}
	}
}

func notifyPartyStatus(partyID string) {
	store.Mu.Lock()
	party, exists := store.Parties[partyID]
	store.Mu.Unlock()

	if !exists {
		return
	}

	for _, memberID := range party.Members {
		connectionsMu.Lock()
		memberConn, memberOnline := connections[memberID]
		connectionsMu.Unlock()

		if memberOnline {
			message := map[string]interface{}{
				"party_id": partyID,
				"members":  party.Members,
			}
			memberConn.WriteJSON(message)
		}
	}
}
