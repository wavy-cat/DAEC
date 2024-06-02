# DAEC

**D**istributed **A**rithmetic **E**xpression **C**alculator

aka Распределенный вычислитель арифметических выражений

> [!IMPORTANT]
> Это не полноценный проект, а всего-лишь решение на финальную задачу в Яндекс Лицее.
> Вся документация ведётся на русском языке.

## Архитектура проекта

В проекте реализована микросервисная архитектура.

Основные части приложения:

* Backend aka оркестратор
* Agent (demon) aka вычислитель
* Frontend (не имеет отдельного сервера) aka веб-интерфейс

В качестве основного веб-сервера используется [Caddy](https://caddyserver.com/).
Он уже присутствует в Docker Compose.

Все идентификаторы используют тип UUID.

*TODO: тут будет диаграмма архитектуры*

Приложение не использует никакую внешнюю СУБД. Все данные хранятся в памяти.

## Как запустить систему

### Запуск в Docker Compose

```bash
git clone https://github.com/wavy-cat/DAEC.git
cd DAEC
docker compose up -d
```

После ввода данных команд будет доступен веб-интерфейс по адресу https://localhost/, а API на https://localhost/api/v1/

Внутренняя часть API (/internal/task) недоступна вне сети контейнеров.

Поддерживаются протоколы http и https, даже в localhost (подробнее об этом на [сайте Caddy](https://caddyserver.com/))

*TODO: а возможно только что-то одно. :)*

### Запуск напрямую

> [!IMPORTANT]
> Этот способ у вас может не заработать по ряду причин. Лучше используйте запуск в Docker Compose.
> Траблшутинг лежит на вашей ответственности.

> [!NOTE]
> Если вы используйте Linux, то сервер у вас не запустится на порту 80!
> Чтобы решить проблему, запускайте сервер с root правами (от пользователя root или через sudo/doas).
> Либо можете просто сменить порт в конфиге сервера и агента
> (`/backend/internal/config/config.go` и `/agent/config/config.go`).
> Но в таком случае запросы нужно тоже отправлять на другой порт.

* Linux

```bash
git clone https://github.com/wavy-cat/DAEC.git
cd DAEC
chmod +x setup.sh
./setup.sh
```

* Windows

```bash
git clone https://github.com/wavy-cat/DAEC.git
cd DAEC
setup.bat
```

*мне лень это делать, оставлю пока как черновик*

# Черновик

### Унарный минус

Это знак "-" перед которым стоит не цифра (пробел, скобка (но не закрывающая) или другой оператор 
либо если он просто находится в самом начале выражения) и после которого идёт цифра:

* `40-2` — от 40 отнять 2 = 38
* `-40-2` — от -40 отнять 2 = -42
* `- 40-2` — неправильная форма записи
* `40- 2` — от 40 отнять 2 = 38
* `40 - ( -2 )` — от 40 отнять -2 = 42
* `40-(-2)` — от 40 отнять -2 = 42

### Чё поддерживается

* Целочисленные числа
* Вещественные числа
* Унарный минус
* Операторы +, -, *, /, ^
* Скобки

### Структура проекта

### API Reference
