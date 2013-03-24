package goku

import (
    "errors"
    // "net/http"
    "time"
)

// session store
var sessionItems map[string]map[string]sessionItem = make(map[string]map[string]sessionItem)

// session expired in sec.
var maxlifetime int64 = 1800 // 30min

type memoryProvider struct {
}

func (m *memoryProvider) StartSession(sid string) (Session, error) {
    if sid == "" {
        return nil, errors.New("session id can't be null")
    }
    s := &memorySession{}
    s.sessionId = sid
    if items, ok := sessionItems[sid]; ok {
        s.items = items
    }
    return s, nil
}

type sessionItem struct {
    timeAccessed time.Time //last accessed time     
    value        interface{}
}

type memorySession struct {
    sessionId string
    items     map[string]sessionItem
}

func (s *memorySession) Set(key string, value interface{}) error {
    if s.items == nil {
        // TODO: lock ?
        s.items = make(map[string]sessionItem)
        sessionItems[s.sessionId] = s.items
    }
    item := sessionItem{}
    item.timeAccessed = time.Now()
    item.value = value
    s.items[key] = item
    return nil
}

func (s *memorySession) Get(key string) interface{} {
    if s.items == nil {
        return nil
    }
    v, ok := s.items[key]
    if ok {
        if time.Now().Unix()-v.timeAccessed.Unix() < maxlifetime {
            return v.value
        } else {
            delete(s.items, key)
        }
    }
    return nil
}

func (s *memorySession) Del(key string) error {
    if s.items == nil {
        return nil
    }
    delete(s.items, key)
    return nil
}

func (s *memorySession) Clear() error {
    if s.items == nil {
        return nil
    }
    delete(sessionItems, s.sessionId)
    return nil
}

func (s *memorySession) ID() string {
    return s.sessionId
}

func init() {
    RegisterSessionProvider("memory", &memoryProvider{})
}
