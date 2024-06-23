async function followExp(id, expression = null) {
    let result;
    do {
        await sleep(1000); // Пауза 1 секунда между запросами
        result = await getExpressionByGet(id);
    } while (result.status === "pending");

    // По готовности изменяем задачу в таблице
    await editTask(id, result.status, result.result, expression)
}