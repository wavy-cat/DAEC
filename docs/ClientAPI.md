# Client API Reference

> Работает по протоколу HTTP. Формат данных - JSON.

## Эндпойнты не требующие авторизации

<details>
<summary><kbd>POST</kbd> <code>/api/v1/register</code></summary>

**Регистрирует пользователя.**

Тело запроса:

```json
{
  "login": "string",
  "password": "string"
}
```

Тело ответа при коде 200 отсутствует.

Возможные коды ответа:

* `200 Ok` - пользователь успешно зарегистрирован.
* `409 Conflict` - переданный логин уже используется.
* `422 Unprocessable Entity` - неправильная семантика данных запроса.
* `500 Internal Server Error` - внутренняя ошибка сервера.

Пример запроса:

[В .http файле](../examples/client.http)

```bash
curl -X POST -H "Content-Type: application/json" -d '{"login": "login", "password": "password"}' http://localhost/api/v1/register
```

</details>

<details>
<summary><kbd>POST</kbd> <code>/api/v1/login</code></summary>

**Выдаёт JWT-токен для авторизации.**

Тело запроса:

```json
{
  "login": "string",
  "password": "string"
}
```

Тело ответа:

```json5
{
  "token": "jwt-token",
  "expire": 1000 // время в формате unix когда истечёт токен
}
```

Возможные коды ответа:

* `200 Ok` - всё окей.
* `403 Forbidden` - неправильный пароль.
* `404 Not found` - пользователь не найден.
* `422 Unprocessable Entity` - неправильная семантика данных запроса.
* `500 Internal Server Error` - внутренняя ошибка сервера.

Пример запроса:

[В .http файле](../examples/client.http)

```bash
curl -X POST -H "Content-Type: application/json" -d '{"login": "admin", "password": "admin"}' http://localhost/api/v1/login
```
</details>

## Эндпойнты требующие авторизации

> В данных эндпойнтах нужно обязательно указывать заголовок `Authorization` с JWT-токеном.
> Иначе будет возвращён код ответа `401 Unauthorized`.
> [Подробнее](Authorization.md).

<details>
<summary><kbd>GET</kbd> <code>/api/v1/expressions</code></summary>

**Отправляет список выражений пользователя.**

Тело ответа:

```json5
{
  "expressions": [
    {
      "id": 1,
      "status": "done", // Возможные значения: `pending`, `error`, `done`
      "result": 10,
      "content": "5+5" // Само выражение
    }
  ]
}
```

Возможные коды ответа:

* `200 Ok` - всё окей.
* `500 Internal Server Error` - внутренняя ошибка сервера.

Пример запроса:

[В .http файле](../examples/client.http)

```bash
curl -X GET -H "Authorization: Bearer YOUR_TOKEN_HERE" http://localhost/api/v1/expressions
```
</details>

<details>
<summary><kbd>GET</kbd> <code>/api/v1/expressions/{id}</code></summary>

**Отправляет информацию о конкретном выражении пользователя.**

Тело ответа:

```json5
{
  "expression": {
    "id": 1,
    "status": "done", // Возможные значения: `pending`, `error`, `done`
    "result": 10,
    "content": "5+5" // Само выражение
  }
}
```

Возможные коды ответа:

* `200 Ok` - всё окей.
* `404 Not found` - выражение не найдено.
* `500 Internal Server Error` - внутренняя ошибка сервера.

Пример запроса:

[В .http файле](../examples/client.http)

```bash
curl -X GET -H "Authorization: Bearer YOUR_TOKEN_HERE" http://localhost/api/v1/expressions/YOUR_ID_HERE
```
</details>

<details>
<summary><kbd>POST</kbd> <code>/api/v1/calculate</code></summary>

**Принимает выражение и отправляет его на обработку.**

Тело запроса:

```json
{
  "expression": "10+10"
}
```

> Разрешённые символы в выражении: 
> `0`, `1`, `2`, `3`, `4`, `5`, `6`, `7`, `8`, `9`, `+`, `-`, `(`, `)`, `/`, `*`, `^`, `.` и пробел.

Тело ответа:

```json
{
  "id": 1
}
```

Возможные коды ответа:

* `201 Created` - всё окей.
* `422 Unprocessable Entity` - неправильная семантика данных запроса, либо выражение содержит запрещённые символы.
* `500 Internal Server Error` - внутренняя ошибка сервера.

Пример запроса:

[В .http файле](../examples/client.http)

```bash
curl -X POST -H "Authorization: Bearer YOUR_TOKEN_HERE" -H "Content-Type: application/json" -d '{"expression": "YOUR_EXPRESSION_HERE"}' http://localhost/api/v1/calculate
```
</details>
