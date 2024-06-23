let statusLocale = new Map();
statusLocale.set('pending', '<span class="badge rounded-pill text-bg-warning">Вычисляется</span>');
statusLocale.set('done', '<span class="badge rounded-pill text-bg-success">Завершено</span>');
statusLocale.set('error', '<span class="badge rounded-pill text-bg-danger">Ошибка</span>');

// Подготавливает таблицу
async function prepareTable() {
    // Удаляем текст "Загрузка задач..."
    const p = document.querySelector('#no-exps');
    if (p) {
        p.remove();
    }

    // Делаем таблицу видимой
    const elem = document.querySelector('#table');
    if (elem.classList.contains('invisible')) {
        elem.classList.remove('invisible');
    }
}

// Собирает HTML текст из данных
async function taskContentBuilder(id, status, result, expression) {
    if (result == null) {
        result = `<span class="text-body-secondary">¯\\_(ツ)_/¯</span>`;
    } else if (status === "error") {
        result = `<span class="text-body-secondary">Безрезультатно</span>`;
    }

    if (expression == null) {
        expression = `<span class="text-body-secondary">¯\\_(ツ)_/¯</span>`;
    }

    return `<tr id="${id}">
        <td><code>${id}</code></td>
        <td>${statusLocale.get(status)}</td>
        <td>${expression}</td>
        <td>${result}</td>
    </tr>`;
}

// Изменяет задачу в таблице по ID
async function editTask(id, status, result, expression = null) {
    const row = document.getElementById(id);
    if (row) {
        row.innerHTML = await taskContentBuilder(id, status, result, expression);
    }
}

// Добавляет задачу в таблицу
async function addTask(id, status, result, expression = null) {
    await prepareTable()
    const thead = document.getElementById('tbody');
    thead.insertAdjacentHTML('afterbegin', await taskContentBuilder(id, status, result, expression));
}