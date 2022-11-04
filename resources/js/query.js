function saveQuery() {
    if (isNull(currentServer) || isNull(currentService) || isNull(currentMethod)) {
        return
    }

    let request = {}
    for (const field of currentMethod.input) {
        let data = getRequestData(field)
        $.extend(request, data)
    }

    let req = {
        name: 'query.update',
        payload: {
            server_id: currentServer.id,
            service: currentService.name,
            method: currentMethod.name,
            request: request,
            title: $('.popover-body .query-popover-query-name').val(),
            description: $('.popover-body .query-popover-query-description').val(),
        }
    }
    if (!isNull(currentQuery)) {
        req.payload.id = currentQuery.id
    }

    astilectron.sendMessage(req, function (message) {
        if (message.payload.status !== "ok") {
            return
        }
        currentQuery = message.payload.data.query
        currentServer = message.payload.data.server
        setRequestTitle(currentQuery.data.service + "." + currentQuery.data.method)
    });
}
