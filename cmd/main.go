package main

import (
	"log"

	"forum/internal/delivery"
	"forum/internal/repository"
	"forum/internal/server"
	"forum/internal/service"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := repository.OpenSqliteDB("store.db")
	if err != nil {
		log.Fatalf("error while opening db: %s", err)
	}
	repo := repository.NewRepository(db)
	service := service.NewService(repo)
	handler := delivery.NewHandler(service)
	server := new(server.Server)
	if err := server.Run("8080", handler.InitRoutes()); err != nil {
		log.Fatalf("error while running the server: %s", err.Error())
	}
}
