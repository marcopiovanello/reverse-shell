package internal

import (
	"os"
	"sync"
)

type State struct {
	wd string
	mu sync.RWMutex
}

func NewState() *State {
	cwd, _ := os.Getwd()
	return &State{
		wd: cwd,
	}
}

func (s *State) GetCurrentDirectory() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.wd
}

func (s *State) SetCurrentDirectory(d string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.wd = d
}
