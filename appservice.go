package main

import (
	"context"
	"fmt"
	"time"

	hook "github.com/robotn/gohook"
	"github.com/wailsapp/wails/v3/pkg/application"
)

type ActiveTab string

const (
	ActiveTabKeyPresser ActiveTab = "key_presser"
	ActiveTabSlime      ActiveTab = "slime"
	ActiveTabFisher     ActiveTab = "fisher"
)

type AppStatus string

const (
	AppStatusIdle    AppStatus = "idle"
	AppStatusRunning AppStatus = "running"
)

type AppService struct {
	app               *application.App
	activeTab         ActiveTab
	appStatus         AppStatus
	keyPresserService *KeyPresserService
}

func NewAppService(app *application.App, keyPresserService *KeyPresserService) *AppService {
	return &AppService{
		app:               app,
		activeTab:         ActiveTabKeyPresser,
		keyPresserService: keyPresserService,
	}
}

func (s *AppService) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
	hook.Register(hook.KeyDown, []string{"f1"}, func(e hook.Event) {
		if s.appStatus == AppStatusRunning {
			return
		}
		fmt.Println("f1 start")
		time.Sleep(100 * time.Millisecond)
		s.appStatus = AppStatusRunning
		switch s.activeTab {
		case ActiveTabKeyPresser:
			s.keyPresserService.start()
		case ActiveTabSlime:
			// call slime service start
		case ActiveTabFisher:
			// call fisher service start
		}
	})
	hook.Register(hook.KeyDown, []string{"f3"}, func(e hook.Event) {
		if s.appStatus == AppStatusIdle {
			return
		}
		fmt.Println("f3 stop")
		s.appStatus = AppStatusIdle
		s.keyPresserService.stop()
	})
	go func() {
		h := hook.Start()
		<-hook.Process(h)
	}()
	return nil
}

func (s *AppService) SetActiveTab(tab ActiveTab) ActiveTab {
	s.activeTab = tab
	return s.activeTab
}
