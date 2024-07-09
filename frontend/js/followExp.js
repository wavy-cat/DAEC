async function followExp(expressionId) {
    const {id, status, result, content} = await getExpressionByID(expressionId)
    if (status !== 'pending') {
        const row = document.getElementById(expressionId);
        if (!row) return;
        return row.innerHTML = await taskContentBuilder(id, status, result, content);
    }
    await new Promise(resolve => setTimeout(resolve, 1000));
    await followExp(id);
}