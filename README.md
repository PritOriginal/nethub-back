# nethub-back

[![wakatime](https://wakatime.com/badge/user/b2a0c08d-61f2-4144-ba78-aab13a59cb9f/project/c78a72e5-c028-4c37-b598-9d3aac550a3c.svg)](https://wakatime.com/badge/user/b2a0c08d-61f2-4144-ba78-aab13a59cb9f/project/c78a72e5-c028-4c37-b598-9d3aac550a3c)

В данном репозитории представлен REST-сервис управления сетевыми устройствами. Данный сервис является частью тестового задания в компании [**Нетхаб**](https://nethub.ru/).

[nethub-front](https://github.com/PritOriginal/nethub-front) - Frontend репозиторий.

[nethub.pritoriginal.ru](https://nethub.pritoriginal.ru/) - Frontend, который общается с данным сервисом. Развёрнут на VPS.

[Swagger документация](https://nethub.pritoriginal.ru/api/swagger/index.html) - документация REST API , развёрнутая на VPS. Для локального просмотра `http://[host]:[port]/swagger/index.html`

## Эндпоинты

- `POST /devices` – создание устройства.
- `GET /devices` – список с фильтрацией по is_active и поиском по hostname (по подстроке).
- `GET /devices/{id}` – получение устройства по id.
- `PUT /devices/{id}` – обновление устройства.
- `DELETE /devices/{id}` – мягкое удаление (пометка флагом, без физического удаления).

## Стек

- `Golang`
- [`Gin`](https://github.com/gin-gonic/gin) - Веб фреймворк
- `PostgreSQL` - СУБД
- [`migrate`](https://github.com/golang-migrate/migrate) - Миграции
- `log/slog` - Логгер
- `GitHub Actions` - CI/CD  
- [`swaggo/swag`](https://github.com/swaggo/swag) - OpenAPI (Swagger)
- `Docker` - Контейнеризация

## Подготовка

Создайте конфиг

Для `.env`

```bash
cp ./configs/.env.example ./configs/.env
```

Для `.yaml`

```bash
cp ./configs/config.yaml.example ./configs/config.yaml
```

## Запуск

```bash
make run
```

Docker

```bash
make run-docker
```

## Тесты

Простой прогон тестов:

```bash
make test
```

Прогон тестов с выводом покрытия:

```bash
make test-cover
```

## Миграции

`migrate create`:

```bash
make migrate NAME_MIGRATION="name_migration" 
```

`migrate up`:

```bash
make migrate-up
```

`migrate up 1`:

```bash
make migrate-up-1
```

`migrate down`:

```bash
make migrate-down
```

`migrate down 1`:

```bash
make migrate-down-1
```

## Swagger

Чтобы обновить документацию после внесения изменений, выполните следующую команду:

```bash
make swag
```


## Примечание

Если в качестве конфигурационного файла был выбран `.yaml`, то замените путь к конфигурационному файлу в `Makefile` либо запускайте приложение командой:

```bash
go run ./cmd/rest/ --config=./configs/config.yaml
```