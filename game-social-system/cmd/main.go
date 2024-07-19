package main

import (
	"game-social-system/pkg/api"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.POST("/add_friend", api.AddFriend)
	r.POST("/accept_friend_request", api.AcceptFriendRequest)
	r.POST("/reject_friend_request", api.RejectFriendRequest)
	r.POST("/remove_friend", api.RemoveFriend)
	r.GET("/view_friends/:user_id", api.ViewFriends)

	r.POST("/create_party", api.CreateParty)
	r.POST("/invite_to_party", api.InviteToParty)
	r.POST("/join_party", api.JoinParty)
	r.POST("/leave_party", api.LeaveParty)
	r.POST("/remove_from_party", api.RemoveUserFromParty)

	r.GET("/ws/:user_id", api.WSHandler)

	r.Run(":8080")
}
