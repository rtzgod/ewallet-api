package app

import (
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	htp "github.com/rtzgod/EWallet/internal/delivery/http"
	"github.com/rtzgod/EWallet/internal/domain/service"
	"github.com/rtzgod/EWallet/internal/repository"
	"github.com/rtzgod/EWallet/internal/repository/psql"
	"log"
	"net/http"
	"os"
)

func Run() {
	r := mux.NewRouter()
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	port := os.Getenv("SERVER_PORT")
	db := psql.Connect()

	storage := repository.NewStorage(db)

	services := service.NewService(storage)

	handler := htp.NewHandler(services)

	handler.InitRoutes(r)

	done := make(chan bool)
	go http.ListenAndServe(":"+port, r)
	log.Println("Server started on port: " + port)
	<-done
}
