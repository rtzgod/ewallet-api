package main

import (
	_ "github.com/rtzgod/EWallet/docs"
	"github.com/rtzgod/EWallet/internal/app"
)

//	@title			Ewallet App API
//	@version		2.0
//	@description	API Server for Ewallet App

//	@host		localhost:8080
//	@BasePath /

func main() {
	app.Run()
}
