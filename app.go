package main

import (
	"context"

	"github.com/wailsapp/wails/v3/pkg/application"
)

type MouseAction struct {
	Type           string `json:"type"`
	Button         string `json:"button"`
	Hold           bool   `json:"hold"`
	Repeat         bool   `json:"repeat"`
	RepeatInterval int    `json:"repeatInterval"`
}

type KeyboardAction struct {
	Type           string `json:"type"`
	Key            string `json:"key"`
	Hold           bool   `json:"hold"`
	Repeat         bool   `json:"repeat"`
	RepeatInterval int    `json:"repeatInterval"`
}
type MyService struct {
	app *application.App
}

func NewMyService(app *application.App) *MyService {
	return &MyService{app: app}
}

func (s *MyService) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
	// Initialisation
	return nil
}

func (s *MyService) Greet(name string) string {
	return "Hello " + name
}
