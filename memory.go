package sessions

import (
	"sync"
)

type memorySession struct {
	value  interface{}
	rwLock sync.RWMutex
	valueType
}

type cookieValueMap map[string]interface{}

var (
	_ Session = &memorySession{}
)

func newMemorySession(valueType valueType) Session {
	var value interface{}
	if valueType == ValueMap {
		value = make(cookieValueMap, 10)
	}
	return &memorySession{
		value:     value,
		valueType: valueType,
	}
}

func (m *memorySession) Get(keys ...string) interface{} {

	m.rwLock.RLock()
	defer m.rwLock.RUnlock()

	if m.value == nil {
		return m.value
	}

	if value, ok := m.value.(cookieValueMap); ok {

		return value[keys[0]]
	}

	return m.value
}

func (m *memorySession) Set(value interface{}, keys ...string) {
	m.rwLock.Lock()
	defer m.rwLock.Unlock()

	if m.valueType == ValueMap && len(keys) == 0 {
		panic("The type is map. Please pass a key")
	} else if m.valueType == ValueString && len(keys) > 0 {
		panic("The type is string. Please don't pass the key")
	}

	if v, ok := m.value.(cookieValueMap); ok {
		v[keys[0]] = value
	}

	m.value = value

}

func (m *memorySession) Delete(keys ...string) {
	m.rwLock.Lock()
	defer m.rwLock.Unlock()

	if m.value == nil {
		return
	}

	if value, ok := m.value.(cookieValueMap); ok && len(keys) > 0 {
		for _, v := range keys {
			delete(value, v)
		}
		return
	}

	var value interface{}
	if m.valueType == ValueMap {
		value = make(cookieValueMap, 10)
	}

	m.value = value

}
