package types

import "fmt"

// Parsing error codes
const (
	// Temporary parsing errors (retryable)
	CodeParsingNetworkTimeout = "PARSING_NETWORK_TIMEOUT"
	CodeParsingRateLimit      = "PARSING_RATE_LIMIT"
	CodeParsingServerError    = "PARSING_SERVER_ERROR"
	CodeParsingConnectionFail = "PARSING_CONNECTION_FAIL"

	// Permanent parsing errors (not retryable)
	CodeParsingNotFound      = "PARSING_NOT_FOUND"
	CodeParsingInvalidFormat = "PARSING_INVALID_FORMAT"
	CodeParsingAccessDenied  = "PARSING_ACCESS_DENIED"
	CodeParsingInvalidData   = "PARSING_INVALID_DATA"
)

// NewParsingTemporaryError создает временную ошибку парсинга
func NewParsingTemporaryError(code, message string) *DomainError {
	return NewDomainError(ErrorTypeParsingTemporary, code, message)
}

// NewParsingPermanentError создает постоянную ошибку парсинга
func NewParsingPermanentError(code, message string) *DomainError {
	return NewDomainError(ErrorTypeParsingPermanent, code, message)
}

// NewNetworkTimeoutError создает ошибку таймаута сети
func NewNetworkTimeoutError(url string, timeout int) *DomainError {
	return NewParsingTemporaryError(
		CodeParsingNetworkTimeout,
		fmt.Sprintf("Network timeout after %ds", timeout),
	).WithContext("url", url).WithContext("timeout_seconds", timeout)
}

// NewRateLimitError создает ошибку превышения лимита запросов
func NewRateLimitError(source string, retryAfter int) *DomainError {
	return NewParsingTemporaryError(
		CodeParsingRateLimit,
		fmt.Sprintf("Rate limit exceeded for %s", source),
	).WithContext("source", source).WithContext("retry_after_seconds", retryAfter)
}

// NewServerError создает ошибку сервера
func NewServerError(url string, statusCode int) *DomainError {
	return NewParsingTemporaryError(
		CodeParsingServerError,
		fmt.Sprintf("Server error: HTTP %d", statusCode),
	).WithContext("url", url).WithContext("status_code", statusCode)
}
