package main

import (
	"context"
	"log"

	"github.com/Daniil-Sakharov/HockeyProject/internal/initializer/app"
)

func main() {
	ctx := context.Background()

	parserApp, err := app.NewParserApp(ctx)
	if err != nil {
		log.Fatalf("Failed to create parser app: %v", err)
	}

	if err := parserApp.Run(ctx); err != nil {
		log.Fatalf("Parser failed: %v", err)
	}
}
