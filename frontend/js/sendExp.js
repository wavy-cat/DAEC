// Отправляет запрос оркестратору на получение выражения по его ID
async function getExpressionByID(id) {
    try {
        const response = await fetch(`${ServerAddress}/api/v1/expressions/${id}`);
        const result = await response.json();
        return result.expression
    } catch (e) {
        alert('Упс! Внутренняя ошибка: ' + e)
    }
}

// Отправляет выражение на вычисление
async function sendExpression() {
    // Получаем выражение из input
    const expression = document.getElementById("exp-data").value;

    // Отправляем запрос
    const reqData = {expression: expression};
    const response = await fetch(
        `${ServerAddress}/api/v1/calculate`,
        {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(reqData)
        });
    const data = await response.json();

    // Показываем ошибку, если есть
    if (!data.hasOwnProperty('id')) return alert("Ошибка в выражении: " + data.error);

    // Очищаем input
    document.getElementById('exp-data').value = '';

    // Добавляем задачу в таблицу
    const tbody = document.getElementById('tbody');
    tbody.insertAdjacentHTML('afterbegin', await taskContentBuilder(data.id, 'pending', null, expression));
    await prepareTable()

    // Следим за статусом задачи 👀
    await followExp(data.id, expression)
}
