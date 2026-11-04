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
	slimeService      *SlimeService
	fisherService     *FisherService
}

func NewAppService(app *application.App, keyPresserService *KeyPresserService, slimeService *SlimeService, fisherService *FisherService) *AppService {
	return &AppService{
		app:               app,
		activeTab:         ActiveTabKeyPresser,
		keyPresserService: keyPresserService,
		slimeService:      slimeService,
		fisherService:     fisherService,
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
			break
		case ActiveTabSlime:
			s.slimeService.start()
			break
		case ActiveTabFisher:
			s.fisherService.start()
			break
		}
	})
	hook.Register(hook.KeyDown, []string{"f3"}, func(e hook.Event) {
		if s.appStatus == AppStatusIdle {
			return
		}
		fmt.Println("f3 stop")
		s.appStatus = AppStatusIdle
		switch s.activeTab {
		case ActiveTabKeyPresser:
			s.keyPresserService.stop()
			break
		case ActiveTabSlime:
			s.slimeService.stop()
			break
		case ActiveTabFisher:
			s.fisherService.stop()
			break
		}
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
