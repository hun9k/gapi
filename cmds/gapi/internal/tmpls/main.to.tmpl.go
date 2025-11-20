package tmpls

const Main = `package main

import (
	_ "{{.Path}}/routers"

	"github.com/hun9k/gapi"
)

func main() {
	// create and get app
	app := gapi.App()

	// run app
	if err := app.Run(); err != nil {
		gapi.Log().Error(err.Error())
	}
}

`
