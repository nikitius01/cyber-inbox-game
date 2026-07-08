# Инспектор входящих

Обучающая игра по распознаванию фишинга: игрок свайпает письма влево как фишинг или вправо как легитимные.

## Структура

- `frontend/` — React + JS, Vite, интерфейс игры, профиль, вход и регистрация.
- `backend/` — Go API, авторизация, задачи, RAW-данные, AI-генерация.
- `docs/` — API, база данных и заметки по безопасности.

## Запуск

Локальная конфигурация находится в `.env`. AI-ключи в нем оставлены пустыми, их нужно указать вручную.

PostgreSQL:

```bash
docker compose up -d postgres
```

Backend:

```bash
cd backend
DATABASE_URL="postgres://cyber:cyber@localhost:5432/cybersecurity_game?sslmode=disable" \
ENABLE_DEV_USER_LIST=true \
go run ./cmd/server
```

Frontend:

```bash
cd frontend
pnpm install
pnpm run dev
```

Переменные окружения:

```bash
SERVER_ADDR=:8080
FRONTEND_ORIGIN=http://localhost:5173
JWT_SECRET=change-me
DATABASE_URL=postgres://cyber:cyber@localhost:5432/cybersecurity_game?sslmode=disable
OPENAI_API_KEY=sk-...
AI_API_KEY=sk-...
AI_BASE_URL=https://api.openai.com/v1/chat/completions
AI_MODEL=gpt-4.1-mini
AI_TIMEOUT_SECONDS=90
```

`AI_BASE_URL` можно заменить на URL другого OpenAI-compatible провайдера. Для старого варианта `OPENAI_API_KEY` тоже поддерживается, но предпочтительнее использовать `AI_API_KEY`.

## Список аккаунтов

В dev-режиме аккаунты можно посмотреть:

```bash
curl http://localhost:8080/api/dev/users
```

Эндпоинт работает только при `ENABLE_DEV_USER_LIST=true`.

Если `DATABASE_URL` задан, аккаунты сохраняются в PostgreSQL. В `docker-compose.yml` база хранит данные в volume `postgres_data`, поэтому перезапуск Go-сервера аккаунты не удаляет.

Напрямую в PostgreSQL:

```bash
docker compose exec postgres psql -U cyber -d cybersecurity_game -c "select id, username, email, created_at from users;"
```

Подробная инструкция по запуску и рестарту: `docs/server.md`.
