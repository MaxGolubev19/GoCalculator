# Калькулятор арифметических выражений

Этот проект представляет собой веб-сервис, который вычисляет арифметические выражения с использованием стандартных операций (сложение, вычитание, умножение, деление) и скобок.

### Статусы HTTP-ответов

- **200 OK** — Выражение успешно вычислено.
- **422 Unprocessable Entity** — Некорректное выражение.
- **500 Internal Server Error** — Внутренняя ошибка сервера.

### Пример успешного запроса

Если выражение вычислено успешно:

```bash
curl --location 'localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": "2+2*2"
}'
```

Ответ:

```json
{
  "result": 6
}
```

### Пример ошибки 422 (некорректное выражение)

Если выражение содержит буквы, лишние скобки, неправильный порядок операндов (например, "", "1+a", "1+1*", "2+2**2", "((2+2-*(2"):

```bash
curl --location 'localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": "1+a"
}'
```

Ответ:

```json
{
  "error": "Expression is not valid"
}
```

### Пример ошибки 500 (внутренняя ошибка сервера)

Если произошла непредвиденная ошибка на сервере, например, проблемы с вычислением выражения:

```bash
curl --location 'localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": "42/0"
}'
```

Ответ:

```json
{
  "error": "Internal server error"
}
```

## Запуск проекта

Для запуска проекта используйте команду:
```bash
go run ./cmd/main.go
```

Сервер будет запущен на порту `8080`, и вы сможете отправлять запросы по адресу `http://localhost:8080/api/v1/calculate`. Если нужно использовать другой порт, вы можете задать переменную окружения PORT:
```bash
export PORT=8081
go run ./cmd/main.go
```

```sh
$env:PORT=8081
go run ./cmd/main.go
```