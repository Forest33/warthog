import {currentMethod, currentQuery, currentServer, currentService, setCurrentQuery, setCurrentServer, setRequestTitle,} from "./server.js";
import {getRequestData, getRequestMetadata} from "./request.js";
import {isNull} from "./index.js";

export {saveQuery};

function saveQuery() {
    if (
        isNull(currentServer) ||
        isNull(currentService) ||
        isNull(currentMethod)
    ) {
        return;
    }

    let request = {};
    for (const field of currentMethod.input) {
        let data = getRequestData(field, undefined, true);
        $.extend(request, data);
    }

    let metadata = [];
    for (let [key, value] of Object.entries(getRequestMetadata())) {
        metadata.push({key: key, value: value});
    }

    let req = {
        name: "query.update",
        payload: {
            server_id: currentServer.id,
            service: currentService.name,
            method: currentMethod.name,
            title: $(".popover-body .query-popover-query-name").val(),
            description: $(".popover-body .query-popover-query-description").val(),
            request: {
                input: request,
                metadata: metadata,
            },
        },
    };
    if (!isNull(currentQuery)) {
        req.payload.id = currentQuery.id;
    }

    astilectron.sendMessage(req, function (message) {
        if (message.payload.status !== "ok") {
            return;
        }
        setCurrentServer(message.payload.data.server);
        setCurrentQuery(message.payload.data.query);
        setRequestTitle(currentQuery.data.service + "." + currentQuery.data.method);
    });
}
