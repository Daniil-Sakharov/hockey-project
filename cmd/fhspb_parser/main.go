package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Daniil-Sakharov/HockeyProject/internal/initializer"
)

func main() {
	// Создаем контекст с отменой по сигналу
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Обработка сигналов завершения
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		log.Println("Received shutdown signal, stopping...")
		cancel()
	}()

	// Создаем приложение
	application, err := initializer.New(ctx)
	if err != nil {
		log.Fatalf("Failed to create initializer: %v", err)
	}

	// Запускаем парсер
	if err := application.RunFHSPBParser(ctx); err != nil {
		log.Fatalf("Parser failed: %v", err)
	}

	log.Println("Parser completed successfully")
}
