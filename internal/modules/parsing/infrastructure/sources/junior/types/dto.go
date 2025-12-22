package types

// PlayerDTO DTO для игрока из junior.fhr.ru
type PlayerDTO struct {
	Number     string // Номер игрока
	ProfileURL string // URL профиля (/player/...)
	Name       string // ФИО
	BirthDate  string // Дата рождения (13.05.2008)
	Position   string // Позиция (Защитник/Нападающий/Вратарь)
	Height     string // Рост
	Weight     string // Вес
	Handedness string // Хват (Левый/Правый)
}

// TeamDTO DTO для команды
type TeamDTO struct {
	URL  string // Относительный URL команды
	Name string // Название команды
	City string // Город
}

// TournamentDTO DTO для турнира
type TournamentDTO struct {
	ID        string // ID из URL
	Name      string // Название турнира
	URL       string // Относительный URL
	Domain    string // Домен
	Season    string // Сезон (2025/2026)
	StartDate string // Дата начала (01.09.2025)
	EndDate   string // Дата окончания (30.04.2025 или пусто)
	IsEnded   bool   // Флаг завершенности (comp-ended class)
}

// SeasonInfo информация о сезоне из дропдауна
type SeasonInfo struct {
	Name    string // "2025/2026"
	AjaxURL string // "/fhr-ajax/.../getTournamentsList/?season=2025-2026"
}

// YearLink информация о годе рождения из AJAX-ссылки
type YearLink struct {
	Year    int    // Год рождения (2008, 2009, ...)
	AjaxURL string // AJAX URL для загрузки команд этого года
}
