package server

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sarkarshuvojit/kafka-sync-proxy/internal/controller"
)

func setupRoutes(ctx context.Context, app *fiber.App) {
	controller.AppContext = ctx
	api := app.Group("/v1")
	api.Post("/", controller.HandleRest)
}

func Start() {
	app := fiber.New(fiber.Config{
		AppName: "Kafka Sync Proxy",
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGABRT, os.Interrupt)
	go func() {
		<-sigchan
		log.Println("Received shutdown signal, starting graceful shutdown...")
		cancel()

		// Give ongoing requests time to complete
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer shutdownCancel()

		if err := app.Server().ShutdownWithContext(shutdownCtx); err != nil {
			log.Printf("Error during graceful shutdown: %v", err)
		} else {
			log.Println("Server shutdown gracefully")
		}
	}()

	setupRoutes(ctx, app)
	app.Listen(":8420")
}
