async function getData() {
    const currentTimeElement = document.getElementById("no-exps");
    const response = await fetch('http://localhost/api/v1/expressions')
    const data = await response.json();
    const tbody = document.getElementById('tbody');

    const expressions = data?.expressions
    if (!expressions || expressions.length === 0) {
        return currentTimeElement.innerHTML = "Вы ещё не отсылали задач";
    }
    
    expressions.map(async ({id, status, result, content}) => {
        tbody.insertAdjacentHTML('afterbegin', await taskContentBuilder(id, status, result, content));
        if (status === 'pending') await followExp(id)
    })

    await prepareTable()
}

document.addEventListener("DOMContentLoaded", function () {
    getData().then();
});