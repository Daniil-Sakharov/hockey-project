package main

import (
	"context"
	"log"

	"github.com/Daniil-Sakharov/HockeyProject/internal/initializer"
)

func main() {
	ctx := context.Background()

	parserApp, err := initializer.NewParserApp(ctx)
	if err != nil {
		log.Fatalf("Failed to create parser initializer: %v", err)
	}

	if err := parserApp.Run(ctx); err != nil {
		log.Fatalf("Parser failed: %v", err)
	}
}
