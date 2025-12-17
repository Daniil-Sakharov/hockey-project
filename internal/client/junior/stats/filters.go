package stats

import (
	"context"
	"net/http"

	"github.com/Daniil-Sakharov/HockeyProject/internal/client/junior/stats/parsing"
	"github.com/PuerkitoBio/goquery"
)

// Combination псевдоним для parsing.StatsCombination
type Combination = parsing.StatsCombination

// ParseCombinations парсит готовые комбинации год+группа из data-ajax атрибутов
// DEPRECATED: Используйте ParseCombinationsWithAjax для полного парсинга
func ParseCombinations(doc *goquery.Document) ([]StatsCombination, error) {
	result, err := parsing.ParseCombinations(doc)
	if err != nil {
		return nil, err
	}

	// Конвертируем parsing.StatsCombination в stats.StatsCombination
	combinations := make([]StatsCombination, len(result))
	for i, r := range result {
		combinations[i] = StatsCombination{
			YearID:    r.YearID,
			YearLabel: r.YearLabel,
			GroupID:   r.GroupID,
			GroupName: r.GroupName,
		}
	}

	return combinations, nil
}

// ParseCombinationsWithAjax парсит комбинации год+группа с двухуровневой логикой
// 1. Извлекает все ГОДЫ из dropdown
// 2. Для КАЖДОГО года делает AJAX запрос
// 3. Из AJAX ответа извлекает ГРУППЫ для этого года
func ParseCombinationsWithAjax(
	ctx context.Context,
	doc *goquery.Document,
	domain string,
	httpClient *http.Client,
) ([]StatsCombination, error) {
	result, err := parsing.ParseCombinationsWithAjax(ctx, doc, domain, httpClient)
	if err != nil {
		return nil, err
	}

	// Конвертируем parsing.StatsCombination в stats.StatsCombination
	combinations := make([]StatsCombination, len(result))
	for i, r := range result {
		combinations[i] = StatsCombination{
			YearID:    r.YearID,
			YearLabel: r.YearLabel,
			GroupID:   r.GroupID,
			GroupName: r.GroupName,
		}
	}

	return combinations, nil
}
