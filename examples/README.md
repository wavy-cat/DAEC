# Примеры

В этой папке есть 2 файла формата `.http`.
Используйте их в какой-нибудь IDE от JetBrains, например, GoLand.
Или в VSCode с плагином [Rest Client](https://marketplace.visualstudio.com/items?itemName=humao.rest-client).
Или где-нибудь ещё.

Также, примеры из задания полностью подходят к этому решению.
Ниже они продублированы.

> [!NOTE]
> Это не API Reference, это всего лишь примеры.
> Чтобы посмотреть эндпойнты детально, перейдите к файлу README.md.

## Примеры на стороне клиента

### Добавление вычисления арифметического выражения

Запрос:

```bash
curl --location 'localhost/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
      "expression": "строка с выражением"
}'
```

Ответ:

```json
{
  "id": "уникальный идентификатор выражения (UUID)"
}
```

### Получение списка выражений

Запрос:

```bash
curl --location 'localhost/api/v1/expressions'
```

Ответ:

```json5
{
  "expressions": [
    {
      "id": "идентификатор выражения (UUID)",
      "status": "статус вычисления выражения",
      "result": 0.0  // результат выражения (вещественное число)
    },
    {
      "id": "идентификатор выражения (UUID)",
      "status": "статус вычисления выражения",
      "result": 1.0  // результат выражения (вещественное число)
    }
  ]
}

```

### Получение выражения по его идентификатору

Запрос:
```bash
curl --location 'localhost/api/v1/expressions/{id}'
```

Ответ:
```json5
{
  "expression": {
    "id": "идентификатор выражения (UUID)",
    "status": "статус вычисления выражения",
    "result": 2.0  // результат выражения (вещественное число)
  }
}
```

## Примеры на стороне агента (демона)

### Получение задачи для выполнения.

Запрос:
```bash
curl --location 'localhost/internal/task'
```

Ответ:
```json5
{
       "task":
             {
                   "id": "идентификатор задачи (JSON)",
                   "arg1": 1.5, // первый аргумент (вещественное число)
                   "arg2": 3, // второй аргумент (вещественное число)
                   "operation": "операция", // Один символ
                   "operation_time": 5000 // время выполнения операции (целочисленное число)
              }
}
 
```

### Прием результата обработки данных.

Запрос:
```bash
curl --location 'localhost/internal/task' \
--header 'Content-Type: application/json' \
--data '{
      "id": "идентификатор задачи",
      "result": <результат (вещественное число)>
}'
```

Тела ответа нет. Ориентируйтесь по HTTP коду ответа.