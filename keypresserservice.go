package main

import (
	"sync"
)

type ActionType string

const (
	ActionTypeKeyboard ActionType = "keyboard"
	ActionTypeMouse    ActionType = "mouse"
)

type Action struct {
	ID             int        `json:"id"`
	Type           ActionType `json:"type"`
	Key            string     `json:"key,omitempty"`
	Button         string     `json:"button,omitempty"`
	Hold           bool       `json:"hold"`
	Repeat         bool       `json:"repeat"`
	RepeatInterval int        `json:"repeatInterval"`
}

type KeyPresserService struct {
	actions []Action
	nextID  int
	mu      sync.RWMutex
}

func NewKeyPresserService() *KeyPresserService {
	return &KeyPresserService{
		actions: []Action{},
		nextID:  1,
	}
}

func (s *KeyPresserService) GetAll() []Action {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.actions
}

func (s *KeyPresserService) AddAction(action Action) (*Action, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	action.ID = s.nextID
	s.nextID++

	s.actions = append(s.actions, action)

	result := action
	return &result, nil
}

func (s *KeyPresserService) AddEmptyAction(actionType ActionType) (*Action, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	action := Action{
		ID:             s.nextID,
		Type:           actionType,
		Key:            "",
		Button:         "",
		Hold:           false,
		Repeat:         false,
		RepeatInterval: 0,
	}

	s.nextID++
	s.actions = append(s.actions, action)

	result := action
	return &result, nil
}

func (s *KeyPresserService) UpdateAction(action Action) (*Action, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, a := range s.actions {
		if a.ID == action.ID {
			s.actions[i] = action
			return &action, nil
		}
	}

	return nil, nil
}

func (s *KeyPresserService) DeleteAction(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, a := range s.actions {
		if a.ID == id {
			s.actions = append(s.actions[:i], s.actions[i+1:]...)
			return nil
		}
	}

	return nil
}
