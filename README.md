# Сервер социальной сети "Social Nyetwork"

## Описание

В данном репозитории предлагается реализация веб-сервера социальной сети "Social Nyetwork" (см. подробное описание в [корневом репозитории](https://github.com/pamugk/go-social-nyetwork)).

Актуальная структура БД может быть полностью восстановлена по предоставляемым миграциям.

Для начала предполагается хранение всех данных в одной БД, в будущем возможно распределение по специализированным хранилищам.

Подготовлена [OpenAPI](https://swagger.io/specification/) спецификации по предоставляемому сервисом REST API.

## Используемые технологии и инструментальные средства

* Язык программирования - [Go](https://go.dev/).
* База данных - [PostgreSQL](https://www.postgresql.org/).
* Применение миграций БД - [migrate](https://github.com/golang-migrate/migrate).
* Драйвер БД - [pgx](https://github.com/jackc/pgx).
* Маршрутизация HTTP-запросов - [chi](https://github.com/go-chi/chi).
* Валидация запросов - [validator](https://github.com/go-playground/validator).
* JWT аутентификация - [jwtauth](https://github.com/go-chi/jwtauth).
* RPC - [gRPC-Go](https://github.com/grpc/grpc-go).
