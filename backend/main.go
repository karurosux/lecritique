package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	_ "lecritique/docs"
	"lecritique/internal/shared/config"
	"lecritique/internal/shared/database"
	"lecritique/internal/shared/server"
)

// @title LeCritique API
// @version 1.0
// @description Organization feedback management system API
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email justdevelopitmx@proton.me

// @host localhost:8080

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	db, err := database.Initialize(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	srv := server.NewWithDI(cfg, db)

	go func() {
		if err := srv.Start(); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
}

