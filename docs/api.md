# API

Бэкенд работает на `http://localhost:8080`.

## Auth

- `POST /api/auth/register`
- `POST /api/auth/login`
- `POST /api/auth/logout`
- `GET /api/auth/me`

Авторизация хранится в `HttpOnly SameSite=Lax` cookie `auth_token`.

## Tasks

- `GET /api/tasks/random?category=easy`
- `POST /api/tasks/{id}/answer`
- `POST /api/ai/tasks/answer`

Публичная задача не содержит `isPhishing`, `redFlags` и `explanation`. Эти данные возвращаются только после ответа.

