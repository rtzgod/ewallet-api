package main

import (
	_ "github.com/rtzgod/ewallet-api/docs"
	"github.com/rtzgod/ewallet-api/internal/app"
)

//	@title			Ewallet App API
//	@version		2.0
//	@description	API Server for Ewallet App

//	@host		localhost:8080
//	@BasePath /

func main() {
	app.Run()
}
