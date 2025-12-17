package callback

// Actions - первый уровень callback data (parts[0])
const (
	ActionMenu   = "menu"
	ActionFilter = "filter"
	ActionSearch = "search"
	ActionPlayer = "player"
	ActionReport = "download_report"
)

// Menu commands (parts[1] для menu:*)
const (
	MenuSearch = "search"
	MenuStats  = "stats"
	MenuTeam   = "team"
	MenuHelp   = "help"
	MenuMain   = "main"
)

// Filter commands (parts[1] для filter:*)
const (
	FilterBack     = "back"
	FilterReset    = "reset"
	FilterApply    = "apply"
	FilterYear     = "year"
	FilterPosition = "position"
	FilterHeight   = "height"
	FilterWeight   = "weight"
	FilterRegion   = "region"
	FilterFio      = "fio"
)

// FIO sub-commands (parts[2] для filter:fio:*)
const (
	FioSelect     = "select"
	FioLastName   = "last_name"
	FioFirstName  = "first_name"
	FioPatronymic = "patronymic"
	FioClearLast  = "clear_last"
	FioClearFirst = "clear_first"
	FioClearPatr  = "clear_patr"
	FioApply      = "apply"
	FioBack       = "back"
)

// Sub-commands (parts[2] для select/value)
const (
	SubCmdSelect = "select"
	ValueAny     = "any"
)

// Search commands (parts[1] для search:*)
const (
	SearchPage          = "page"
	SearchBackToFilters = "back_to_filters"
	SearchBackToResults = "back_to_results"
)

// Search pagination (parts[2] для search:page:*)
const (
	PageNext = "next"
	PagePrev = "prev"
)

// Player commands (parts[1] для player:*)
const (
	PlayerProfile = "profile"
)

// Builder helpers - упрощают создание callback data

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
