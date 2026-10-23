package main

import (
	"embed"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed frontend/dist
var assets embed.FS

func main() {
	app := application.New(application.Options{
		Name:     "AutoKeyPresser",
		Services: []application.Service{},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	keyPresserService := NewKeyPresserService()
	appService := NewAppService(app, keyPresserService)

	app.RegisterService(application.NewService(keyPresserService))
	app.RegisterService(application.NewService(appService))

	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:  "AutoKeyPresser",
		Width:  700,
		Height: 500,
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 50,
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarDefault,
		},
		BackgroundColour: application.NewRGB(6, 7, 15),
		URL:              "/",
	})

	err := app.Run()
	if err != nil {
		panic(err)
	}
}
