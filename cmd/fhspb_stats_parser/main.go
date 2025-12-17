package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Daniil-Sakharov/HockeyProject/internal/initializer/app"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		log.Println("Received shutdown signal, stopping...")
		cancel()
	}()

	statsApp, err := app.NewFHSPBStatsParserApp(ctx)
	if err != nil {
		log.Fatalf("Failed to create stats parser app: %v", err)
	}

	if err := statsApp.Run(ctx); err != nil {
		log.Fatalf("Stats parser failed: %v", err)
	}

	log.Println("Stats parser completed successfully")
}
