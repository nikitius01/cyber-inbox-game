# Database

Если задан `DATABASE_URL`, backend использует PostgreSQL для аккаунтов, профиля и статистики.

Для локального запуска через Docker Compose используется сервис `postgres`:

```text
postgres://cyber:cyber@localhost:5432/cybersecurity_game?sslmode=disable
```

Данные хранятся в Docker volume:

```text
cybersecuritygame_postgres_data
```

Поэтому перезапуск Go-сервера аккаунты больше не удаляет.

Если `DATABASE_URL` не задан, backend включает in-memory fallback. В этом режиме данные исчезают после перезапуска.

Схема PostgreSQL находится в `backend/internal/database/migrations/001_init.sql`.

Ключевая идея хранения задач: `email_payload JSONB`, где лежат письмо, RAW-заголовки, ссылки и вложения.

Посмотреть аккаунты в dev-режиме можно через API:

```bash
curl http://localhost:8080/api/dev/users
```

Или напрямую в контейнере:

```bash
docker compose exec postgres psql -U cyber -d cybersecurity_game -c "select id, username, email, created_at from users;"
```
