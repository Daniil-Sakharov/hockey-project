package main

import (
	"context"
	"log"

	"github.com/Daniil-Sakharov/HockeyProject/internal/initializer/app"
)

func main() {
	ctx := context.Background()

	statsParserApp, err := app.NewStatsParserApp(ctx)
	if err != nil {
		log.Fatalf("Failed to create stats parser app: %v", err)
	}

	if err := statsParserApp.Run(ctx); err != nil {
		log.Fatalf("Stats parser failed: %v", err)
	}
}
