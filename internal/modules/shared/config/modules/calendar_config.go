package modules

import "time"

// CalendarConfig конфигурация для парсеров календарей
type CalendarConfig struct {
	Junior   JuniorCalendarConfig   `env:"JUNIOR_CALENDAR"`
	FHSPB    FHSPBCalendarConfig    `env:"FHSPB_CALENDAR"`
	MIHF     MIHFCalendarConfig     `env:"MIHF_CALENDAR"`
	FHMoscow FHMoscowCalendarConfig `env:"FHMOSCOW_CALENDAR"`
}

// JuniorCalendarConfig конфигурация Junior календаря
type JuniorCalendarConfig struct {
	RequestDelay      time.Duration `env:"JUNIOR_CALENDAR_REQUEST_DELAY" default:"150ms"`
	TournamentWorkers int           `env:"JUNIOR_CALENDAR_TOURNAMENT_WORKERS" default:"3"`
	GameWorkers       int           `env:"JUNIOR_CALENDAR_GAME_WORKERS" default:"5"`
	ParseProtocol     bool          `env:"JUNIOR_CALENDAR_PARSE_PROTOCOL" default:"true"`
	ParseLineups      bool          `env:"JUNIOR_CALENDAR_PARSE_LINEUPS" default:"true"`
	ParseVideo        bool          `env:"JUNIOR_CALENDAR_PARSE_VIDEO" default:"true"`
	SkipExisting      bool          `env:"JUNIOR_CALENDAR_SKIP_EXISTING" default:"true"`
}

// FHSPBCalendarConfig конфигурация FHSPB календаря
type FHSPBCalendarConfig struct {
	RequestDelay      time.Duration `env:"FHSPB_CALENDAR_REQUEST_DELAY" default:"150ms"`
	TournamentWorkers int           `env:"FHSPB_CALENDAR_TOURNAMENT_WORKERS" default:"3"`
	GameWorkers       int           `env:"FHSPB_CALENDAR_GAME_WORKERS" default:"5"`
	ParseProtocol     bool          `env:"FHSPB_CALENDAR_PARSE_PROTOCOL" default:"true"`
	ParseLineups      bool          `env:"FHSPB_CALENDAR_PARSE_LINEUPS" default:"true"`
	SkipExisting      bool          `env:"FHSPB_CALENDAR_SKIP_EXISTING" default:"true"`
}

// MIHFCalendarConfig конфигурация MIHF календаря
type MIHFCalendarConfig struct {
	RequestDelay      time.Duration `env:"MIHF_CALENDAR_REQUEST_DELAY" default:"150ms"`
	TournamentWorkers int           `env:"MIHF_CALENDAR_TOURNAMENT_WORKERS" default:"3"`
	GameWorkers       int           `env:"MIHF_CALENDAR_GAME_WORKERS" default:"5"`
	ParseProtocol     bool          `env:"MIHF_CALENDAR_PARSE_PROTOCOL" default:"true"`
	SkipExisting      bool          `env:"MIHF_CALENDAR_SKIP_EXISTING" default:"true"`
	MinBirthYear      int           `env:"MIHF_CALENDAR_MIN_BIRTH_YEAR" default:"2008"`
	MaxSeasons        int           `env:"MIHF_CALENDAR_MAX_SEASONS" default:"0"`
	TestSeason        string        `env:"MIHF_CALENDAR_TEST_SEASON" default:""`
	RetryEnabled      bool          `env:"MIHF_CALENDAR_RETRY_ENABLED" default:"true"`
	RetryMaxAttempts  int           `env:"MIHF_CALENDAR_RETRY_MAX_ATTEMPTS" default:"3"`
	RetryDelay        time.Duration `env:"MIHF_CALENDAR_RETRY_DELAY" default:"5s"`
}

// FHMoscowCalendarConfig конфигурация FHMoscow календаря
type FHMoscowCalendarConfig struct {
	RequestDelay time.Duration `env:"FHMOSCOW_CALENDAR_REQUEST_DELAY" default:"150ms"`
	SkipExisting bool          `env:"FHMOSCOW_CALENDAR_SKIP_EXISTING" default:"true"`
}

// Методы для реализации интерфейса CalendarConfig

func (c *JuniorCalendarConfig) RequestDelayMs() int {
	return int(c.RequestDelay.Milliseconds())
}

func (c *JuniorCalendarConfig) GetTournamentWorkers() int {
	return c.TournamentWorkers
}

func (c *JuniorCalendarConfig) GetGameWorkers() int {
	return c.GameWorkers
}

func (c *JuniorCalendarConfig) GetParseProtocol() bool {
	return c.ParseProtocol
}

func (c *JuniorCalendarConfig) GetParseLineups() bool {
	return c.ParseLineups
}

func (c *JuniorCalendarConfig) GetSkipExisting() bool {
	return c.SkipExisting
}

// Методы для FHSPBCalendarConfig

func (c *FHSPBCalendarConfig) RequestDelayMs() int {
	return int(c.RequestDelay.Milliseconds())
}

func (c *FHSPBCalendarConfig) GetGameWorkers() int {
	return c.GameWorkers
}

func (c *FHSPBCalendarConfig) GetParseProtocol() bool {
	return c.ParseProtocol
}

func (c *FHSPBCalendarConfig) GetParseLineups() bool {
	return c.ParseLineups
}

func (c *FHSPBCalendarConfig) GetSkipExisting() bool {
	return c.SkipExisting
}

// Методы для MIHFCalendarConfig

func (c *MIHFCalendarConfig) RequestDelayMs() int {
	return int(c.RequestDelay.Milliseconds())
}

func (c *MIHFCalendarConfig) GetMinBirthYear() int {
	return c.MinBirthYear
}

func (c *MIHFCalendarConfig) GetRequestDelay() int {
	return int(c.RequestDelay.Milliseconds())
}

func (c *MIHFCalendarConfig) GetGameWorkers() int {
	return c.GameWorkers
}

func (c *MIHFCalendarConfig) GetParseProtocol() bool {
	return c.ParseProtocol
}

func (c *MIHFCalendarConfig) GetSkipExisting() bool {
	return c.SkipExisting
}

func (c *MIHFCalendarConfig) GetMaxSeasons() int {
	return c.MaxSeasons
}

func (c *MIHFCalendarConfig) GetTestSeason() string {
	return c.TestSeason
}

func (c *MIHFCalendarConfig) GetRetryEnabled() bool {
	return c.RetryEnabled
}

func (c *MIHFCalendarConfig) GetRetryMaxAttempts() int {
	return c.RetryMaxAttempts
}

func (c *MIHFCalendarConfig) GetRetryDelay() time.Duration {
	return c.RetryDelay
}
