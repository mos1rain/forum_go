# ForumGo

Микросервисный форум с авторизацией, чатом (WebSocket), ролевым доступом, swagger-документацией, логированием и покрытием тестами.

## Структура проекта

- `cmd/` — точка входа для каждого микросервиса (`auth`, `forum`, `chat`)
- `internal/` — бизнес-логика и сервисы
- `pkg/` — общие пакеты (jwt, database)
- `migrations/` — SQL-миграции для PostgreSQL
- `frontend/` — фронтенд (статические файлы)
- `docs/` — swagger-документация
- `proto/` — gRPC-протоколы

## Микросервисы

- **Auth Service** — аутентификация, JWT, роли
- **Forum Service** — категории, посты, комментарии, ролевой доступ
- **Chat Service** — WebSocket-чат, автоудаление старых сообщений, ролевой доступ

## Быстрый старт

1. **Установите зависимости:**
   ```bash
   go mod download
   npm install --prefix frontend
   ```
2. **Запустите миграции:**
   ```bash
   migrate -path migrations -database "postgres://postgres:YOURPASS@localhost:5432/forum?sslmode=disable" up
   ```
3. **Запустите сервисы (в отдельных терминалах):**
   ```bash
   go run cmd/auth/main.go
   go run cmd/forum/main.go
   go run cmd/chat/main.go
   ```
4. **Запустите фронтенд:**
   - Откройте `frontend/public/index.html` в браузере

## Тесты и покрытие

- Запуск всех тестов:
  ```bash
  go test ./...
  ```
- Проверка покрытия для сервисов:
  ```bash
  go test -cover ./internal/auth/service
  go test -cover ./internal/forum/service
  go test -cover ./internal/chat/service
  ```
- Покрытие тестами для каждого микросервиса > 80%

## Swagger-документация
- Swagger UI доступен по пути `/swagger/` для каждого сервиса (например, http://localhost:8081/swagger/)
- Описание API — в `docs/swagger.yaml` и `docs/swagger.json`

## Архитектура и требования
- Микросервисная архитектура (auth, forum, chat)
- Взаимодействие между сервисами через gRPC
- Миграции через golang-migrate
- Покрытие тестами (unit, моки, >80%)
- Логирование через zerolog
- Swagger/OpenAPI для всех сервисов
- Чат через WebSocket, автоудаление старых сообщений
- Ролевой доступ (admin, user, guest)
- Чистая архитектура, разделение бизнес-логики и инфраструктуры

## Пример .env (если используете)
```
POSTGRES_DSN=postgres://postgres:YOURPASS@localhost:5432/forum?sslmode=disable
JWT_SECRET=your_jwt_secret
```

## Автор
- [Ваше имя]

---
Готово к сдаче! Если возникнут вопросы по запуску — пишите. 