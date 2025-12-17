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

	fhspbApp, err := app.NewFHSPBParserApp(ctx)
	if err != nil {
		log.Fatalf("Failed to create fhspb parser app: %v", err)
	}

	if err := fhspbApp.Run(ctx); err != nil {
		log.Fatalf("Parser failed: %v", err)
	}

	log.Println("Parser completed successfully")
}
