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
		html := mewn.String("./frontend/dist/index.html")

		app := wails.CreateApp(&wails.AppConfig{
			Width:     1024,
			Height:    768,
			Title:     "Proiect SGBD",
			JS:        js,
			CSS:       css,
			HTML:      html,
			Colour:    "#131313",
			Resizable: true,
		})

		//Bind Exposed Functions
		app.Bind(expose.CreateConnection)
		app.Bind(expose.SwitchConnection)
		app.Bind(expose.RemoveConnection)
		app.Bind(expose.GetConnections)
		app.Bind(expose.AddNotNull)
		app.Bind(expose.FixPrimaryKey)
		app.Bind(expose.AddPrimaryKey)
		//Run SQL
		app.Bind(expose.Run)
		app.Bind(expose.GetTables)
		app.Bind(expose.ResetTables)

		APP = app

		app.Run()
	}
}
