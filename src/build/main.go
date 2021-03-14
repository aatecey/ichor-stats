package main

import (
	_application "ichor-stats/src/app/application/implementation"
)

func main() {
	app := _application.NewApplication()
	app.Run()
}