# Распределённый вычислитель арифметических выражений

## Описание проекта

Это распределённая система вычисления арифметических выражений, в которой вычисления масштабируемы, асинхронны и надёжны за счёт:

- взаимодействия между компонентами через gRPC;

- хранения состояния в SQLite (персистентность);

- поддержки многопользовательского режима с авторизацией через JWT.


Пользователь может регистрироваться, входить в систему, отправлять выражения на вычисление и получать результат после завершения вычислений.

Система состоит из двух основных компонентов:
1. **Оркестратор** — принимает выражения от пользователей, разбивает их на отдельные операции и распределяет задачи между вычислителями (агентами).
2. **Агент** — получает задачи от оркестратора по gRPC, выполняет операции (сложение, вычитание, умножение, деление) и отправляет результат обратно.

## Основные возможности

- Поддержка базовых арифметических операций: +, -, *, /

- Хранение истории вычислений по пользователю

- Персистентность (SQLite)

- gRPC-связь между оркестратором и агентом

- JWT-авторизация

- REST API для взаимодействия с пользователями

## Архитектура

### 1. Оркестратор
Оркестратор управляет задачами вычислений, а также предоставляет REST API для пользователей и gRPC API для агентов. Он:
- Принимает выражения и разбивает их на задачи.
- Раздаёт задачи агентам и получает результаты.
- Хранит статус выражения и пользователей в SQLite.

### 2. Агент
Агент запрашивает задачи у оркестратора, выполняет их и возвращает результаты. Он работает многопоточно, создавая несколько вычислительных горутин (количество настраивается переменной среды `COMPUTING_POWER`).

## Структура проекта

Проект состоит из следующих директорий:

- **`cmd/orchestrator/`** — точка входа оркестратора.
- **`cmd/agent/`** — точка входа агента.
- **`internal/orchestrator/`** — внутренняя логика оркестратора (управление пользователями, выражениями, задачами, вычислениями, связь с базой данных).
- **`internal/agent/`** — внутренняя логика агента (получение задач, вычисления, отправка результатов).
- **`pkg/`** — вспомогательные библиотеки и модули.
- **`Dockerfile`** и **`docker-compose.yml`** — файлы для контейнеризации и развертывания системы.

## Запуск проекта

1. Клонируй репозиторий:
    ```sh
    git clone https://github.com/MaxGolubev19/GoCalculator.git
    cd GoCalculator
    ```

2. Создай файл `.env` в корне проекта:
    ```
    PUBLIC_PORT=8080
    GRPC_PORT=50051

    SECRET_KEY = "secret key"

    COMPUTING_POWER=10

    TIME_ADDITION_MS=100
    TIME_SUBTRACTION_MS=200
    TIME_MULTIPLICATIONS_MS=300
    TIME_DIVISIONS_MS=400
    ```

3. Запусти систему с помощью Docker:
    ```sh
    docker-compose up --build
    ```

## Переменные окружения

| Переменная                | Описание                                     |
|---------------------------|----------------------------------------------|
| `PUBLIC_PORT`             | Порт для REST API оркестратора               |
| `GRPC_PORT`               | Порт, на котором агент слушает gRPC-запросы  |
| `SECRET_KEY`              | Секретный ключ для подписи JWT               |
| `COMPUTING_POWER`         | Количество горутин для вычислений            |
| `TIME_ADDITION_MS`        | Время выполнения операции сложения (мс)      |
| `TIME_SUBTRACTION_MS`     | Время выполнения операции вычитания (мс)     |
| `TIME_MULTIPLICATIONS_MS` | Время выполнения операции умножения (мс)     |
| `TIME_DIVISIONS_MS`       | Время выполнения операции деления (мс)       |

## API

### Регистрация
```sh
curl --location --request POST 'localhost:8080/api/v1/register' \
--header 'Content-Type: application/json' \
--data '{
  "login": "user1",
  "password": "123456"
}'
```
**Коды ответов:**
- `201 Created` — пользователь зарегистрирован.
- `409 Conflict` — пользователь уже существует.
- `500 Internal Server Error` — внутренняя ошибка сервера.

### Вход (авторизация)
```sh
curl --location --request POST 'localhost:8080/api/v1/login' \
--header 'Content-Type: application/json' \
--data '{
  "login": "user1",
  "password": "123456"
}'
```
**Коды ответов:**
- `200 OK` — выражение принято для вычисления.

    ```json
    {
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6..."
    }
    ```

- `401 Unauthorized` — некорректный логин или пароль.

- `500 Internal Server Error` — внутренняя ошибка сервера.

### Добавление выражения
```sh
curl --location 'localhost/api/v1/calculate' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer <jwt_token>' \
--data '{
  "expression": "2+2*2"
}'
```
**Коды ответов:**
- `201 Created` — выражение принято для вычисления.
    ```json
    {
        "id": 42
    }
    ```
- `422 Unprocessable Entity` — некорректные данные.

- `500 Internal Server Error` — внутренняя ошибка сервера.


### Получение списка выражений
```sh
curl --location 'localhost/api/v1/expressions' \
--header 'Authorization: Bearer <jwt_token>'
```

**Коды ответов:**
- `200 OK` — успешно получен список выражений.

    ```json
    {
        "expressions": [
            {
                "id": 42,
                "status": "IN PROGRESS",
                "result": 0
            }
        ]
    }
    ```

- `500 Internal Server Error` — внутренняя ошибка сервера.

### Получение выражения по ID
```sh
curl --location 'localhost:8080/api/v1/expressions/42' \
--header 'Authorization: Bearer <jwt_token>'
```

**Коды ответов:**
- `200 OK` — успешно получено выражение.
    ```json
    {
        "expression": {
            "id": 42,
            "status": "DONE",
            "result": 6
        }
    }
    ```

- `404 Not Found` — выражение не найдено.

- `500 Internal Server Error` — внутренняя ошибка сервера.

## Тестирование
Проект покрыт тестами, которые можно запустить командой:
```sh
go test ./...
```