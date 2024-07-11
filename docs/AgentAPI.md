# Agent API Reference

> Работает по протоколу gRPC. Формат данных - Protobuf.

Структура сообщений находится в файлах *tasks.proto*.
Их можно найти в папках *proto* [бекенда](../backend/proto) и [агента](../agent/proto).

<details>
<summary><kbd>Pull</kbd></summary>

**Метод `Pull` возвращает задачу, которую агент должен решить.**

```protobuf
rpc Pull (Empty) returns (PullTaskResponse);
```

Принимает пустое сообщение (`Empty`), возвращает `PullTaskResponse`.

В случае если задач нет возвращает ошибку `no task yet`.
</details>

<details>
<summary><kbd>Push</kbd></summary>

**Метод `Push` принимает результат задачи, решённой агентом.**

```protobuf
rpc Push (PushTaskRequest) returns (Empty);
```

Принимает `PushTaskRequest`, возвращает пустое сообщение (`Empty`).

Может вернуть ошибку. Например, в случае, если задача не найдена: `task not found`.
</details>