package main

import (
	"fmt"
	"forum/internal/delivery"
	"forum/internal/repository"
	"forum/internal/server"
	"forum/internal/service"
	"log"

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

	fmt.Printf("Starting server at port 8081\nhttp://localhost:8081/\n")

	if err := server.Run("8081", handler.InitRoutes()); err != nil {
		log.Fatalf("error while running the server: %s", err.Error())
	}
}
