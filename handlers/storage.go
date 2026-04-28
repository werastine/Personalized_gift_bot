package handlers

import (
	tele "gopkg.in/telebot.v4"
	"sync"
)

type GlobalStorage struct {
	userIDGift map[int64]UserGiftToSend
	userState  map[int64]string
	mu         sync.RWMutex
}

func NewStorage() *GlobalStorage {
	g := &GlobalStorage{
		userIDGift: make(map[int64]UserGiftToSend),
		userState:  make(map[int64]string),
	}
	return g
}

func (a *GlobalStorage) getSession(userID int64) (UserGiftToSend, bool) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	session, ok := a.userIDGift[userID]
	if !ok {
		return session, false
	}
	return session, true
}

func (a *GlobalStorage) setSession(userID int64, session *UserGiftToSend) {
	a.mu.Lock()
	a.userIDGift[userID] = *session
	a.mu.Unlock()
}

func (a *GlobalStorage) setID(giftID string, userID int64) {
	a.mu.Lock()
	defer a.mu.Unlock()
	session := a.userIDGift[userID]
	session.GiftID = giftID
	a.userIDGift[userID] = session
}

func (a *GlobalStorage) DeleteMessage(userID int64, message *tele.Message) {
	a.mu.Lock()
	defer a.mu.Unlock()

	session := a.userIDGift[userID]
	session.MessageIdToDelete = message
	a.userIDGift[userID] = session
}

func (a *GlobalStorage) setState(userID int64, state string) {
	a.mu.Lock()
	defer a.mu.Unlock()

	if state == "" {
		delete(a.userState, userID)
		return
	}
	a.userState[userID] = state
}

func (a *GlobalStorage) getState(userID int64) string {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return a.userState[userID]
}

func (a *GlobalStorage) deleteStorage(userID int64) {
	a.mu.Lock()
	defer a.mu.Unlock()
	delete(a.userState, userID)
}
