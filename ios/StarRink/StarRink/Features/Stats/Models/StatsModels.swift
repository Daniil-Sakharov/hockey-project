import Foundation
import SwiftUI

struct GameStat: Identifiable {
    let id = UUID()
    let gameNumber: Int
    let date: Date
    let goals: Int
    let assists: Int
    let points: Int
    let plusMinus: Int
    let penaltyMinutes: Int
    let shotsOnGoal: Int
    let timeOnIce: TimeInterval

    var formattedDate: String {
        let formatter = DateFormatter()
        formatter.dateFormat = "dd.MM"
        return formatter.string(from: date)
    }

    var formattedTimeOnIce: String {
        let minutes = Int(timeOnIce) / 60
        let seconds = Int(timeOnIce) % 60
        return String(format: "%d:%02d", minutes, seconds)
    }
}

struct SeasonTrend: Identifiable {
    let id = UUID()
    let month: String
    let goals: Int
    let assists: Int

    var points: Int { goals + assists }
}

struct StatCategory: Identifiable {
    let id = UUID()
    let name: String
    let value: Double
    let maxValue: Double
    let color: StatColor

    enum StatColor {
        case cyan, purple, amber, green, red

        var swiftUIColor: Color {
            switch self {
            case .cyan: return .srCyan
            case .purple: return .srPurple
            case .amber: return .srAmber
            case .green: return .srSuccess
            case .red: return .srError
            }
        }
    }
}

// MARK: - Sample Data

extension GameStat {
    static let sampleData: [GameStat] = [
        GameStat(gameNumber: 1, date: Date().addingTimeInterval(-86400*30), goals: 1, assists: 2, points: 3, plusMinus: 1, penaltyMinutes: 0, shotsOnGoal: 5, timeOnIce: 1080),
        GameStat(gameNumber: 2, date: Date().addingTimeInterval(-86400*25), goals: 0, assists: 1, points: 1, plusMinus: -1, penaltyMinutes: 2, shotsOnGoal: 3, timeOnIce: 960),
        GameStat(gameNumber: 3, date: Date().addingTimeInterval(-86400*20), goals: 2, assists: 0, points: 2, plusMinus: 2, penaltyMinutes: 0, shotsOnGoal: 7, timeOnIce: 1200),
        GameStat(gameNumber: 4, date: Date().addingTimeInterval(-86400*15), goals: 1, assists: 3, points: 4, plusMinus: 3, penaltyMinutes: 0, shotsOnGoal: 4, timeOnIce: 1140),
        GameStat(gameNumber: 5, date: Date().addingTimeInterval(-86400*10), goals: 0, assists: 0, points: 0, plusMinus: -2, penaltyMinutes: 4, shotsOnGoal: 2, timeOnIce: 780),
        GameStat(gameNumber: 6, date: Date().addingTimeInterval(-86400*5), goals: 3, assists: 1, points: 4, plusMinus: 2, penaltyMinutes: 0, shotsOnGoal: 8, timeOnIce: 1260),
        GameStat(gameNumber: 7, date: Date().addingTimeInterval(-86400*2), goals: 1, assists: 2, points: 3, plusMinus: 1, penaltyMinutes: 2, shotsOnGoal: 6, timeOnIce: 1100),
    ]
}

extension SeasonTrend {
    static let sampleData: [SeasonTrend] = [
        SeasonTrend(month: "Сен", goals: 2, assists: 3),
        SeasonTrend(month: "Окт", goals: 4, assists: 5),
        SeasonTrend(month: "Ноя", goals: 3, assists: 4),
        SeasonTrend(month: "Дек", goals: 5, assists: 6),
        SeasonTrend(month: "Янв", goals: 4, assists: 7),
    ]
}
