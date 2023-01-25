package main

import (
	"fmt"
	"log"

	"forum/internal/delivery"
	"forum/internal/repository"
	"forum/internal/server"
	"forum/internal/service"

	_ "github.com/mattn/go-sqlite3"
)

const port = "8080"

func main() {
	db, err := repository.OpenSqliteDB("store.db")
	if err != nil {
		log.Fatalf("error while opening db: %s", err)
	}

	repo := repository.NewRepository(db)
	service := service.NewService(repo)
	handler := delivery.NewHandler(service)
	server := new(server.Server)

	fmt.Printf("Starting server at port %s\nhttp://localhost:%s/\n", port, port)

	if err := server.Run(port, handler.InitRoutes()); err != nil {
		log.Fatalf("error while running the server: %s", err.Error())
	}
}
