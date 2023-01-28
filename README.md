# Сервис, предоставляющий API по созданию сокращённых ссылок

Ссылка должна быть:

- Уникальной: на один оригинальный URL должна ссылаться только одна сокращенная ссылка;

- Длиной 10 символов;

- Из символов латинского алфавита в нижнем и верхнем регистре, цифр и символа _ (подчеркивание).

# Использование

**Клонирование репозитория**

`$ git clone git@github.com:kuzkuss/url_service.git && cd url_service`

**Сборка**

`$ docker compose build --no-cache`

**Запуск**

`$ docker compose up -d`

Для изменения используемого хранилища необходимо изменить файл config/config.toml:
```
database = "postgres" - для использования Postgres

database = "in_memory" - для использования in memory
```

**Отправление запросов**

- POST запрос:

`$ curl -X POST http://0.0.0.0:8080/create -H 'Content-Type: application/json' -d '{"original_link":"https://www.golang.org"}'`

Ответ:

`{"body":{"short_link":"uXQ71UxAzr"}}`

- GET запрос:

`curl -X GET http://127.0.0.1:8080/get/uXQ71UxAzr`

Ответ:

`{"body":{"original_link":"https://www.golang.org"}}`

Более подробно описано в swagger документации в `docs/swagger.yaml`

**Тестирование**

`$ go test ./... -coverpkg ./... -coverprofile=c.out`

Для проверки покрытия кода тестами:

`$ cat c.out | grep -v -E ".*/mocks|.*/proto"  > c_res.out`

`$ go tool cover -func=c_res.out`

