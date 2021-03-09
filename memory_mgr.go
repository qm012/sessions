package sessions

import (
	"sync"
)

type memoryMgr struct {
	Data   map[string]Session
	RwLock sync.RWMutex
}

func newMemoryMgr() SessionMgr {
	return &memoryMgr{
		Data: make(map[string]Session, 100000),
	}
}

var _ SessionMgr = &memoryMgr{}

func (m *memoryMgr) GetSession(cookValue string) Session {
	m.RwLock.RLock()
	defer m.RwLock.RUnlock()
	session, _ := m.Data[cookValue]
	return session
}

func (m *memoryMgr) CreateSession(cookValue string, valueType valueType, expire int) Session {
	m.RwLock.Lock()

	defer m.RwLock.Unlock()
	session, ok := m.Data[cookValue]
	if !ok {
		session = newMemorySession(valueType)
	}
	m.Data[cookValue] = session
	return session
}
