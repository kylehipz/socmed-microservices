package main

import (
	"log"

	"github.com/kylehipz/socmed-microservices/users/config"
	"github.com/kylehipz/socmed-microservices/users/internal/db"
	"github.com/kylehipz/socmed-microservices/users/internal/server"
)

func main() {
	dsn := config.Settings.DATABASE_URL

	gormDB := db.NewGormDB(dsn)

	e := server.NewEchoServer(gormDB)

	if err := e.Start(":8080"); err != nil {
		log.Fatal(err)
	}
}
