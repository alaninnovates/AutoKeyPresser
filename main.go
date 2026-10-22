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
	})

	app.RegisterService(application.NewService(NewMyService(app)))

	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:  "AutoKeyPresser",
		Width:  700,
		Height: 500,
	})

	err := app.Run()
	if err != nil {
		panic(err)
	}
}
