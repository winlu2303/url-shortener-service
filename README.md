# URL Shortener Service

Сервис для сокращения URL-адресов на Go с полным покрытием unit-тестами и тестами HTTP-обработчиков.

## Содержание

- [Возможности](#возможности)
- [Установка](#установка)
- [API Эндпоинты](#api-эндпоинты)
- [Примеры использования cURL](#примеры-использования)
- [Тестирование и проверки покрытия кода](#тестирование)

## Возможности

- **Сокращение URL** - генерация уникального короткого идентификатора (6 символов)
- **Редирект** - автоматическое перенаправление на оригинальный URL
- **Валидация URL** - проверка корректности HTTP/HTTPS адресов
- **Защита от дубликатов** - одинаковые URL получают одинаковый короткий ID
- **Хранение в памяти** - быстрый доступ через map с RWMutex
- **Полное тестирование** - unit-тесты и HTTP-тесты с покрытием >80%

## Установка

```bash
# Клонирование репозитория
git clone https://github.com/YOUR_USERNAME/urlshortener-service.git
cd urlshortener-service

# Инициализация модуля
go mod init urlshortener-service
go mod tidy

# Запуск сервера
go run .
```

## API Эндпоинты

### POST /shorten - Сокращение URL

Создает короткий идентификатор для переданного URL.

Запрос:

```json
{
    "url": "https://example.com/very/long/path"
}
```
Ответ (201 Created)

```json
{
    "short_url": "http://localhost:8080/abc123",
    "original_url": "https://example.com/very/long/path"
}
```
Ошибки:
400 Bad Request - невалидный URL или пустая строка
405 Method Not Allowed - не POST запрос
500 Internal Server Error - внутренняя ошибка
GET /{short_url} - Получение оригинального URL


Ответ (200 OK):
```json
{
    "status": "ok"
}
```

## Примеры использования cURL

```bash
# 1. Сокращение URL
curl -X POST http://localhost:8080/shorten \
  -H "Content-Type: application/json" \
  -d '{"url":"https://google.com"}'

# 2. Переход по короткой ссылке
curl -L http://localhost:8080/abc123

# 3. Проверка здоровья
curl http://localhost:8080/health
```

## Тестирование и проверки покрытия кода

```bash
# Unit-тесты бизнес-логики
go test -v -run TestURLShortener

# HTTP-тесты обработчиков
go test -v -run TestShortenHandler
go test -v -run TestRedirectHandler

# Проверка покрытия
go test -coverprofile=coverage.out

# Просмотр покрытия в терминале
go tool cover -func=coverage.out

# Открытие HTML отчета
go tool cover -html=coverage.out
```