package main

import (
	"context"

	"github.com/wailsapp/wails/v3/pkg/application"
)

type AppService struct {
	app *application.App
}

func NewAppService(app *application.App) *AppService {
	return &AppService{app: app}
}

func (s *AppService) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
	// Initialisation
	return nil
}
