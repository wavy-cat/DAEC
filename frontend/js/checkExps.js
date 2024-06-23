async function getData() {
    const currentTimeElement = document.getElementById("no-exps");

    await fetch('http://localhost/api/v1/expressions')
        .then(response => response.json())
        .then(data => {
            if (data.expressions.length === 0) {
                currentTimeElement.innerHTML = "Вы ещё не отсылали задач";
            } else {
                data.expressions.forEach(async function (item, _) {
                    await addTask(item.id, item.status, item.result)
                })
            }
        })
        .catch(error => console.error('Error:', error));
}

document.addEventListener("DOMContentLoaded", function () {
    getData().then();
});