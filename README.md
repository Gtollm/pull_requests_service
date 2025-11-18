# Pull Request Service

## Запуск

Требуется установленный Docker и Docker Compose.

Проект использует стандартные значения переменных среды из `.env.example`. Для запуска выполните команду:

```bash
docker-compose up --build
```

Сервис будет доступен по адресу http://localhost:8080.

При необходимости можно переопределить переменные окружения, создав файл `.env` с параметрами:
- `POSTGRES_USER` (по умолчанию: username)
- `POSTGRES_PASSWORD` (по умолчанию: password)
- `POSTGRES_DB` (по умолчанию: pull_requests_reviewer)
