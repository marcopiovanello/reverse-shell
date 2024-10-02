package internal

import (
	"bytes"
	"io"
	"sync"
)

type ClientState struct {
	store map[string][]byte
	mu    sync.RWMutex
}

func NewClientState() *ClientState {
	return &ClientState{
		store: make(map[string][]byte),
	}
}

func (s *ClientState) Store(key string, value []byte) {
	s.mu.Lock()
	s.store[key] = value
	s.mu.Unlock()
}

func (s *ClientState) Get(key string) []byte {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.store[key]
}

func (s *ClientState) Read(key string) io.Reader {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return bytes.NewReader(s.store[key])
}
