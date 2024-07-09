function sleep(ms) {
    return new Promise(resolve => setTimeout(resolve, ms));
}

// Отправляет запрос оркестратору на получение выражения по его ID
async function getExpressionByGet(id) {
    let result;
    await fetch(`${ServerAddress}/api/v1/expressions/${id}`)
        .then(response => response.json())
        .then(data => {
            result = data.expression
        })
        .catch(error => alert('Error: ' + error));
    return result
}

// Отправляет выражение на вычисление
async function sendExpression() {
    // Получаем выражение из input
    const expression = document.getElementById("exp-data").value;

    // Отправляем запрос
    let reqData = {expression: expression};
    let response = await fetch(
        `${ServerAddress}/api/v1/calculate`,
        {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(reqData)
        });
    let data = await response.json();

    // Показываем ошибку, если есть
    if (!data.hasOwnProperty('id')) {
        alert("Ошибка в выражении: " + data.error);
        return;
    }

    // Очищаем input
    document.getElementById('exp-data').value = '';

    // Добавляем задачу в таблицу
    await addTask(data.id, "pending", null, expression);

    // Следим за статусом задачи 👀
    await followExp(data.id, expression)
}
