package models

// User represents a user in the system
type User struct {
	ID      string   `json:"id"`
	Friends []string `json:"friends"`
}

// Party represents a game party
type Party struct {
	ID      string   `json:"id"`
	Leader  string   `json:"leader"`
	Members []string `json:"members"`
}

// FriendRequest represents a friend request
type FriendRequest struct {
	From string `json:"from"`
	To   string `json:"to"`
}

// CreatePartyRequest represents a request to create a party
type CreatePartyRequest struct {
	Creator string `json:"creator"`
}

// InviteRequest represents a request to invite a user to a party
type InviteRequest struct {
	PartyID string `json:"party_id"`
	Inviter string `json:"inviter"`
	Invitee string `json:"invitee"`
}

// JoinPartyRequest represents a request for a user to join a party
type JoinPartyRequest struct {
	PartyID string `json:"party_id"`
	UserID  string `json:"user_id"`
}

// LeavePartyRequest represents a request for a user to leave a party
type LeavePartyRequest struct {
	PartyID string `json:"party_id"`
	UserID  string `json:"user_id"`
}

// AcceptPartyInvitationRequest represents a request to accept a party invitation
type AcceptPartyInvitationRequest struct {
	PartyID string `json:"party_id"`
	UserID  string `json:"user_id"`
}

// RejectPartyInvitationRequest represents a request to reject a party invitation
type RejectPartyInvitationRequest struct {
	PartyID string `json:"party_id"`
	UserID  string `json:"user_id"`
}

// RemoveUserFromPartyRequest represents a request to remove a user from a party
type RemoveUserFromPartyRequest struct {
	PartyID   string `json:"party_id"`
	UserID    string `json:"user_id"`
	Requester string `json:"requester"`
}

// RemoveFriendRequest represents a request to remove a friend
type RemoveFriendRequest struct {
	From string `json:"from"`
	To   string `json:"to"`
}
