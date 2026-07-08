# Локальный Запуск Сервера

## Что Где Хранится

PostgreSQL поднимается через Docker Compose.

Локальная база:

```text
host: localhost
port: 5432
database: cybersecurity_game
user: cyber
password: cyber
```

Строка подключения для backend:

```text
postgres://cyber:cyber@localhost:5432/cybersecurity_game?sslmode=disable
```

Данные сохраняются в Docker volume:

```text
cybersecuritygame_postgres_data
```

## Первый Запуск

Из корня проекта:

```bash
cd "/Users/nikita/учеба/CyberSecurity Game"
docker compose up -d postgres
cd backend
set -a
source ../.env
set +a
go run ./cmd/server
```

Во втором окне терминала:

```bash
cd "/Users/nikita/учеба/CyberSecurity Game/frontend"
PATH="/Users/nikita/.cache/codex-runtimes/codex-primary-runtime/dependencies/node/bin:/Users/nikita/.cache/codex-runtimes/codex-primary-runtime/dependencies/bin:$PATH"
pnpm run dev
```

Сайт:

```text
http://localhost:5173/
```

API:

```text
http://localhost:8080/
```

## Рестарт Backend

Останови backend в терминале через `Ctrl+C`, затем:

```bash
cd "/Users/nikita/учеба/CyberSecurity Game/backend"
set -a
source ../.env
set +a
go run ./cmd/server
```

## Рестарт PostgreSQL

```bash
cd "/Users/nikita/учеба/CyberSecurity Game"
docker compose restart postgres
```

## Полный Рестарт Через Docker Compose

```bash
cd "/Users/nikita/учеба/CyberSecurity Game"
docker compose up -d postgres
docker compose restart postgres
```

## Посмотреть Аккаунты

Через API:

```bash
curl http://localhost:8080/api/dev/users
```

Напрямую в PostgreSQL:

```bash
cd "/Users/nikita/учеба/CyberSecurity Game"
docker compose exec postgres psql -U cyber -d cybersecurity_game -c "select id, username, email, created_at from users;"
```

## AI Ключи

В `.env` оставлены пустые поля:

```text
AI_API_KEY=
OPENAI_API_KEY=
```

Для OpenAI-compatible провайдера можно поменять:

```text
AI_BASE_URL=
AI_MODEL=
AI_TIMEOUT_SECONDS=
```
