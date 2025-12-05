package storage

import (
	"fmt"
	"sync"
	"sync/atomic"
)

const (
	StatusAvalible    = "avalible"
	StatusNotAvalible = "not avalible"
)

type Storage struct {
	counter atomic.Int64
	links   map[int64]map[string]string
	mu      sync.RWMutex
}

func MakeStorage() *Storage {
	s := &Storage{
		links: make(map[int64]map[string]string),
	}
	return s
}

func (s *Storage) RegistrationLinks(links map[string]string) int64 {
	if links == nil {
		return -1
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.counter.Add(1)
	s.links[s.counter.Load()] = links
	return s.counter.Load()
}

// error check in handler
func (s *Storage) GetLiinksInfo(id int64) map[string]string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	links, ok := s.links[id]
	if !ok {
		return nil
	}
	return links
}

func (s *Storage) GetBucketsInfo(IDs ...int64) (map[int64]map[string]string, error) {
	res := make(map[int64]map[string]string)

	for _, id := range IDs {
		info := s.GetLiinksInfo(id)
		if info == nil {
			return nil, fmt.Errorf("invalid id - %d", id)
		}
		res[id] = info
	}
	return res, nil
}
