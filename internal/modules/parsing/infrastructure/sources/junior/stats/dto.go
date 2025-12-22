package stats

// YearFilter фильтр года рождения
type YearFilter struct {
	ID    string // "16756920"
	Label string // "2008"
}

// GroupFilter фильтр группы
type GroupFilter struct {
	ID   string // "16771932" или "all"
	Name string // "Группа А" или "Общая статистика"
}

// StatsCombination представляет готовую комбинацию год+группа для запроса к API
// Извлекается напрямую из data-ajax атрибутов (option или .filter-btn)
type StatsCombination struct {
	YearID    string // "16743907" - ID года для API
	YearLabel string // "2009" - отображаемый год
	GroupID   string // "all" или "16743965" - ID группы для API
	GroupName string // "Общая статистика" или "Первый этап"
}

// PlayerStatisticDTO DTO статистики игрока из JSON API
type PlayerStatisticDTO struct {
	GroupName string // из контекста запроса
	BirthYear string // из контекста запроса

	// HTML поля из JSON (требуют парсинга)
	Surname  string `json:"surname"`   // <a href="/player/...">...</a>
	TeamName string `json:"team_name"` // <a href="/tournaments/.../team_123/">...</a>

	// Основная статистика
	GP        string `json:"gp"`        // <div>11</div>
	G         string `json:"g"`         // <div>13</div>
	A         string `json:"a"`         // <div>18</div>
	PTS       string `json:"pts"`       // <div>31</div>
	Plus      string `json:"plus"`      // <div>43</div>
	Minus     string `json:"minus"`     // <div>18</div>
	PlusMinus string `json:"plusminus"` // <div>25</div>
	PIM       string `json:"pim"`       // <div>0</div>

	// Детальная статистика голов
	ESG string `json:"esg"` // <div>7</div> - равенство
	PPG string `json:"ppg"` // <div>5</div> - большинство
	SHG string `json:"shg"` // <div>1</div> - меньшинство
	G1P string `json:"g1p"` // <div>5</div> - 1 период
	G2P string `json:"g2p"` // <div>3</div> - 2 период
	G3P string `json:"g3p"` // <div>5</div> - 3 период
	GOT string `json:"got"` // <div>0</div> - овертайм
	HT  string `json:"ht"`  // <div>3</div> - хет-трики
	WB  string `json:"wb"`  // <div>0</div> - решающие буллиты

	// Средние показатели
	GPG   string `json:"gpg"`   // <div>1,18</div>
	PTSPG string `json:"ptspg"` // <div>2,82</div>
	PPM   string `json:"ppm"`   // <div>0,00</div>
}

// StatsResponse ответ от JSON API
type StatsResponse struct {
	RecordsTotal    int                  `json:"recordsTotal"`
	RecordsFiltered int                  `json:"recordsFiltered"`
	Data            []PlayerStatisticDTO `json:"data"`
}
