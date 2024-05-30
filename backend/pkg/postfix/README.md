### Это пакет

который реализует преобразование инфиксной записи в постфиксную aka обратную польскую запись.

Использует стек из пакета pkg/stack (backend)

> [!WARNING]
> Код может быть сложным для понимания, но он рабочий.
> Не лезьте в него без крайней необходимости.

Основные инструменты:

* Токенизатор инфиксной записи - <kbd>tokenizer.go</kbd>
* Конвертор инфиксной записи в постфиксную - <kbd>convertor.go</kbd>
* Декомпозитор постфиксной записи - <kbd>postfix.go</kbd>

Материалы:

* https://ru.wikipedia.org/wiki/%D0%90%D0%BB%D0%B3%D0%BE%D1%80%D0%B8%D1%82%D0%BC_%D1%81%D0%BE%D1%80%D1%82%D0%B8%D1%80%D0%BE%D0%B2%D0%BE%D1%87%D0%BD%D0%BE%D0%B9_%D1%81%D1%82%D0%B0%D0%BD%D1%86%D0%B8%D0%B8
* https://habr.com/ru/articles/596925/
* https://foxford.ru/wiki/informatika/obratnaya-polskaya-notatsiya
