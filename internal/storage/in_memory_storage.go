package storage

import (
	"github.com/google/uuid"
	"golang.org/x/net/context"
	"sync"
)

type Storage interface {
	Store(ctx context.Context, key int) (string, error)
	Retrieve(ctx context.Context, key string) (int, bool)
}

type inMemoryStore struct {
	data map[string]int
	mu   sync.RWMutex
}

func NewInMemoryStore() Storage {
	return &inMemoryStore{
		data: make(map[string]int),
	}
}

func (s *inMemoryStore) Store(ctx context.Context, points int) (string, error) {
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
		s.mu.Lock()
		defer s.mu.Unlock()
		id := uuid.New().String()
		s.data[id] = points
		return id, nil
	}
}

func (s *inMemoryStore) Retrieve(ctx context.Context, id string) (int, bool) {
	select {
	case <-ctx.Done():
		return 0, false
	default:
		s.mu.RLock()
		defer s.mu.RUnlock()
		points, exists := s.data[id]
		return points, exists
	}
}
