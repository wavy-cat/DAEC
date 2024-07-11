# Авторизация в оркестраторе

Некоторые эндпойнты требуют чтобы пользователь был зарегистрирован и предоставлял им свой JWT-токен в заголовке запроса.

Для таких эндпойнтов просто передавайте в заголовке `Authorization` свой JWT-токен.

```
Authorization: Bearer {jwt-token}
```

JWT-токен действует 1 час с момента его получения.

Получить его можно отправив запрос на эндпойнт `/api/v1/login` (см. [Client API](ClientAPI.md)).