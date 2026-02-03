package types

import "net/http"

// HTTPRequester интерфейс для выполнения HTTP запросов
type HTTPRequester interface {
	MakeRequest(url string) (*http.Response, error)
	MakeRequestWithHeaders(url string, headers map[string]string) (*http.Response, error)
}
