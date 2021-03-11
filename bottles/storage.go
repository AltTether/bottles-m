package bottles

import (
	"fmt"
	"sync"
)

type Storage struct {
	messages []*Message
	mux      *sync.Mutex
}

func NewStorage() *Storage {
	return &Storage{
		mux: &sync.Mutex{},
	}
}

func (s *Storage) Get() (*Message, error) {
	s.mux.Lock()
	if len(s.messages) == 0 {
		s.mux.Unlock()
		return nil, fmt.Errorf("No Messages")
	}

	m := s.messages[0]
	s.messages = s.messages[1:]
	s.mux.Unlock()
	return m, nil
}

func (s *Storage) Add(m *Message) error {
	s.mux.Lock()
	s.messages = append(s.messages, m)
	s.mux.Unlock()
	return nil
}
