package main

import (
	"context"
	"log"
	"time"

	"github.com/chromedp/chromedp"
)

func main() {
	log.Println("🏒 Hockey Stats Parser - Starting...")

	// Настройки chromedp - ВИДИМЫЙ браузер для тестирования
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false), // false = видим браузер!
		chromedp.Flag("disable-gpu", false),
		chromedp.Flag("no-sandbox", true),
		chromedp.WindowSize(1920, 1080),
	)

	// Создаём контекст с нашими настройками
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	// Создаём контекст браузера
	ctx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	defer cancel()

	// Устанавливаем таймаут
	ctx, cancel = context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	log.Println("🌐 Открываем registry.fhr.ru...")

	// Открываем сайт
	err := chromedp.Run(ctx,
		chromedp.Navigate("https://registry.fhr.ru"),
	)
	if err != nil {
		log.Fatalf("❌ Ошибка при открытии сайта: %v", err)
	}

	log.Println("✅ Сайт успешно загружен!")
	log.Println("🔍 Браузер останется открытым 30 секунд для проверки...")

	// Держим браузер открытым чтобы ты видел результат
	time.Sleep(30 * time.Second)

	log.Println("👋 Закрываем браузер...")
}
