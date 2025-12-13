package main

import (
	"context"
	"log"

	"github.com/Daniil-Sakharov/HockeyProject/internal/initializer"
)

func main() {
	ctx := context.Background()

	statsParserApp, err := initializer.NewStatsParserApp(ctx)
	if err != nil {
		log.Fatalf("Failed to create stats parser initializer: %v", err)
	}

	if err := statsParserApp.Run(ctx); err != nil {
		log.Fatalf("Stats parser failed: %v", err)
	}
}
