package api

import (
	"game-social-system/pkg/models"
	"game-social-system/pkg/store"
	"game-social-system/pkg/utils" // Import utils package
	"net/http"

	"github.com/gin-gonic/gin"
)

// AddFriend handles adding a new friend
func AddFriend(c *gin.Context) {
	var req models.FriendRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	store.Mu.Lock()
	defer store.Mu.Unlock()

	fromUser, exists := store.Users[req.From]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	toUser, exists := store.Users[req.To]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Friend not found"})
		return
	}

	if utils.Contains(fromUser.Friends, req.To) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Already friends"})
		return
	}

	fromUser.Friends = append(fromUser.Friends, req.To)
	toUser.Friends = append(toUser.Friends, req.From)
	store.Users[req.From] = fromUser
	store.Users[req.To] = toUser

	c.JSON(http.StatusOK, gin.H{"message": "Friend added"})
}

// RemoveFriend handles removing a friend
func RemoveFriend(c *gin.Context) {
	var req models.FriendRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	store.Mu.Lock()
	defer store.Mu.Unlock()

	fromUser, exists := store.Users[req.From]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	toUser, exists := store.Users[req.To]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Friend not found"})
		return
	}

	fromUser.Friends = utils.Remove(fromUser.Friends, req.To)
	toUser.Friends = utils.Remove(toUser.Friends, req.From)
	store.Users[req.From] = fromUser
	store.Users[req.To] = toUser

	c.JSON(http.StatusOK, gin.H{"message": "Friend removed"})
}

// ViewFriends handles viewing the list of friends
func ViewFriends(c *gin.Context) {
	userID := c.Param("userID")

	store.Mu.Lock()
	defer store.Mu.Unlock()

	user, exists := store.Users[userID]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"friends": user.Friends})
}

// AcceptFriendRequest handles accepting a friend request
func AcceptFriendRequest(c *gin.Context) {
	var req models.FriendRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	store.Mu.Lock()
	defer store.Mu.Unlock()

	fromUser, exists := store.Users[req.From]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	toUser, exists := store.Users[req.To]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Friend not found"})
		return
	}

	if !utils.Contains(fromUser.Friends, req.To) {
		fromUser.Friends = append(fromUser.Friends, req.To)
		store.Users[req.From] = fromUser
	}

	if !utils.Contains(toUser.Friends, req.From) {
		toUser.Friends = append(toUser.Friends, req.From)
		store.Users[req.To] = toUser
	}

	c.JSON(http.StatusOK, gin.H{"message": "Friend request accepted"})
}

// RejectFriendRequest handles rejecting a friend request
func RejectFriendRequest(c *gin.Context) {
	var req models.FriendRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Implement rejection logic if necessary

	c.JSON(http.StatusOK, gin.H{"message": "Friend request rejected"})
}
