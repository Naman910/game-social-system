package api

import (
	"game-social-system/pkg/models"
	"game-social-system/pkg/store"
	"game-social-system/pkg/utils" // Import utils package
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateParty handles the creation of a new party
func CreateParty(c *gin.Context) {
	var req models.CreatePartyRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Ensure the creator exists
	if _, exists := store.Users[req.Creator]; !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Creator not found"})
		return
	}

	partyID := utils.GenerateID()
	newParty := &models.Party{ // Create a pointer to the Party struct
		ID:      partyID,
		Leader:  req.Creator,
		Members: []string{req.Creator},
	}

	store.Mu.Lock()
	store.Parties[partyID] = newParty
	store.Mu.Unlock()

	// Notify all party members about the newly created party
	notifyPartyStatus(partyID)

	c.JSON(http.StatusOK, gin.H{"party_id": partyID})
}

// InviteToParty handles inviting a user to a party
func InviteToParty(c *gin.Context) {
	var req models.InviteRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	store.Mu.Lock()
	defer store.Mu.Unlock()

	party, exists := store.Parties[req.PartyID]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Party not found"})
		return
	}

	// Ensure the inviter and invitee exist
	_, exists = store.Users[req.Inviter]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Inviter not found"})
		return
	}

	_, exists = store.Users[req.Invitee]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Invitee not found"})
		return
	}

	// Check if inviter is the leader of the party
	if party.Leader != req.Inviter {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only party leader can invite users"})
		return
	}

	// Add invitee to party members if not already present
	if !utils.Contains(party.Members, req.Invitee) {
		party.Members = append(party.Members, req.Invitee)
		store.Parties[req.PartyID] = party
	}

	// Notify all party members about the updated party status
	notifyPartyStatus(req.PartyID)

	c.JSON(http.StatusOK, gin.H{"message": "Invitation sent"})
}

// JoinParty handles a user joining a party
func JoinParty(c *gin.Context) {
	var req models.JoinPartyRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	store.Mu.Lock()
	defer store.Mu.Unlock()

	party, exists := store.Parties[req.PartyID]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Party not found"})
		return
	}

	// Ensure the user exists
	if _, exists := store.Users[req.UserID]; !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Add user to party members if not already present
	if !utils.Contains(party.Members, req.UserID) {
		JoinParty(c)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Joined the party"})
}

// LeaveParty handles a user leaving a party
func LeaveParty(c *gin.Context) {
	var req models.LeavePartyRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	store.Mu.Lock()
	defer store.Mu.Unlock()

	party, exists := store.Parties[req.PartyID]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Party not found"})
		return
	}

	if !utils.Contains(party.Members, req.UserID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User is not in the party"})
		return
	}

	// Remove user from party members
	party.Members = utils.Remove(party.Members, req.UserID)
	store.Parties[req.PartyID] = party

	c.JSON(http.StatusOK, gin.H{"message": "Left the party"})
}

// AcceptPartyInvitation handles accepting a party invitation
func AcceptPartyInvitation(c *gin.Context) {
	var req models.AcceptPartyInvitationRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	store.Mu.Lock()
	defer store.Mu.Unlock()

	party, exists := store.Parties[req.PartyID]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Party not found"})
		return
	}

	// Ensure the user exists
	if _, exists := store.Users[req.UserID]; !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Add user to party members if not already present
	if !utils.Contains(party.Members, req.UserID) {
		party.Members = append(party.Members, req.UserID)
		store.Parties[req.PartyID] = party
	}

	c.JSON(http.StatusOK, gin.H{"message": "Party invitation accepted"})
}

// RejectPartyInvitation handles rejecting a party invitation
func RejectPartyInvitation(c *gin.Context) {
	var req models.RejectPartyInvitationRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Implement rejection logic if necessary

	c.JSON(http.StatusOK, gin.H{"message": "Party invitation rejected"})
}

// RemoveUserFromParty handles removing a user from a party
func RemoveUserFromParty(c *gin.Context) {
	var req models.RemoveUserFromPartyRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	store.Mu.Lock()
	defer store.Mu.Unlock()

	party, exists := store.Parties[req.PartyID]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Party not found"})
		return
	}

	if party.Leader != req.Requester {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only party leader can remove users"})
		return
	}

	if !utils.Contains(party.Members, req.UserID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User is not in the party"})
		return
	}

	party.Members = utils.Remove(party.Members, req.UserID)
	store.Parties[req.PartyID] = party

	c.JSON(http.StatusOK, gin.H{"message": "User removed from party"})
}
