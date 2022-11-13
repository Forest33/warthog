export { hideQueryError, query, getRequestData };
import {
  currentMethod,
  currentService,
  protoTypeBool,
  protoTypeEnum,
  protoTypeMessage,
} from "./server.js";
import { isNull } from "./index.js";
import { template } from "./template.js";

function query() {
  if (currentService === undefined || currentMethod === undefined) {
    return;
  }

  if (isQueryRun()) {
    astilectron.sendMessage({ name: "query.cancel" }, function () {});
    setQueryRunButton();
    return;
  }

  hideQueryError();

  let request = {};
  if (currentMethod.input !== undefined) {
    for (const field of currentMethod.input) {
      let data = getRequestData(field, undefined, false);
      $.extend(request, data);
    }
  }

  console.log("request: " + JSON.stringify(request, null, 1));

  let req = {
    name: "query.run",
    payload: {
      service: currentService.name,
      method: currentMethod.name,
      metadata: getRequestMetadata(),
      data: request,
    },
  };

  setQueryCancelButton();

  astilectron.sendMessage(req, function (message) {
    //console.log("query result: " + JSON.stringify(message, null, 1));
    setQueryRunButton();
    if (message.payload.status !== "ok") {
      showQueryError(message.payload.error);
      return;
    }
    $("#badge-result")
      .html("0: OK")
      .css("visibility", "visible")
      .addClass("bg-success");
    $("#time-spent").html(message.payload.data.spent_time);
    $("#query-result").html(syntaxHighlight(message.payload.data.json_string));
    showHeadersTrailers(
      message.payload.data.header,
      message.payload.data.trailer
    );
  });
}

function getRequestData(field, root, disableProtoFQN) {
  let data = {};
  if (root === undefined) {
    root = $("#nav-request-data");
  }

  switch (field.type) {
    case protoTypeEnum:
      let enumValues = [];
      root.find('[data-field-fqn="' + field.fqn + '"]').each(function () {
        enumValues.push($(this).find(":selected").val());
      });
      if (enumValues.length > 0) {
        if (!field.repeated) {
          data[field.fqn] = enumValues[0];
        } else {
          data[field.fqn] = enumValues;
        }
      }
      break;
    case protoTypeBool:
      let boolValues = [];
      root.find('[data-field-fqn="' + field.fqn + '"]').each(function () {
        boolValues.push($(this).is(":checked"));
      });
      if (boolValues.length > 0) {
        if (!field.repeated) {
          data[field.fqn] = boolValues[0];
        } else {
          data[field.fqn] = boolValues;
        }
      }
      break;
    case protoTypeMessage:
      if (field.map !== undefined) {
        // todo validate dup keys
        if (field.map.fields === undefined) {
          let items = {};
          root
            .find(
              '.request-repeated-message[data-repeated-fqn="' + field.fqn + '"]'
            )
            .each(function () {
              let key = $(this)
                .find('[data-map-key="' + field.fqn + '"]')
                .val();
              items[key] = $(this)
                .find('[data-map-value="' + field.fqn + '"]')
                .val();
            });
          if (Object.keys(items).length > 0) {
            data[getProtoFQN(field, disableProtoFQN)] = items;
          }
        } else {
          let objects = {};
          root
            .find(
              '.request-repeated-message[data-repeated-fqn="' + field.fqn + '"]'
            )
            .each(function () {
              let key = $(this)
                .find('[data-map-key="' + field.fqn + '"]')
                .val();
              let items = {};
              for (const f of field.map.fields) {
                let v = getRequestData(f, $(this), disableProtoFQN);
                if (Object.keys(v).length > 0) {
                  $.extend(items, v);
                }
              }
              if (Object.keys(items).length > 0) {
                objects[key] = items;
              }
            });
          if (Object.keys(objects).length > 0) {
            data[getProtoFQN(field, disableProtoFQN)] = objects;
          }
        }
      } else if (field.message !== undefined) {
        let messageValues = [];
        root
          .find(
            '.request-repeated-message[data-repeated-fqn="' + field.fqn + '"]'
          )
          .each(function () {
            let messageObjects = {};
            for (const f of field.message.fields) {
              let v = getRequestData(f, $(this), disableProtoFQN);
              if (Object.keys(v).length !== 0) {
                $.extend(messageObjects, v);
              }
            }
            messageValues.push(messageObjects);
          });
        if (messageValues.length > 0) {
          let fqn = getProtoFQN(field, disableProtoFQN);
          if (!field.repeated) {
            data[fqn] = messageValues[0];
          } else {
            data[fqn] = messageValues;
          }
        }
      }
      break;
    default:
      let textValues = [];
      root.find('[data-field-fqn="' + field.fqn + '"]').each(function () {
        textValues.push($(this).val());
      });
      if (textValues.length > 0) {
        if (!field.repeated) {
          data[field.fqn] = textValues[0];
        } else {
          data[field.fqn] = textValues;
        }
      }
  }

  return data;
}

function getProtoFQN(field, disableProtoFQN) {
  if (disableProtoFQN || isNull(field.proto_fqn) || field.proto_fqn === "") {
    return field.fqn;
  }
  return field.proto_fqn;
}

function getRequestMetadata() {
  let metadata = {};
  let keys = [];
  $("#nav-request-metadata .metadata-key").each(function () {
    keys.push($(this).val());
  });
  $("#nav-request-metadata .metadata-value").each(function (i) {
    if (keys[i] !== "") {
      metadata[keys[i]] = $(this).val();
    }
  });
  return metadata;
}

function syntaxHighlight(json) {
  json = json
    .replace(/&/g, "&amp;")
    .replace(/</g, "&lt;")
    .replace(/>/g, "&gt;");
  return json.replace(
    /("(\\u[a-zA-Z0-9]{4}|\\[^u]|[^\\"])*"(\s*:)?|\b(true|false|null)\b|-?\d+(?:\.\d*)?(?:[eE][+\-]?\d+)?)/g,
    function (match) {
      let cls = "number";
      if (/^"/.test(match)) {
        if (/:$/.test(match)) {
          cls = "key";
        } else {
          cls = "string";
        }
      } else if (/true|false/.test(match)) {
        cls = "boolean";
      } else if (/null/.test(match)) {
        cls = "null";
      }
      return '<span class="' + cls + '">' + match + "</span>";
    }
  );
}

function showQueryError(err) {
  let tmpl = $(template["query-error"]);
  $(tmpl.find(".code")[0]).html(err.code);
  $(tmpl.find(".message")[0]).html(err.message);
  $("#query-error").html(tmpl).show();
  $("#query-result").html("").hide();
  $("#time-spent").html("");
  $("#badge-result")
    .html(err.code + ": " + err.code_description)
    .css("visibility", "visible")
    .removeClass("bg-success")
    .addClass("bg-danger");
  showHeadersTrailers(null, null);
}

function hideQueryError() {
  $("#query-result").show();
  let badge = $("#badge-result");
  if (badge.hasClass("bg-danger")) {
    $("#query-error").html("").hide();
    badge.removeClass("bg-danger").css("visibility", "hidden");
  }
}

function showHeadersTrailers(header, trailer) {
  let headers = $("#nav-result-headers")
  headers.find(".header").css("visibility", "hidden");
  headers.find(".trailer").css("visibility", "hidden");
  if (header !== null && Object.keys(header).length !== 0) {
    $("#query-result-headers").html(
      syntaxHighlight(JSON.stringify(header, null, 1))
    );
    $("#nav-result-headers .header").css("visibility", "visible");
  }
  if (trailer !== null && Object.keys(trailer).length !== 0) {
    $("#query-result-trailers").html(
      syntaxHighlight(JSON.stringify(trailer, null, 1))
    );
    $("#nav-result-headers .trailer").css("visibility", "visible");
  }
}

function isQueryRun() {
  return $("#request-run").hasClass("btn-danger");
}

function setQueryRunButton() {
  $("#request-run")
    .removeClass("btn-danger bi-stop-fill")
    .addClass("btn-success bi-play-fill");
}

function setQueryCancelButton() {
  $("#request-run")
    .removeClass("btn-success bi-play-fill")
    .addClass("btn-danger bi-stop-fill");
}
