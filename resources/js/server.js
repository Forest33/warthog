export {
  currentServer,
  currentService,
  currentServices,
  currentMethod,
  currentQuery,
  protoTypeEnum,
  protoTypeBool,
  protoTypeMessage,
  loadServer,
  saveRequest,
  createRequestForm,
  setRequestTitle,
  setCurrentServer,
  setCurrentQuery,
};
import { isNull } from "./index.js";
import {
  getRequestData,
  getRequestMetadata,
  hideQueryError,
  showQueryError,
} from "./request.js";
import { WorkspaceTypeQuery } from "./tree.js";
import { template } from "./template.js";

let currentServer = undefined;
let currentQuery = undefined;
let currentService = undefined;
let currentMethod = undefined;
let currentServices = {};
let currentRequest = {};
let currentInputCount = 1;
let oneOfNodes = new Map();

const protoTypeEnum = "enum";
const protoTypeBool = "bool";
const protoTypeMessage = "message";

const protoType2inputType = {
  string: "text",
  bytes: "text",
  int32: "number",
  int64: "number",
  uint32: "number",
  uint64: "number",
  sint32: "number",
  sint64: "number",
  fixed32: "number",
  fixed64: "number",
  sfixed32: "number",
  sfixed64: "number",
  double: "number",
  float: "number",
  bool: "checkbox",
  message: "checkbox",
};

function loadServer(srv, show) {
  currentQuery = undefined;
  currentRequest = {};

  let req = {
    name: "server.load",
    payload: {
      id: srv.id,
    },
  };

  hideQueryError();
  clearRequestPanel();

  astilectron.sendMessage(req, function (message) {
    if (message.payload.status !== "ok") {
      showQueryError(message.payload.error);
      currentServer = undefined;
      return;
    }

    if (!isNull(message.payload.data.warning)) {
      showServerWarning(message.payload.data.warning);
    }

    currentServer = message.payload.data.server;
    if (srv.type === WorkspaceTypeQuery) {
      currentQuery = message.payload.data.query;
      currentRequest = currentQuery.data.request;
      show = {
        service: { name: currentQuery.data.service },
        method: { name: currentQuery.data.method },
      };
    }

    let selServices = $("#sidebar-services-list");
    let selMethods = $("#sidebar-methods-list");
    selServices.children().remove();
    selMethods.children().remove();

    if (message.payload.data.services.length > 0) {
      let showServiceIdx = 0;
      message.payload.data.services.forEach(function (service, s_idx) {
        if (show !== undefined && show.service.name === service.name) {
          showServiceIdx = s_idx;
        }
      });
      message.payload.data.services.forEach(function (service, s_idx) {
        selServices.append(
          $("<option>", {
            value: service.name,
            text: service.name,
            selected: showServiceIdx === s_idx,
          })
        );

        currentServices[service.name] = {
          name: service.name,
          methods: {},
        };

        if (!isNull(service.methods)) {
          service.methods.forEach(function (method, m_idx) {
            currentServices[service.name].methods[method.name] = method;
            if (showServiceIdx === s_idx) {
              selMethods.append(
                $("<option>", {
                  value: method.name,
                  text: method.name,
                  selected:
                    (show !== undefined && show.method.name === method.name) ||
                    (show === undefined && m_idx === 0),
                })
              );
            }
          });
        }
      });

      if (
        show === undefined &&
        !isNull(message.payload.data.services[0].methods) &&
        message.payload.data.services[0].methods.length > 0
      ) {
        createRequestForm(
          message.payload.data.services[0],
          message.payload.data.services[0].methods[0]
        );
        return;
      }

      if (
        !isNull(show) &&
        currentServices[show.service.name] !== undefined &&
        !isNull(currentServices[show.service.name].methods) &&
        currentServices[show.service.name].methods[show.method.name] !==
          undefined
      ) {
        createRequestForm(
          currentServices[show.service.name],
          currentServices[show.service.name].methods[show.method.name]
        );
      } else if (
        message.payload.length > 0 &&
        !isNull(message.payload[0].methods) &&
        message.payload[0].methods.length > 0
      ) {
        createRequestForm(
          message.payload.data.services[0],
          message.payload.data.services[0].methods[0]
        );
      }
    }
  });

  $("#offcanvasTree").offcanvas("hide");
}

function createRequestForm(service, method) {
  currentService = service;
  currentMethod = method;
  oneOfNodes.clear();

  if (
    isNull(currentQuery) &&
    !isNull(currentServer.data.request) &&
    !isNull(currentServer.data.request[service.name]) &&
    !isNull(currentServer.data.request[service.name][method.name])
  ) {
    currentRequest = currentServer.data.request[service.name][method.name];
  }

  setRequestTitle(service.name + "." + method.name);
  enableRequestPanel();

  let request = $("#nav-request-data");
  request.html("").hide();

  if (method.input === undefined) {
    return;
  }

  let fields = [];
  for (const field of method.input) {
    let tmpl = getFieldTemplate(field, {});
    if (tmpl !== undefined) {
      request.append(tmpl);
      fields.push({ field: field, tmpl: tmpl });
    }
  }

  if (!isNull(currentRequest) && !isNull(currentRequest.input)) {
    for (const f of fields) {
      if (currentRequest.input[f.field.fqn] !== undefined) {
        setRequestData(f.field, f.tmpl, currentRequest.input[f.field.fqn]);
      }
    }
  }

  setRequestMetadata(currentRequest);

  request.show();
}

function getFieldTemplate(field, attr, showOneOf) {
  let tmpl;

  if (field.oneof !== undefined && !showOneOf) {
    let oneOfFqn = field.oneof.fqn + "-" + currentInputCount;
    if (!oneOfNodes.has(oneOfFqn)) {
      tmpl = $(template["request-oneof-select"]);
      let container = $(tmpl.find(".request-message-container")[0]);
      let select = $(tmpl.find(".field-value")[0]);
      select.attr("data-oneof-fqn", field.oneof.fqn);

      container.attr("data-container-fqn", field.oneof.fqn);
      $(tmpl.find(".label-name")[0]).html(field.oneof.name);
      $(tmpl.find(".label-type")[0]).html("oneof");
      select.append($("<option>", { value: "", text: "", selected: true }));

      oneOfNodes.set(oneOfFqn, {
        tmpl: tmpl,
        container: container,
        select: select,
        choices: new Map(),
      });

      select.on("change", function () {
        container.html("");
        if (this.value === "") {
          return;
        }
        let oneOf = oneOfNodes.get(oneOfFqn);
        let input = getFieldTemplate(oneOf.choices.get(this.value), {}, true);
        oneOf.container.append(input);
        oneOf.container.show();
      });
    }

    let oneOf = oneOfNodes.get(oneOfFqn);
    oneOf.choices.set(field.fqn, field);
    oneOf.select.append($("<option>", { value: field.fqn, text: field.name }));

    return oneOf.tmpl;
  }

  switch (field.type) {
    case protoTypeEnum:
      if (!field.repeated) {
        tmpl = $(template["request-select-input"]);
        $(tmpl.find(".label-name")[0]).html(field.name);
        $(tmpl.find(".label-type")[0]).html(field.enum.value_type);
        let sel = $(tmpl.find(".field-value")[0]);
        field.enum.values.forEach(function (item, idx) {
          sel.append(
            $("<option>", {
              value: item.number,
              text: item.name,
              selected: idx === 0,
            })
          );
        });
      } else {
        tmpl = getRepeatedInput(
          field,
          "request-select-input",
          field.enum.value_type,
          field.name,
          currentInputCount
        );
      }
      break;
    case protoTypeBool:
      if (!field.repeated) {
        tmpl = $(template["request-bool-input"]);
        $(tmpl.find(".label-name")[0]).attr("for", field.fqn).html(field.name);
        $(tmpl.find(".label-type")[0]).html(field.type);
      } else {
        tmpl = getRepeatedInput(
          field,
          "request-bool-input",
          field.type,
          field.name,
          currentInputCount
        );
      }
      break;
    case protoTypeMessage:
      let fieldType = field.type;
      if (field.map !== undefined) {
        fieldType =
          "map&lt;" + field.map.key_type + "," + field.map.value_type + "&gt;";
      } else if (field.message !== undefined) {
        fieldType = field.message.name;
      }

      tmpl = $(template["request-message-input"]);
      let container = $(tmpl.find(".request-message-container")[0]);
      let btn = $(tmpl.find(".request-message-button-add")[0]);
      container.attr("data-container-fqn", field.fqn);
      $(tmpl.find(".label-name")[0]).attr("for", field.fqn).html(field.name);
      $(tmpl.find(".label-type")[0]).html(fieldType);
      btn.attr("data-button-fqn", field.fqn);
      btn.click(function () {
        let inputCount = parseInt($(this).attr("data-input-count"), 10);
        $(this).attr("data-input-count", inputCount + 1);
        container.attr("data-container-input-count", inputCount + 1);
        let buttonAdd = $(this);

        let inputDelete = $(template["request-message-input-delete"]);
        let attr = {
          "data-input-fqn": field.fqn,
          "data-input-id": currentInputCount,
        };
        for (let key in attr) {
          inputDelete.attr(key, attr[key]);
        }

        $(inputDelete.find(".request-message-button-delete")[0]).click(
          function () {
            let id = $(this).parent().attr("data-input-id");
            let inputCount =
              parseInt(buttonAdd.attr("data-input-count"), 10) - 1;
            $(this)
              .closest(".request-message-container")
              .attr("data-container-input-count", inputCount);
            $('.request-input[data-input-id="' + id + '"]').remove();
            buttonAdd.attr("data-input-count", inputCount);
          }
        );

        if (field.map !== undefined) {
          container.append(inputDelete);
          let repeated = $(
            '<div class="request-repeated-message request-input" data-input-id="' +
              currentInputCount +
              '" data-repeated-fqn="' +
              field.fqn +
              '">'
          );
          if (field.map.fields === undefined) {
            let map = [
              {
                type: field.map.key_type,
                name: "key",
                "data-map-key": field.fqn,
              },
              {
                type: field.map.value_type,
                name: "value",
                "data-map-value": field.fqn,
                fqn: field.fqn + "-map-value",
              },
            ];
            for (const f of map) {
              repeated.append(getFieldTemplate(f, attr));
            }
          } else {
            repeated.append(
              getFieldTemplate(
                {
                  type: field.map.key_type,
                  name: "key",
                  "data-map-key": field.fqn,
                },
                attr
              )
            );
            for (const f of field.map.fields) {
              repeated.append(getFieldTemplate(f, attr));
            }
          }
          container.append(repeated);
        } else if (field.message !== undefined) {
          if (!field.repeated && inputCount > 0) {
            return;
          }
          container.append(inputDelete);
          let repeated = $(
            '<div class="request-repeated-message request-input" data-input-id="' +
              currentInputCount +
              '" data-repeated-fqn="' +
              field.fqn +
              '">'
          );
          for (const f of field.message.fields) {
            repeated.append(getFieldTemplate(f, attr));
          }
          container.append(repeated);
        }
        container.show();
        currentInputCount++;
      });
      break;
    default:
      if (!field.repeated) {
        tmpl = $(template["request-text-input"]);
        $(tmpl.find(".label-name")[0]).html(field.name);
        $(tmpl.find(".label-type")[0]).html(field.type);
      } else {
        tmpl = getRepeatedInput(
          field,
          "request-text-input",
          field.type,
          field.name,
          currentInputCount
        );
      }
  }

  let fv = $(tmpl.find(".field-value")[0]);
  fv.attr("data-field-name", field.name);
  fv.attr("data-field-type", field.type);
  fv.attr("data-field-fqn", field.fqn);
  fv.attr("id", field.fqn);
  if (field["data-map-key"] !== undefined) {
    fv.attr("data-map-key", field["data-map-key"]);
  }
  if (field["data-map-value"] !== undefined) {
    fv.attr("data-map-value", field["data-map-value"]);
  }

  if (attr !== undefined) {
    for (let key in attr) {
      tmpl.attr(key, attr[key]);
    }
  }

  let inputType = getInputType(field.type);
  if (inputType !== "") {
    fv.attr("type", inputType);
  }

  return tmpl;
}

function getRepeatedInput(field, inputTmpl, inputType, inputName, dataID) {
  let tmpl = $(template["request-message-input"]);
  let container = $(tmpl.find(".request-message-container")[0]);
  let btn = $(tmpl.find(".request-message-button-add")[0]);
  container.attr("data-container-fqn", field.fqn);
  container.attr("data-parent-id", dataID);
  $(tmpl.find(".label-name")[0]).attr("for", field.fqn).html(field.name);
  $(tmpl.find(".label-type")[0]).html(field.type);
  btn.attr("data-button-fqn", field.fqn);
  btn.click(function () {
    let inputCount = parseInt($(this).attr("data-input-count"), 10);
    $(this).attr("data-input-count", inputCount + 1);
    container.attr("data-container-input-count", inputCount + 1);
    let buttonAdd = $(this);

    let inputDelete = $(template["request-message-input-delete"]);
    let attr = {
      "data-input-fqn": field.fqn,
      "data-input-id": currentInputCount,
    };
    for (let key in attr) {
      inputDelete.attr(key, attr[key]);
    }
    container.append(inputDelete);

    $(inputDelete.find(".request-message-button-delete")[0]).click(function () {
      let id = $(this).parent().attr("data-input-id");
      let inputCount = parseInt(buttonAdd.attr("data-input-count"), 10) - 1;
      $(this)
        .closest(".request-message-container")
        .attr("data-container-input-count", inputCount);
      $('.request-input[data-input-id="' + id + '"]').remove();
      buttonAdd.attr(
        "data-input-count",
        parseInt(buttonAdd.attr("data-input-count"), 10) - 1
      );
    });

    let input = $(template[inputTmpl]);
    $(input.find(".label-name")[0])
      .attr("for", field.fqn + "." + currentInputCount)
      .html(inputName);
    $(input.find(".label-type")[0]).html(inputType);
    $(input.find(".field-value")[0]).attr("data-field-fqn", field.fqn);
    $(input.find(".field-value")[0]).attr("data-parent-id", dataID);

    if (field.enum !== undefined) {
      let sel = $(input.find(".field-value")[0]);
      field.enum.values.forEach(function (item, idx) {
        sel.append(
          $("<option>", {
            value: item.number,
            text: item.name,
            selected: idx === 0,
          })
        );
      });
    }

    $(input.find(".field-value")[0]).attr(
      "id",
      field.fqn + "." + currentInputCount
    );
    for (let key in attr) {
      input.attr(key, attr[key]);
    }

    container.append(input);
    container.show();
    currentInputCount++;
  });
  return tmpl;
}

function getInputType(type) {
  let t = protoType2inputType[type];
  if (t === undefined) {
    return "";
  }
  return t;
}

function setRequestTitle(title) {
  let pageTitle = "Warthog";
  if (!isNull(currentServer) && title !== "") {
    pageTitle = currentServer.breadcrumb.join(" » ") + " » " + title;
  }
  document.title = pageTitle;
  $("#nav-request-data-title").html(title);
}

function enableRequestPanel() {
  $("#workspace-reload").removeClass("disabled");
  $("#request-run").removeClass("disabled");
  $("#sidebar-services-list").removeAttr("disabled");
  $("#sidebar-methods-list").removeAttr("disabled");
  $("#request-edit-buttons").show();
}

function clearRequestPanel() {
  $("#workspace-reload").addClass("disabled");
  $("#request-run").addClass("disabled");
  $("#sidebar-services-list").attr("disabled", true);
  $("#sidebar-methods-list").attr("disabled", true);
  $("#request-edit-buttons").hide();
  $("#nav-request-data").html("");
  $("#sidebar-services-list option").each(function () {
    $(this).remove();
  });
  $("#sidebar-methods-list option").each(function () {
    $(this).remove();
  });
  setRequestTitle("");
}

function saveRequest() {
  if (
    !isNull(currentQuery) ||
    isNull(currentServer) ||
    isNull(currentService) ||
    isNull(currentMethod)
  ) {
    return;
  }

  let request = {};
  if (currentMethod.input !== undefined) {
    for (const field of currentMethod.input) {
      let data = getRequestData(field, undefined, true);
      $.extend(request, data);
    }
  }

  if (
    isNull(currentServer.data.request) ||
    isNull(currentServer.data.request)
  ) {
    currentServer.data.request = {};
  }
  if (
    isNull(currentServer.data.request[currentService.name]) ||
    isNull(currentServer.data.request[currentService.name])
  ) {
    currentServer.data.request[currentService.name] = {};
  }

  let metadata = [];
  for (let [key, value] of Object.entries(getRequestMetadata())) {
    metadata.push({ key: key, value: value });
  }

  currentServer.data.request[currentService.name][currentMethod.name] = {
    input: request,
    metadata: metadata,
  };

  let req = {
    name: "server.update.request",
    payload: {
      id: currentServer.id,
      service: currentService.name,
      method: currentMethod.name,
      request:
        currentServer.data.request[currentService.name][currentMethod.name],
    },
  };

  console.log(req);

  astilectron.sendMessage(req, function () {});
}

function setRequestData(field, tmpl, data) {
  if (field.oneof !== undefined) {
    tmpl
      .find('.field-value[data-oneof-fqn="' + field.oneof.fqn + '"]')
      .val(field.fqn)
      .change();
  }

  switch (field.type) {
    case protoTypeEnum:
      if (!field.repeated) {
        $(
          tmpl.find(
            '.field-value[data-field-fqn="' +
              field.fqn +
              '"] option[value="' +
              data +
              '"]'
          )[0]
        ).attr("selected", "selected");
      } else {
        let btn = tmpl.find(
          '.request-message-button-add[data-button-fqn="' + field.fqn + '"]'
        );
        data.forEach(function (d, idx) {
          btn.trigger("click");
          $(
            tmpl.find(
              '.field-value[data-field-fqn="' +
                field.fqn +
                '"] option[value="' +
                d +
                '"]'
            )[idx]
          ).attr("selected", "selected");
        });
      }
      break;
    case protoTypeBool:
      if (!field.repeated) {
        $(
          tmpl.find('.field-value[data-field-fqn="' + field.fqn + '"]')[0]
        ).prop("checked", data);
      } else {
        let btn = tmpl.find(
          '.request-message-button-add[data-button-fqn="' + field.fqn + '"]'
        );
        data.forEach(function (d, idx) {
          btn.trigger("click");
          $(
            tmpl.find('.field-value[data-field-fqn="' + field.fqn + '"]')[idx]
          ).prop("checked", d);
        });
      }
      break;
    case protoTypeMessage:
      let btn = tmpl.find(
        '.request-message-button-add[data-button-fqn="' + field.fqn + '"]'
      );
      if (field.map !== undefined) {
        if (field.map.fields === undefined) {
          let idx = 0;
          for (let key in data) {
            btn.trigger("click");
            let container = $(
              tmpl.find(
                '.request-repeated-message[data-repeated-fqn="' +
                  field.fqn +
                  '"]'
              )[idx++]
            );
            container
              .find('.field-value[data-map-key="' + field.fqn + '"]')
              .val(key);
            let f = {
              fqn: field.fqn + "-map-value",
              type: field.map.value_type,
            };
            setRequestData(f, container, data[key]);
          }
        } else {
          let idx = 0;
          for (let key in data) {
            btn.trigger("click");
            let container = $(
              tmpl.find(
                '.request-repeated-message[data-repeated-fqn="' +
                  field.fqn +
                  '"]'
              )[idx++]
            );
            container
              .find('.field-value[data-map-key="' + field.fqn + '"]')
              .val(key);
            for (const f of field.map.fields) {
              if (data[key][f.fqn] !== undefined) {
                setRequestData(f, container, data[key][f.fqn]);
              }
            }
          }
        }
      } else if (field.message !== undefined) {
        if (!field.repeated) {
          btn.trigger("click");
          for (const f of field.message.fields) {
            if (data[f.fqn] !== undefined) {
              setRequestData(f, tmpl, data[f.fqn]);
            }
          }
        } else {
          data.forEach(function (d, idx) {
            btn.trigger("click");
            let container = $(
              tmpl.find(
                '.request-repeated-message[data-repeated-fqn="' +
                  field.fqn +
                  '"]'
              )[idx]
            );
            for (const f of field.message.fields) {
              if (d[f.fqn] !== undefined) {
                setRequestData(f, container, d[f.fqn]);
              }
            }
          });
        }
      }
      break;
    default:
      if (!field.repeated) {
        $(tmpl.find('.field-value[data-field-fqn="' + field.fqn + '"]')[0]).val(
          data
        );
      } else {
        let btn = tmpl.find(
          '.request-message-button-add[data-button-fqn="' + field.fqn + '"]'
        );
        data.forEach(function (d, idx) {
          btn.trigger("click");
          $(
            tmpl.find('.field-value[data-field-fqn="' + field.fqn + '"]')[idx]
          ).val(d);
        });
      }
  }
}

function setRequestMetadata(request) {
  let metadata = $("#nav-request-metadata");

  metadata.find(".metadata-row").each(function (i, row) {
    if (i > 0) {
      $(row).remove();
    } else {
      $(row).find(".metadata-key").val("");
      $(row).find(".metadata-value").val("");
    }
  });

  if (isNull(request) || isNull(request.metadata)) {
    return;
  }

  request.metadata.forEach(function (d, i) {
    let lastKey = metadata.find(".metadata-key:last");
    lastKey.val(d.key);
    metadata.find(".metadata-value:last").val(d.value);
    if (i + 1 < request.metadata.length) {
      lastKey.trigger("click");
    }
  });
}

function showServerWarning(warn) {
  let error = $("#query-error");
  for (const w of warn) {
    let tmpl = $(template["protobuf-error"]);
    tmpl.addClass("alert-warning");
    $(tmpl.find(".file")[0]).html(w.pos.Filename);
    $(tmpl.find(".line")[0]).html(w.pos.Line);
    $(tmpl.find(".column")[0]).html(w.pos.Col);
    $(tmpl.find(".message")[0]).html(w.warning);
    error.append(tmpl);
  }
  error.show();
}

function setCurrentServer(s) {
  currentServer = s;
}

function setCurrentQuery(q) {
  currentQuery = q;
}
