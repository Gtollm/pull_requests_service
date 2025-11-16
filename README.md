# Pull Request Service

## Запуск

Требуется установленный Docker и Docker Compose.

Перед запуском укажите параметры БД в файле `.env` (POSTGRES_USER, POSTGRES_PASSWORD, POSTGRES_DB).

Из директории `deployments/` выполните команду:

```bash
docker-compose up --build
```

Сервис будет доступен по адресу http://localhost:8080.
