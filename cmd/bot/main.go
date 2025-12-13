package main

import (
	"context"
	"log"

	"github.com/Daniil-Sakharov/HockeyProject/internal/initializer"
)

func main() {
	ctx := context.Background()

	botApp, err := initializer.NewBotApp(ctx)
	if err != nil {
		log.Fatalf("Failed to create bot initializer: %v", err)
	}

	if err := botApp.Run(ctx); err != nil {
		log.Fatalf("Bot failed: %v", err)
	}
}
