# Weather Bot

Telegram-бот для получения текущей погоды. Пользователь указывает город командой `/city`, после чего может в любой момент запросить погоду через `/weather`.

## Используемые библиотеки

- [go-telegram-bot-api/v5](https://github.com/go-telegram-bot-api/telegram-bot-api) — работа с Telegram Bot API
- [pgx/v5](https://github.com/jackc/pgx) — драйвер и пул соединений для PostgreSQL
- [godotenv](https://github.com/joho/godotenv) — загрузка переменных окружения из `.env`
- [goose](https://github.com/pressly/goose) — миграции базы данных

## Команды

| Команда | Описание |
|---------|----------|
| `/city <название>` | Сохраняет указанный город для пользователя |
| `/weather` | Показывает текущую температуру в сохранённом городе |

## Требования

- Go 1.26+
- PostgreSQL (через Docker Compose или локально)
- Токен Telegram-бота (через [@BotFather](https://t.me/BotFather))
- API-ключ [OpenWeatherMap](https://openweathermap.org/api)

## Запуск

1. Клонируйте репозиторий и перейдите в директорию проекта.

2. Создайте файл `.env` из примера:

```bash
cp .env.example .env
```

3. Укажите в `.env` токен бота (`BOT_TOKEN`) и ключ OpenWeatherMap (`OPENWEATHERAPI_KEY`).

4. Запустите PostgreSQL:

```bash
docker compose up -d
```

5. Примените миграции:

```bash
go run github.com/pressly/goose/v3/cmd/goose up
```

6. Запустите бота:

```bash
go run .
```
