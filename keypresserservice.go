package main

import (
	"context"
	"sync"

	"github.com/go-vgo/robotgo"
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
	ctx     context.Context
	cancel  context.CancelFunc
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

func (s *KeyPresserService) start() {
	ctx, cancel := context.WithCancel(context.Background())
	s.ctx = ctx
	s.cancel = cancel

	go func() {
		select {
		case <-s.ctx.Done():
			return
		default:
			s.mu.RLock()
			for _, action := range s.actions {
				if action.Repeat {
					if action.Type == ActionTypeKeyboard {
						err := robotgo.KeyTap(action.Key)
						if err != nil {
							return
						}
					} else if action.Type == ActionTypeMouse {
						err := robotgo.Click(action.Button)
						if err != nil {
							return
						}
					}
				}
			}
			s.mu.RUnlock()
		}
	}()
}

func (s *KeyPresserService) stop() {
	s.cancel()
}
