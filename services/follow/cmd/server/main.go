package main

import (
	"log"

	"github.com/kylehipz/socmed-microservices/follow/config"
	"github.com/kylehipz/socmed-microservices/follow/internal/db"
	"github.com/kylehipz/socmed-microservices/follow/internal/server"
)

func main() {
	dsn := config.Settings.DATABASE_URL

	gormDB := db.NewGormDB(dsn)

	e := server.NewEchoServer(gormDB)

	log.Println("Starting follow service...")
	if err := e.Start(":8080"); err != nil {
		log.Fatal(err)
	}
}
