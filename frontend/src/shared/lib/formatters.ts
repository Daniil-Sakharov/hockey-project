export function formatNumber(value: number): string {
  return new Intl.NumberFormat('ru-RU').format(value)
}

export function formatDate(date: string | Date): string {
  return new Intl.DateTimeFormat('ru-RU', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  }).format(new Date(date))
}

const ABBREVIATIONS: Record<string, string> = {
  пфо: 'ПФО',
  ппфо: 'ППФО',
  урфо: 'УрФО',
  цфо: 'ЦФО',
  сзфо: 'СЗФО',
  юфо: 'ЮФО',
  сфо: 'СФО',
  дфо: 'ДФО',
  мо: 'МО',
  ло: 'ЛО',
  спб: 'СПб',
}

// Regex для римских цифр (I, V, X, L, C, D, M)
const ROMAN_NUMERAL_REGEX = /^[ivxlcdm]+$/i

/**
 * Проверяет является ли слово римской цифрой и возвращает её в верхнем регистре
 */
function fixRomanNumeral(word: string): string | null {
  if (ROMAN_NUMERAL_REGEX.test(word)) {
    return word.toUpperCase()
  }
  return null
}

/**
 * Убирает возрастные суффиксы и исправляет регистр аббревиатур и римских цифр.
 * "Xiii Зимняя Спартакиада" → "XIII Зимняя Спартакиада"
 * "Первенство Пфо 18/17/16/15 лет" → "Первенство ПФО"
 */
export function cleanTournamentName(name: string): string {
  return name
    .replace(/\s+до\s+\d+\s*(?:лет|Лет)/i, '')
    .replace(/\s+\d+(?:\/\d+)*\s*(?:лет|Лет)/i, '')
    .trim()
    .replace(/\S+/g, (word) => {
      // Сначала проверяем римские цифры
      const roman = fixRomanNumeral(word)
      if (roman) return roman
      // Затем аббревиатуры
      return ABBREVIATIONS[word.toLowerCase()] ?? word
    })
}

export function pluralize(count: number, one: string, few: string, many: string): string {
  const mod10 = count % 10
  const mod100 = count % 100

  if (mod100 >= 11 && mod100 <= 19) return many
  if (mod10 === 1) return one
  if (mod10 >= 2 && mod10 <= 4) return few
  return many
}

/**
 * Форматирует название группы, добавляя "Группа" если нужно.
 * "А" → "Группа А"
 * "B" → "Группа B"
 * "1" → "Группа 1"
 * "Группа А" → "Группа А" (без изменений)
 * "Плей-офф" → "Плей-офф" (без изменений)
 * "Финал" → "Финал" (без изменений)
 */
export function formatGroupName(name: string): string {
  if (!name) return name

  // Если уже содержит "Группа" - вернуть как есть
  if (/группа/i.test(name)) return name

  // Специальные названия которые не нужно менять
  const specialNames = ['плей-офф', 'финал', 'полуфинал', 'четвертьфинал', 'квалификация', 'утешительный']
  if (specialNames.some(s => name.toLowerCase().includes(s))) return name

  // Если это одна буква (кириллица или латиница) или цифра - добавить "Группа"
  if (/^[А-ЯA-Z0-9]$/i.test(name.trim())) {
    return `Группа ${name.trim().toUpperCase()}`
  }

  // Иначе вернуть как есть
  return name
}
