package main

import (
	"context"
	"fmt"
	"log"
	"strings"
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
	ctx, cancel = context.WithTimeout(ctx, 600*time.Second)
	defer cancel()

	// ======== АВТОРИЗАЦИЯ ========
	username := "MN_SAHAROV"  // ← ЗАМЕНИТЬ
	password := "necv8iniWr5" // ← ЗАМЕНИТЬ

	log.Println("🌐 Открываем registry.fhr.ru...")

	// Открываем сайт
	err := chromedp.Run(ctx,
		chromedp.Navigate("https://registrynew.fhr.ru/"),
	)
	if err != nil {
		log.Fatalf("❌ Ошибка при открытии сайта: %v", err)
	}

	log.Println("✅ Сайт загружен")
	log.Println("🔐 Начинаем авторизацию...")

	// Ждём появления формы логина и заполняем её
	err = chromedp.Run(ctx,
		// Ждём появления поля логина (новый селектор: input#input-vaadin-text-field-10)
		chromedp.WaitVisible(`input[name="username"]`, chromedp.ByQuery),
		chromedp.Sleep(500*time.Millisecond), // Небольшая пауза для стабильности

		// Вводим логин
		chromedp.SendKeys(`input[name="username"]`, username, chromedp.ByQuery),
		chromedp.Sleep(300*time.Millisecond),

		// Вводим пароль (новый селектор: input с name="password")
		chromedp.SendKeys(`input[name="password"]`, password, chromedp.ByQuery),
		chromedp.Sleep(300*time.Millisecond),

		// Кликаем на кнопку "Войти" (vaadin-button с slot="submit")
		chromedp.Click(`vaadin-button[slot="submit"]`, chromedp.ByQuery),
	)
	if err != nil {
		log.Fatalf("❌ Ошибка при заполнении формы логина: %v", err)
	}

	log.Println("⏳ Ждём результата авторизации...")

	// Ждём успешной авторизации - проверяем появление главного меню с навигацией
	err = chromedp.Run(ctx,
		chromedp.Sleep(3*time.Second), // Ждём перехода
		// Проверяем что мы на главной странице (появилось боковое меню с элементом "Школа" или "Главное меню")
		chromedp.WaitVisible(`vaadin-side-nav-item`, chromedp.ByQuery),
	)
	if err != nil {
		log.Fatalf("❌ Авторизация не удалась: %v", err)
	}

	log.Println("✅ Авторизация успешна!")

	// ======== НАВИГАЦИЯ: ПРЯМОЙ ПЕРЕХОД НА СТРАНИЦУ СПОРТСМЕНОВ ========
	log.Println("🔗 Переходим напрямую на страницу 'Спортсмены' по URL...")

	// Переходим на прямой URL страницы спортсменов
	targetURL := "https://registrynew.fhr.ru/m/b/712458/Спортсмены%20Федерации%20хоккея%20России"
	err = chromedp.Run(ctx,
		chromedp.Navigate(targetURL),
	)
	if err != nil {
		log.Fatalf("❌ Ошибка при переходе на URL: %v", err)
	}

	log.Println("✅ Переход выполнен, ждём загрузки страницы...")

	// Ждём загрузки страницы спортсменов
	err = chromedp.Run(ctx,
		chromedp.Sleep(10*time.Second), // Даём время Vaadin для обработки hash navigation и загрузки контента
	)
	if err != nil {
		log.Fatalf("❌ Ошибка при ожидании загрузки: %v", err)
	}

	log.Println("✅ Страница 'Спортсмены' должна быть загружена")

	// ======== ПРИМЕНЯЕМ ФИЛЬТР ========
	log.Println("🔎 Ищем кнопку 'Применить фильтр'...")

	// Сначала ждём появления диалога с фильтрами (vaadin-dialog-overlay)
	err = chromedp.Run(ctx,
		chromedp.WaitVisible(`vaadin-dialog-overlay`, chromedp.ByQuery),
	)
	if err != nil {
		log.Fatalf("❌ Не удалось дождаться появления диалога фильтров: %v", err)
	}

	log.Println("✅ Диалог фильтров открыт")

	// Селектор для кнопки "Применить фильтр" (vaadin-button с текстом внутри)
	// Используем XPath для поиска по тексту внутри кнопки
	filterButtonSelector := `//vaadin-button[contains(., 'Применить фильтр')]`

	// Ждём появления кнопки фильтра
	err = chromedp.Run(ctx,
		chromedp.WaitVisible(filterButtonSelector, chromedp.BySearch),
	)
	if err != nil {
		log.Fatalf("❌ Не удалось найти кнопку 'Применить фильтр': %v", err)
	}

	log.Println("✅ Кнопка 'Применить фильтр' найдена")

	// Кликаем на кнопку
	log.Println("🖱️ Кликаем на 'Применить фильтр'...")
	err = chromedp.Run(ctx,
		chromedp.Click(filterButtonSelector, chromedp.BySearch),
	)
	if err != nil {
		log.Fatalf("❌ Ошибка при клике на 'Применить фильтр': %v", err)
	}

	log.Println("✅ Фильтр применён, ждём загрузки списка игроков...")

	// Ждём загрузки списка игроков
	err = chromedp.Run(ctx,
		chromedp.Sleep(5*time.Second), // Даём время на загрузку данных
	)
	if err != nil {
		log.Fatalf("❌ Ошибка при ожидании загрузки списка: %v", err)
	}

	// Ждём появления СТРОК в таблице (не просто самой таблицы)
	log.Println("⏳ Ждём появления строк игроков в таблице...")
	// Попробуем несколько селекторов на случай изменения структуры
	err = chromedp.Run(ctx,
		chromedp.WaitVisible(`vaadin-grid-cell-content`, chromedp.ByQuery),
	)
	if err != nil {
		log.Fatalf("❌ Строки игроков не появились: %v", err)
	}

	// Дополнительная пауза для стабильности Vaadin
	time.Sleep(2 * time.Second)

	log.Println("✅ Список игроков загружен!")

	// ======== ПРОВЕРКА ВИДИМЫХ СТРОК ========
	log.Println("📊 Проверяем количество видимых строк в таблице...")

	var visibleRowsCount int
	// Подсчитываем строки через tbody tr внутри vaadin-grid
	err = chromedp.Run(ctx,
		chromedp.Evaluate(`
			const grid = document.querySelector('vaadin-grid');
			if (!grid) {
				throw new Error('Vaadin Grid не найден');
			}
			
			// Пробуем найти tbody tr напрямую
			let rows = grid.querySelectorAll('tbody tr');
			
			// Если не нашли, пробуем через Shadow DOM
			if (rows.length === 0 && grid.shadowRoot) {
				const tbody = grid.shadowRoot.querySelector('tbody');
				if (tbody) {
					rows = tbody.querySelectorAll('tr');
				}
			}
			
			// Если всё ещё 0, пробуем через slot table
			if (rows.length === 0) {
				const table = grid.querySelector('table');
				if (table) {
					rows = table.querySelectorAll('tbody tr');
				}
			}
			
			rows.length;
		`, &visibleRowsCount),
	)
	if err != nil {
		log.Printf("⚠️ Не удалось получить количество строк: %v", err)
	} else {
		log.Printf("📊 Видимых строк в таблице: %d", visibleRowsCount)
	}

	// Проверка что таблица не пустая
	if visibleRowsCount == 0 {
		log.Println("❌ Таблица пустая! Данные не загрузились.")
		log.Fatal("❌ Невозможно продолжить: таблица пустая (0 строк)")
	}

	// ======== СОХРАНЕНИЕ URL СПИСКА ========
	log.Println("💾 Сохраняем URL списка игроков для возврата...")

	var savedListURL string
	err = chromedp.Run(ctx,
		chromedp.Location(&savedListURL),
	)
	if err != nil {
		log.Printf("⚠️ Не удалось получить URL списка: %v", err)
	} else {
		log.Printf("💾 Сохранён URL: %s", savedListURL)
	}

	// Проверяем что URL правильный (содержит маркер списка спортсменов)
	if !strings.Contains(savedListURL, "712458") && !strings.Contains(savedListURL, "Спортсмены") {
		log.Println("⚠️ URL не похож на страницу списка игроков!")
		log.Printf("⚠️ Получили: %s", savedListURL)
		log.Println("⚠️ Ожидали URL с '712458' или 'Спортсмены'")

		// Пробуем перейти напрямую
		log.Println("🔄 Пробуем перейти напрямую к списку...")
		savedListURL = "https://registrynew.fhr.ru/m/b/712458/Спортсмены%20Федерации%20хоккея%20России"
		err = chromedp.Run(ctx,
			chromedp.Navigate(savedListURL),
		)
		if err != nil {
			log.Fatalf("❌ Ошибка при переходе: %v", err)
		}

		time.Sleep(5 * time.Second)

		// Повторно применяем фильтр
		log.Println("🔄 Повторно применяем фильтр...")
		err = chromedp.Run(ctx,
			chromedp.Click(filterButtonSelector, chromedp.BySearch),
		)
		if err != nil {
			log.Fatalf("❌ Ошибка при повторном применении фильтра: %v", err)
		}

		time.Sleep(5 * time.Second)

		// Проверяем строки снова
		err = chromedp.Run(ctx,
			chromedp.Evaluate(`
				const grid = document.querySelector('vaadin-grid');
				let rows = grid ? grid.querySelectorAll('tbody tr') : [];
				if (rows.length === 0 && grid && grid.shadowRoot) {
					const tbody = grid.shadowRoot.querySelector('tbody');
					if (tbody) rows = tbody.querySelectorAll('tr');
				}
				rows.length;
			`, &visibleRowsCount),
		)
		log.Printf("📊 Строк после повторного фильтра: %d", visibleRowsCount)

		if visibleRowsCount == 0 {
			log.Fatal("❌ Таблица всё ещё пустая после повторной попытки!")
		}
	}

	// ======== ПАРСИНГ ИГРОКОВ БЛОКАМИ ========
	log.Println("\n🎯 Начинаем парсинг игроков блоками...")
	log.Println("=" + strings.Repeat("=", 60))

	// Настройки парсинга
	const totalPlayers = 102555    // Всего игроков
	const playersPerBlock = 35     // Видимых строк в таблице
	const blocksToTest = 2         // Для теста: 2 блока
	const playersPerBlockTest = 15 // Для теста: 15 игроков в блоке (было 3)

	currentBlockStart := 0

	// Цикл по блокам
	for blockNum := 0; blockNum < blocksToTest; blockNum++ {
		log.Printf("\n📦 БЛОК #%d (глобальный индекс %d-%d)", blockNum+1, currentBlockStart, currentBlockStart+playersPerBlockTest-1)
		log.Println("─" + strings.Repeat("─", 60))

		// Скроллим к началу блока (один раз на блок!)
		if currentBlockStart > 0 {
			log.Printf("⬇️ Скроллим к индексу %d...", currentBlockStart)
			err = chromedp.Run(ctx,
				chromedp.Evaluate(fmt.Sprintf(`
					(() => {
						const grid = document.querySelector('vaadin-grid');
						if (grid && typeof grid.scrollToIndex === 'function') {
							grid.scrollToIndex(%d);
							return 'Scrolled to index %d';
						}
						return 'scrollToIndex not found';
					})();
				`, currentBlockStart, currentBlockStart), nil),
			)
			if err != nil {
				log.Fatalf("❌ Ошибка при скролле: %v", err)
			}

			log.Println("⏳ Ждём подгрузки данных после скролла...")
			time.Sleep(3 * time.Second)
			log.Println("✅ Скролл выполнен!")

			// Дополнительная пауза перед получением координат
			log.Println("⏸️  Дополнительная пауза 2 секунды для стабилизации...")
			time.Sleep(2 * time.Second)
		}

		// Парсим игроков в текущем блоке
		for localIndex := 0; localIndex < playersPerBlockTest; localIndex++ {
			globalIndex := currentBlockStart + localIndex

			log.Printf("\n👤 Игрок #%d (блок %d, локальный индекс %d)", globalIndex, blockNum+1, localIndex)
			log.Println("  " + strings.Repeat("·", 58))

			// Получаем координаты ячейки с Фамилией через JavaScript + детальная информация
			jsGetCoords := fmt.Sprintf(`
				(function() {
					const grid = document.querySelector('vaadin-grid');
					if (!grid) {
						throw new Error('Vaadin Grid не найден');
					}
					
					let rows = grid.querySelectorAll('tbody tr');
					
					if (rows.length === 0 && grid.shadowRoot) {
						const tbody = grid.shadowRoot.querySelector('tbody');
						if (tbody) {
							rows = tbody.querySelectorAll('tr');
						}
					}
					
					if (rows.length === 0) {
						const table = grid.querySelector('table');
						if (table) {
							rows = table.querySelectorAll('tbody tr');
						}
					}
					
					if (!rows || rows.length === 0) {
						throw new Error('Нет видимых строк в таблице');
					}
					
					const targetRow = rows[%d];
					if (!targetRow) {
						throw new Error('Целевая строка не найдена');
					}
					
					const cells = targetRow.querySelectorAll('td');
					if (!cells || cells.length < 2) {
						throw new Error('Недостаточно ячеек в строке');
					}
					
					const familyCell = cells[1];
					const rect = familyCell.getBoundingClientRect();
					
					// Дополнительная информация для отладки
					const cellText = familyCell.textContent.trim();
					const isVisible = rect.top >= 0 && rect.top < window.innerHeight;
					const firstVisibleIndex = grid._firstVisibleIndex || 'unknown';
					
					return {
						x: rect.left + 10,
						y: rect.top + rect.height / 2,
						cellText: cellText,
						rectTop: rect.top,
						rectHeight: rect.height,
						isVisible: isVisible,
						windowHeight: window.innerHeight,
						firstVisibleIndex: firstVisibleIndex,
						totalRows: rows.length
					};
				})();
			`, localIndex)

			var coords struct {
				X                 float64     `json:"x"`
				Y                 float64     `json:"y"`
				CellText          string      `json:"cellText"`
				RectTop           float64     `json:"rectTop"`
				RectHeight        float64     `json:"rectHeight"`
				IsVisible         bool        `json:"isVisible"`
				WindowHeight      float64     `json:"windowHeight"`
				FirstVisibleIndex interface{} `json:"firstVisibleIndex"`
				TotalRows         int         `json:"totalRows"`
			}

			err = chromedp.Run(ctx,
				chromedp.Evaluate(jsGetCoords, &coords),
			)
			if err != nil {
				log.Printf("  ❌ Ошибка при получении координат: %v", err)
				continue
			}

			log.Printf("  📍 Координаты: x=%.1f, y=%.1f", coords.X, coords.Y)
			log.Printf("  📋 Текст ячейки: '%s'", coords.CellText)
			log.Printf("  🔍 Детали: rectTop=%.1f, height=%.1f, видна=%v", coords.RectTop, coords.RectHeight, coords.IsVisible)
			log.Printf("  🗂️  Grid: firstVisibleIndex=%v, всего строк=%d", coords.FirstVisibleIndex, coords.TotalRows)

			// Шаг 1: Одинарный клик (выделение строки)
			log.Println("  👆 Шаг 1: Одинарный клик (выделение строки)...")
			err = chromedp.Run(ctx,
				chromedp.MouseClickXY(coords.X, coords.Y),
			)
			if err != nil {
				log.Printf("  ❌ Ошибка при одинарном клике: %v", err)
				continue
			}

			// Пауза для обработки выделения
			log.Println("  ⏸️  Пауза 300ms...")
			time.Sleep(300 * time.Millisecond)

			// Шаг 2: Двойной клик (открытие профиля)
			log.Println("  👆 Шаг 2: Двойной клик (открытие профиля)...")
			err = chromedp.Run(ctx,
				chromedp.MouseClickXY(coords.X, coords.Y, chromedp.ClickCount(2)),
			)
			if err != nil {
				log.Printf("  ❌ Ошибка при двойном клике: %v", err)
				continue
			}

			log.Println("  ⏳ Ждём загрузки профиля...")
			time.Sleep(5 * time.Second)

			// Проверяем URL
			var currentURL string
			err = chromedp.Run(ctx,
				chromedp.Location(&currentURL),
			)
			if err == nil {
				if strings.Contains(currentURL, "Form-S-Update-Web") || strings.Contains(currentURL, "Web-Спортсмены") {
					log.Printf("  ✅ Профиль открыт! URL: %s", currentURL)
				} else {
					log.Printf("  ⚠️ URL не изменился: %s", currentURL)
				}
			}

			// TODO: Здесь будет парсинг данных профиля
			log.Println("  📋 Парсинг данных (TODO)...")

			// Возвращаемся к списку
			log.Println("  ◀️ Возвращаемся к списку...")
			err = chromedp.Run(ctx,
				chromedp.Navigate(savedListURL),
			)
			if err != nil {
				log.Fatalf("  ❌ Ошибка при возврате: %v", err)
			}

			log.Println("  ⏳ Ждём загрузки списка...")
			time.Sleep(5 * time.Second) // Увеличено с 3 до 5 секунд

			// Ждём появления строк в таблице
			log.Println("  🔄 Ждём готовности таблицы...")
			err = chromedp.Run(ctx,
				chromedp.WaitVisible(`vaadin-grid-cell-content`, chromedp.ByQuery),
			)
			if err != nil {
				log.Printf("  ⚠️ Таблица не готова: %v", err)
			}

			// Дополнительная пауза для Vaadin рендеринга
			time.Sleep(2 * time.Second)

			log.Printf("  ✅ Игрок #%d завершён!", globalIndex)
		}

		// Переходим к следующему блоку
		currentBlockStart += playersPerBlockTest
	}

	// ИТОГИ
	log.Println("\n" + strings.Repeat("=", 60))
	log.Println("🎉 Тестирование парсинга блоками завершено успешно!")
	log.Printf("✅ Обработано блоков: %d", blocksToTest)
	log.Printf("✅ Обработано игроков: %d", blocksToTest*playersPerBlockTest)
	log.Println(strings.Repeat("=", 60))

	// ======== ТЕСТИРОВАНИЕ СКРОЛЛА ========
	// 📝 ДОКУМЕНТАЦИЯ: Рабочий код для скролла таблицы игроков
	//
	// Назначение:
	//   - Прокрутка виртуального списка игроков (Vaadin Grid)
	//   - Подгрузка новых данных через виртуальный рендеринг
	//
	// Как работает:
	//   1. Находит скроллер: `.v-grid-scroller-vertical`
	//   2. Выполняет N итераций скролла через JS: scrollTop += pixels
	//   3. Ждёт 1 секунду между скроллами для подгрузки данных
	//   4. Логирует текущую/новую позицию для отладки
	//
	// Параметры:
	//   - scrollIterations: количество прокруток (по умолчанию 10)
	//   - scrollStep: пикселей за одну прокрутку (по умолчанию 800)
	//
	// Применение:
	//   - Для парсинга всех игроков нужно будет:
	//     * Увеличить scrollIterations (для 102k игроков ~25-30 строк = ~4000 итераций)
	//     * Добавить извлечение данных из видимых строк после каждого скролла
	//     * Добавить проверку на достижение конца списка
	//
	// Статус: ✅ ПРОТЕСТИРОВАНО, РАБОТАЕТ
	// ============================================================

	/*
		log.Println("📜 Начинаем тестирование скролла...")

		// Селектор для вертикального скроллера таблицы
		scrollerSelector := `.v-grid-scroller-vertical`

		// Проверяем что скроллер существует
		err = chromedp.Run(ctx,
			chromedp.WaitVisible(scrollerSelector, chromedp.ByQuery),
		)
		if err != nil {
			log.Fatalf("❌ Не удалось найти скроллер таблицы: %v", err)
		}

		log.Println("✅ Скроллер найден, начинаем прокрутку...")

		// Делаем несколько итераций скролла для проверки
		scrollIterations := 10
		scrollStep := 800 // Пикселей за одну прокрутку

		for i := 1; i <= scrollIterations; i++ {
			log.Printf("🔄 Итерация скролла %d/%d...", i, scrollIterations)

			// Получаем текущую позицию скролла (для отладки)
			var currentScrollTop int
			err = chromedp.Run(ctx,
				chromedp.Evaluate(fmt.Sprintf(`document.querySelector('%s').scrollTop`, scrollerSelector), &currentScrollTop),
			)
			if err != nil {
				log.Printf("⚠️ Не удалось получить позицию скролла: %v", err)
			} else {
				log.Printf("📍 Текущая позиция скролла: %d пикселей", currentScrollTop)
			}

			// Скроллим вниз
			err = chromedp.Run(ctx,
				chromedp.Evaluate(fmt.Sprintf(`document.querySelector('%s').scrollTop += %d`, scrollerSelector, scrollStep), nil),
			)
			if err != nil {
				log.Printf("❌ Ошибка при скролле на итерации %d: %v", i, err)
				break
			}

			log.Printf("✅ Скролл выполнен на +%d пикселей", scrollStep)

			// Пауза для подгрузки данных
			time.Sleep(1 * time.Second)

			// Получаем новую позицию после скролла
			var newScrollTop int
			err = chromedp.Run(ctx,
				chromedp.Evaluate(fmt.Sprintf(`document.querySelector('%s').scrollTop`, scrollerSelector), &newScrollTop),
			)
			if err == nil {
				log.Printf("📍 Новая позиция скролла: %d пикселей (сдвиг: +%d)", newScrollTop, newScrollTop-currentScrollTop)
			}
		}

		log.Println("✅ Тестирование скролла завершено!")
	*/

	log.Println("🔍 Браузер останется открытым 30 секунд для визуальной проверки...")

	// Держим браузер открытым для визуальной проверки
	time.Sleep(30 * time.Second)

	log.Println("👋 Закрываем браузер...")
}
