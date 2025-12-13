package main

import (
	"context"
	"log"

	"github.com/Daniil-Sakharov/HockeyProject/internal/initializer/app"
)

func main() {
	ctx := context.Background()

	botApp, err := app.NewBotApp(ctx)
	if err != nil {
		log.Fatalf("Failed to create bot app: %v", err)
	}

	if err := botApp.Run(ctx); err != nil {
		log.Fatalf("Bot failed: %v", err)
	}
}
