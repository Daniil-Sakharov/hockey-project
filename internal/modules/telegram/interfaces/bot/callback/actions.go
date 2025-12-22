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
