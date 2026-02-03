import Foundation

enum APIError: LocalizedError {
    case invalidURL
    case unauthorized
    case forbidden
    case notFound
    case validationError(String)
    case serverError(statusCode: Int, message: String?)
    case decodingError(Error)
    case encodingError(Error)
    case networkError(Error)
    case unknown

    var errorDescription: String? {
        switch self {
        case .invalidURL:
            return "Неверный URL"
        case .unauthorized:
            return "Необходима авторизация"
        case .forbidden:
            return "Доступ запрещён"
        case .notFound:
            return "Не найдено"
        case .validationError(let message):
            return message
        case .serverError(_, let message):
            return message ?? "Ошибка сервера"
        case .decodingError:
            return "Ошибка обработки данных"
        case .encodingError:
            return "Ошибка отправки данных"
        case .networkError:
            return "Ошибка сети"
        case .unknown:
            return "Неизвестная ошибка"
        }
    }
}

struct APIErrorResponse: Decodable {
    let error: String?
    let message: String
}
