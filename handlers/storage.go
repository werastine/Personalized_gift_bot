package handlers

import (
	"sync"
)

type GlobalStorage struct {
	userIDGift map[int64]UserGiftToSend
	userState  map[int64]string
	mu         sync.RWMutex
}

func (a *GlobalStorage) getSession(userID int64) UserGiftToSend {
	a.mu.RLock()
	defer a.mu.RUnlock()
	session := a.userIDGift[userID]
	return session
}

func (a *GlobalStorage) setID(giftID string, userID int64) {
	a.mu.Lock()
	defer a.mu.Unlock()
	session := a.userIDGift[userID]
	session.GiftID = giftID
	a.userIDGift[userID] = session
}

func NewStorage() *GlobalStorage {
	g := &GlobalStorage{
		userIDGift: make(map[int64]UserGiftToSend),
		userState:  make(map[int64]string),
	}
	return g
}
