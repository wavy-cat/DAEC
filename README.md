# DAEC

**D**istributed **A**rithmetic **E**xpression **C**alculator

Aka Распределенный вычислитель арифметических выражений.

> [!IMPORTANT]
> Это не полноценный проект, а всего-лишь решение на финальную задачу в Яндекс Лицее.
> Вся документация и комментарии в коде ведутся на русском языке.

## Архитектура проекта

Основные части приложения:

* Backend aka оркестратор
* Agent (demon) aka вычислитель
* Frontend (не имеет отдельного сервера) aka веб-интерфейс (не реализовано :) )

Вычисления происходят в вычислителе, оркестратор лишь распараллеливает задачу и передаёт её частями агенту.
Также оркестратор хранит данные о задачах.

В качестве основного веб-сервера предлагается использовать [Caddy](https://caddyserver.com/).
Он уже присутствует в Docker Compose.

Все идентификаторы используют тип UUID.

Приложение не использует никакую внешнюю СУБД. Все данные хранятся в памяти.

## Что поддерживается

* Целочисленные числа
* Вещественные числа
* Унарный минус
* Операторы +, -, *, /, ^
* Скобки

## Как запустить систему

### Запуск в Docker Compose

> [!NOTE]
> Для этого способа настраивать прокси **не** нужно, если вы из России!
> В Dockerfile уже указаны образы с [прокси Timeweb Cloud](https://dockerhub.timeweb.cloud/).

Для этого способа у вас должен быть установлен Git, Docker и Docker Compose.
Инструкции по установке Docker доступны на [официальном сайте](https://www.docker.com/get-started/) (заходить с VPN из
России).
Можно также использовать Docker Desktop.
В случае чего находите инструкцию по установке на вашей системе в гугле.

Перед запуском настройте переменные среды в файле .env.

```bash
git clone https://github.com/wavy-cat/DAEC.git
cd DAEC
docker compose up -d
```

Для выключения:

```bash
cd DAEC # Если ещё не в папке с проектом
docker compose down
```

После запуска будет доступно API на http://localhost/api/v1/

Внутренняя часть API (/internal/task) недоступна вне сети контейнеров.

> [!WARNING]
> Если вы видите код ответа 502 Bad Gateway, то не спешите паниковать!
> Это происходит потому что Caddy запускается раньше оркестратора.
> Поэтому просто подождите пока сервер запустится полностью.

### Запуск напрямую

> [!IMPORTANT]
> Этот способ у вас может не заработать по ряду причин. Лучше используйте запуск в Docker Compose.
> Траблшутинг лежит на вашей ответственности.

Шаги по запуску:

1. Склонируйте репозиторий на свой компьютер и перейдите в него

Для этого у вас должен быть установлен [Git](https://git-scm.com/).
Либо можете просто скачать [архив с репозиторием](https://github.com/wavy-cat/DAEC/archive/refs/heads/main.zip) и
распаковать его.

А также само самой должен быть установлен [Go](https://go.dev), версии не ниже 1.22.

```bash
git clone https://github.com/wavy-cat/DAEC.git
cd DAEC
```

2. Установите зависимости сервера

```bash
cd backend
go mod download
```

3. Укажите задержку выполнения в переменных среды

Везде для примера задержка 5000мс, вы можете её изменить.

В Linux / macOS:
```bash
export TIME_ADDITION_MS=5000
export TIME_SUBTRACTION_MS=5000
export TIME_MULTIPLICATIONS_MS=5000
export TIME_DIVISIONS_MS=5000
export TIME_EXPONENTIATION_MS=5000
```

В Windows:
```powershell
SET TIME_ADDITION_MS=5000
SET TIME_SUBTRACTION_MS=5000
SET TIME_MULTIPLICATIONS_MS=5000
SET TIME_DIVISIONS_MS=5000
SET TIME_EXPONENTIATION_MS=5000
```

4. Запустите сервер

```bash
go run backend/cmd/backend
```

> [!NOTE]
> Если вы используйте Linux, то сервер у вас не запустится на портах меньше 1024 (в т.ч. 80 и 443)!
> Чтобы решить проблему, запускайте сервер с root правами (с пользователя root или через sudo/doas).
> Но не забудьте указать переменные среды для правильного пользователя.
> Либо можете просто сменить порт в конфиге сервера и агента
> (`/backend/internal/config/config.go` и `/agent/config/config.go`).
> Но в таком случае запросы нужно тоже отправлять на другой порт.

Пример для запуска через sudo (в Linux):
```bash
sudo \    
TIME_ADDITION_MS=5000 \
TIME_SUBTRACTION_MS=5000 \
TIME_MULTIPLICATIONS_MS=5000 \
TIME_DIVISIONS_MS=5000 \
TIME_EXPONENTIATION_MS=5000 \
go run backend/cmd/backend
```

В примере выше переменные указываются во время запуска.

3. *В другом окне* установите зависимости агента

```bash
cd DAEC/agent
go mod download
```

5. Укажите computing power:

В командах указывается 10, вы можете заменить на сколько угодно, но больше 0.

В Linux / macOS:

```bash
export COMPUTING_POWER=10
```

В Windows:

```powershell
SET COMPUTING_POWER=10
```

6. Запустите агента

*Можно даже парочку разом.*

```bash
go run agent
```

После этого API будет доступно по адресу http://localhost/api/v1 (если вы не указали другой порт в конфиге).

> [!NOTE]
> Если вы видите при запуске ошибку `Failed to get КАКОЕ-ЛИБО_НАЗВАНИЕ_ПЕРЕМЕННОЙ value: strconv.Atoi`,
> то скорее всего вы либо неправильно указали необходимую переменную среду,
> либо указали для неправильного пользователя
> (например, если в Linux запускаете через sudo или su, но export делаете в своём пользователе),
> либо не указали вообще.

## Унарный минус

Это знак "-" перед которым стоит не цифра (пробел, скобка (но не закрывающая) или другой оператор
либо если он просто находится в самом начале выражения) и после которого идёт цифра:

* `40-2` — от 40 отнять 2 = 38
* `-40-2` — от -40 отнять 2 = -42
* `- 40-2` — неправильная форма записи
* `40- 2` — от 40 отнять 2 = 38
* `40 - ( -2 )` — от 40 отнять -2 = 42
* `40-(-2)` — от 40 отнять -2 = 42

## API Reference

> [!NOTE]
> Примеры находятся в папке [examples](examples).

API полностью соответствует представленному в задаче (за исключением идентификаторов, они используют uuid).

### Client-side

#### Добавление вычисления арифметического выражения

`POST /api/v1/calculate`

Тело запроса:

```json5
{
  "expression": "строка с выражением"
}
```

Тело ответа:

```json5
{
  "id": "уникальный идентификатор выражения (UUID)"
}
```

Возможные коды ответа:

* 201 - выражение принято для вычисления
* 422 - невалидные данные / ошибка в выражении
* 500 - что-то пошло не так

#### Получение списка выражений

`GET localhost/api/v1/expressions`

Тело ответа:
```json5
{
  "expressions": [
    {
      "id": "идентификатор выражения (UUID)",
      "status": "статус вычисления выражения", // допустимые значения: pending, error, done
      "result": 0.0  // результат выражения (вещественное число)
    },
    {
      "id": "идентификатор выражения (UUID)",
      "status": "статус вычисления выражения", // допустимые значения: pending, error, done
      "result": 1.0  // результат выражения (вещественное число)
    }
  ]
}
```

Возможные коды ответа:
* 200 - успешно получен список выражений
* 500 - что-то пошло не так

#### Получение выражения по его идентификатору

`GET /api/v1/expressions/{id}`

Тело ответа:
```json5
{
  "expression": {
    "id": "идентификатор выражения (UUID)",
    "status": "статус вычисления выражения", // допустимые значения: queue, processing, done
    "result": 2.5  // результат выражения (вещественное число)
  }
}
```

Возможные коды ответа:
* 200 - успешно получено выражение
* 404 - выражение не найдено
* 500 - что-то пошло не так

### Agent-side

#### Получение задачи для выполнения

`GET /internal/task`

Тело ответа:
```json5
{
  "task": {
    "id": "идентификатор задачи (UUID)",
    "arg1": 1.5, // первый аргумент (вещественное число)
    "arg2": 3, // второй аргумент (вещественное число)
    "operation": "операция", // Один символ
    "operation_time": 5000 // время выполнения операции (целочисленное число)
  }
}
```

Возможные коды ответа:
* 200 - успешно получена задача
* 404 - нет задачи
* 500 - что-то пошло не так

#### Приём результата обработки данных

`POST /internal/task`

Тело запроса:
```json5
{
  "id": "идентификатор задачи (UUID)",
  "result": 10.5 // результат (вещественное число), либо null
}
```

Возможные коды ответа:
* 200 - успешно записан результат
* 404 - нет такой задачи
* 422 - невалидные данные
* 500 - что-то пошло не так
