export { modalTreeExpandedNodes, treeRootNodes, dataIdToNode, isNull };
import { saveQuery } from "./query.js";
import { showTree } from "./tree.js";
import { editServer, initWorkspaceModal } from "./workspace.modal.js";
import { query } from "./request.js";
import {
  createRequestForm,
  currentMethod,
  currentQuery,
  currentServer,
  currentService,
  currentServices,
  loadServer,
  saveRequest,
  setRequestTitle,
} from "./server.js";

let treeRootNodes = new Set();
let modalTreeExpandedNodes = new Set();
let dataIdToNode = new Map();

$(document).ready(function () {
  document.addEventListener("astilectron-ready", function () {
    astilectron.onMessage(function (message) {
      console.log("server message: " + JSON.stringify(message, null, 1));
      if (isNull(message.name)) {
        return;
      }
      switch (message.name) {
        case "server.load":
          loadServer({ id: message.payload.id });
          break;
        case "menu.about":
          showAbout(message.payload);
      }
    });

    astilectron.sendMessage({ name: "workspace.state" }, function (message) {
      if (message.payload.status !== "ok") {
        return;
      }
      if (message.payload.data.servers === 0) {
        $("#workspaceModal").modal("show");
        return;
      }
      if (
        isNull(message.payload.data.startup_workspace_id) ||
        message.payload.data.startup_workspace_id === 0
      ) {
        $("#offcanvasTree").offcanvas("show");
      }
    });
  });

  let includes = $("[data-include]");
  $.each(includes, function () {
    $(this).load($(this).data("include"), function () {
      switch ($(this).attr("data-include")) {
        case "modal.workspace.html":
          initWorkspaceModal();
          break;
      }
    });
  });

  initSidebars();
  initOffcanvas();
  initQueryPopover();

  const shell = require("electron").shell;
  $(document).on("click", 'a[href^="http"]', function (event) {
    event.preventDefault();
    shell.openExternal(this.href);
  });

  $("#workspace-reload").click(function () {
    if (currentService === undefined || currentMethod === undefined) return;
    saveRequest();
    if (!isNull(currentQuery)) {
      loadServer(currentQuery, {
        service: currentService,
        method: currentMethod,
      });
    } else {
      loadServer(currentServer, {
        service: currentService,
        method: currentMethod,
      });
    }
  });

  $("#request-run").click(function () {
    query();
  });

  $("#edit-server").click(function () {
    editServer(currentServer);
  });

  $("#nav-request form").submit(function () {
    query();
    return false;
  });

  $("#nav-request form").keypress(function (e) {
    if (e.which === 13) {
      $("#nav-request form").submit();
      return false;
    }
  });

  $("#sidebar-services-list").on("change", function () {
    saveRequest();
    let selMethods = $("#sidebar-methods-list");
    let service = currentServices[this.value];
    selMethods.children().remove();
    setRequestTitle("");

    if (isNull(service)) {
      return;
    }

    let first = true;
    for (let [name, method] of Object.entries(service.methods)) {
      selMethods.append(
        $("<option>", {
          value: name,
          text: name,
          selected: first,
        })
      );
      if (first) {
        createRequestForm(service, method);
      }
      first = false;
    }
  });

  $("#sidebar-methods-list").on("change", function () {
    saveRequest();
    let service = $("#sidebar-services-list option:selected").val();
    createRequestForm(
      currentServices[service],
      currentServices[service].methods[$(this).val()]
    );
  });

  $(".metadata-key").click(function () {
    addMetadataRow(this);
  });

  $(".request-metadata-button-delete").click(function () {
    removeMetadataRow(this);
  });

  const { app } = require("electron").remote;
  app.on("second-instance", (event, commandLine, workingDirectory) => {
    if (commandLine.length <= 4) {
      return;
    }
    for (const arg of commandLine.slice(4)) {
      const result = arg.split(/[=\s]/);
      if (result.length === 2) {
        switch (result[0]) {
          case "--workspace-id":
          case "-workspace-id":
            loadServer({ id: parseInt(result[1], 10) });
            break;
        }
      }
    }
  });
});

function initQueryPopover() {
  $("body").on("click", function (e) {
    $("[data-toggle=popover]").each(function () {
      if (
        !$(this).is(e.target) &&
        $(this).has(e.target).length === 0 &&
        $(".popover").has(e.target).length === 0
      ) {
        $(this).popover("hide");
      }
    });
  });

  document
  .getElementById("save-query")
  .addEventListener("shown.bs.popover", function () {
    if (!isNull(currentQuery)) {
      $(".popover-body .query-popover-query-name").val(currentQuery.text);
      $(".popover-body .query-popover-query-description").val(
        currentQuery.data.description
      );
    }
    $(".popover-body form").submit(function () {
      $("#save-query").popover("hide");
      saveQuery();
      return false;
    });
    $(".popover-body .query-popover-close").click(function () {
      $("#save-query").popover("hide");
    });
  });

  $("#save-query").popover({
    html: true,
    trigger: "click",
    placement: "bottom",
    container: "body",
    sanitize: false,
    content: function () {
      return $("#query-popover").html();
    },
  });
}

function initOffcanvas() {
  let offcanvasTree = document.getElementById("offcanvasTree");
  offcanvasTree.addEventListener("show.bs.offcanvas", function (event) {
    astilectron.sendMessage({ name: "workspace.get" }, function (message) {
      showTree(message.payload.data);
    });
  });
}

function initSidebars() {
  const resizerLeft = document.querySelector("#resizerLeft");
  const sidebarLeft = document.querySelector("#sidebarLeft");

  resizerLeft.addEventListener("mousedown", (event) => {
    $("#right-nav-tab").addClass("noselect");
    document.addEventListener("mousemove", resizeLeft, false);
    document.addEventListener(
      "mouseup",
      () => {
        $("#right-nav-tab").removeClass("noselect");
        document.removeEventListener("mousemove", resizeLeft, false);
      },
      false
    );
  });

  sidebarLeft.style.flexBasis = "50%";

  function resizeLeft(e) {
    sidebarLeft.style.flexBasis = `${e.x}px`;
  }
}

function addMetadataRow(elm) {
  let keys = $(".metadata-key");
  if (
    $(keys[keys.length - 1]).attr("data-metadata-key-id") !==
    $(elm).attr("data-metadata-key-id")
  ) {
    return;
  }
  let metadata = $($("#nav-request-metadata .metadata-row")[0]).clone();
  let key = $(metadata.find(".metadata-key")[0]);
  let remove = $(metadata.find(".request-metadata-button-delete")[0]);
  key.attr("data-metadata-key-id", Date.now()).val("");
  key.click(function () {
    addMetadataRow(this);
  });
  remove.click(function () {
    removeMetadataRow(this);
  });
  $(metadata.find(".metadata-value")[0]).val("");
  $("#nav-request-metadata").append(metadata);
}

function removeMetadataRow(elm) {
  if ($("#nav-request-metadata .metadata-row").length <= 1) {
    return;
  }
  $(elm).closest(".metadata-row").remove();
}

function showAbout(data) {
  let modal = $("#aboutModal");
  modal.find(".app-name").html(data.app_name);
  modal.find(".app-version").html(data.app_version);
  modal.find(".app-url").html(data.app_url).attr("href", data.app_url);
  modal.find(".go-version").html(data.go_version);
  modal.find(".astilectron-version").html(data.astilectron_version);
  modal.find(".electron-version").html(data.electron_version);
  modal.find(".built-at").html(data.built_at);
  modal.modal("show");
}

function isNull(v) {
  return v === undefined || v === null;
}
