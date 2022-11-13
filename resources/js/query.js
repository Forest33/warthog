export {saveQuery};
import {
  currentMethod,
  currentQuery,
  currentServer,
  currentService,
  setCurrentQuery,
  setCurrentServer,
  setRequestTitle,
} from "./server.js";
import {getRequestData} from "./request.js";
import {isNull} from "./index.js";

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

  let req = {
    name: "query.update",
    payload: {
      server_id: currentServer.id,
      service: currentService.name,
      method: currentMethod.name,
      request: request,
      title: $(".popover-body .query-popover-query-name").val(),
      description: $(".popover-body .query-popover-query-description").val(),
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
