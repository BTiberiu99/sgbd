package app

import (
	"sgbd4/go/expose"

	"github.com/leaanthony/mewn"
	"github.com/wailsapp/wails"
)

var APP *wails.App

func LoadApp() {

	if APP == nil {

		js := mewn.String("./frontend/dist/app.js")
		css := mewn.String("./frontend/dist/app.css")

		app := wails.CreateApp(&wails.AppConfig{
			Width:     1024,
			Height:    768,
			Title:     "sgbd",
			JS:        js,
			CSS:       css,
			Colour:    "#131313",
			Resizable: true,
		})

		//Bind Exposed Functions
		app.Bind(expose.CreateConnection)
		app.Bind(expose.SwitchConnection)
		app.Bind(expose.RemoveConnection)
		app.Bind(expose.GetConnections)
		app.Bind(expose.AddNotNull)

		//Run SQL
		app.Bind(expose.Run)
		app.Bind(expose.GetTables)

		APP = app
		app.Run()
	}
}
