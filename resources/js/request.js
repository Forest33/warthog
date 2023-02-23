export {
    initStreamControl,
    hideQueryError,
    hideStreamControl,
    query,
    response,
    getRequestData,
    getRequestMetadata,
    showQueryError,
};
import {currentMethod, currentQuery, currentServer, currentService, protoTypeBool, protoTypeBytes, protoTypeEnum, protoTypeMessage, saveRequest,} from "./server.js";
import {isNull} from "./index.js";
import {template} from "./template.js";

const MethodTypeUnary = "u";
const MethodTypeClientStream = "cs";
const MethodTypeServerStream = "ss";
const MethodTypeBidiStream = "css";

let streamStopped = false;

function query() {
    if (currentService === undefined || currentMethod === undefined) {
        return;
    }

    if (isQueryRun()) {
        astilectron.sendMessage({name: "query.cancel"}, function () {
        });
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
            server_id: currentServer.id,
            service: currentService.name,
            method: currentMethod.name,
            metadata: getRequestMetadata(),
            data: request,
        },
    };

    if (currentMethod.type === MethodTypeUnary) {
        setQueryCancelButton();
    } else {
        showStreamControl();
        streamStopped = false;
    }

    if (isNull(currentQuery)) {
        saveRequest();
    }

    astilectron.sendMessage(req, function (message) {
        if (currentMethod.type === MethodTypeUnary) {
            setQueryRunButton();
        }
        if (message.payload.status !== "ok") {
            showQueryError(message.payload.error);
            return
        }
        $("#stream-sent .sent").html(message.payload.data.sent);
    });
}

function initStreamControl() {
    $("#stream-stop").on("click", function () {
        streamStopped = true;
        astilectron.sendMessage({name: "query.close.stream"}, function () {
        });
        hideStreamControl();
    });

    $("#stream-cancel").on("click", function () {
        streamStopped = true;
        astilectron.sendMessage({name: "query.cancel"}, function () {
        });
        hideStreamControl();
    });
}

function showStreamControl() {
    switch (currentMethod.type) {
        case MethodTypeClientStream:
            $("#stream-sent").show();
            $("#stream-stop").html("Receive & Close").css("visibility", "visible");
            break;
        case MethodTypeServerStream:
            $("#stream-received").show();
            break;
        case MethodTypeBidiStream:
            $("#stream-sent").show();
            $("#stream-received").show();
            $("#stream-stop").html("Send & Close").css("visibility", "visible");
            break;
    }

    $("#stream-info").show();
    $("#stream-control").css("visibility", "visible");
    $("#time-spent").html("");
}

function hideStreamControl() {
    $("#stream-info").hide();
    $("#stream-info .direction").hide();
    $("#stream-stop").css("visibility", "hidden");
    $("#stream-control").css("visibility", "hidden");
}

function response(data) {
    if (isNull(data.error)) {
        $("#stream-received .received").html(data.received);
        $("#query-error").html("");
        if (
            currentMethod.type === MethodTypeUnary ||
            currentMethod.type === MethodTypeClientStream ||
            streamStopped
        ) {
            $("#badge-result")
                .html("0: OK")
                .css("visibility", "visible")
                .removeClass("bg-danger")
                .addClass("bg-success");
            $("#time-spent").html(data.spent_time);
            if (data.json_string !== "") {
                $("#query-result").html(syntaxHighlight(data.json_string));
            }
            streamStopped = false;
        } else if (data.json_string !== "") {
            let sep = '<div class="hr"><span>' + data.time + "</span></div>";
            if (data.received > 1) {
                $("#query-result").prepend(sep + syntaxHighlight(data.json_string));
            } else {
                $("#query-result").html(sep + syntaxHighlight(data.json_string));
            }
        }
    } else {
        if (
            currentMethod.type === MethodTypeUnary ||
            currentMethod.type === MethodTypeClientStream
        ) {
            $("#query-result").html("");
        }
        $("#badge-result")
            .html(data.error.code + ": " + data.error.code_description)
            .css("visibility", "visible")
            .removeClass("bg-success")
            .addClass("bg-danger");
        let tmpl = $(template["query-error"]);
        $(tmpl.find(".code")[0]).html(data.error.code);
        $(tmpl.find(".message")[0]).html(data.error.message);
        $("#query-error").html(tmpl).show();
        $("#time-spent").html(data.spent_time);
        hideStreamControl();
    }

    showHeadersTrailers(data.header, data.trailer);
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
        case protoTypeBytes:
            let bytesValues = [];
            root.find('[data-field-fqn="' + field.fqn + '"]').each(function () {
                let val = {
                    value: "",
                    file: "",
                }
                if ($(this).prop("disabled")) {
                    val.file = $(this).val();
                } else {
                    val.value = $(this).val();
                }
                bytesValues.push(val);
            });
            if (bytesValues.length > 0) {
                if (!field.repeated) {
                    data[field.fqn] = bytesValues[0];
                } else {
                    data[field.fqn] = bytesValues;
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
    let tmpl;
    if (!isNull(err.pos) && err.pos.Filename !== "") {
        tmpl = $(template["protobuf-error"]);
        $(tmpl.find(".alert")[0]).addClass("alert-danger");
        $(tmpl.find(".file")[0]).html(err.pos.Filename);
        $(tmpl.find(".line")[0]).html(err.pos.Line);
        $(tmpl.find(".column")[0]).html(err.pos.Col);
    } else {
        tmpl = $(template["query-error"]);
        $(tmpl.find(".code")[0]).html(err.code);
    }

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
    $("#time-spent").html("");
    let badge = $("#badge-result");
    let error = $("#query-error");
    if (badge.hasClass("bg-danger") || error.find(".alert-warning")) {
        error.html("").hide();
        badge.removeClass("bg-danger").css("visibility", "hidden");
    }
}

function showHeadersTrailers(header, trailer) {
    let headers = $("#nav-result-headers");
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
    console.log(currentMethod);
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
