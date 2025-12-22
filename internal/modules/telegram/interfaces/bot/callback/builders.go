package callback

// Menu создает callback data для главного меню
func Menu(cmd string) string {
	return ActionMenu + ":" + cmd
}

// Filter создает callback data для фильтров
func Filter(filterType, value string) string {
	return ActionFilter + ":" + filterType + ":" + value
}

// Search создает callback data для поиска
func Search(cmd string) string {
	return ActionSearch + ":" + cmd
}

// SearchPageDirection создает callback data для пагинации
func SearchPageDirection(direction string) string {
	return ActionSearch + ":" + SearchPage + ":" + direction
}

// Player создает callback data для профиля игрока
func Player(action, playerID string) string {
	return ActionPlayer + ":" + action + ":" + playerID
}

// Report создает callback data для скачивания отчета
func Report(playerID string) string {
	return ActionReport + ":" + playerID
}
