package store

import (
	"game-social-system/pkg/models"
	"sync"
)

var (
	Mu      sync.Mutex
	Users   = make(map[string]*models.User)
	Parties = make(map[string]*models.Party)
)
