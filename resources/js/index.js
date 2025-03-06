export {
    currentSettings,
    modalTreeExpandedNodes,
    treeRootNodes,
    treeNodeNames,
    dataIdToNode,
    isNull,
    setCurrentSettings,
    loadFile,
    showModalError
};

import {saveQuery} from "./query.js";
import {showTree} from "./tree.js";
import {editServer, initWorkspaceModal} from "./workspace.modal.js";
import {initSettingsModal, showSettingsModal} from "./settings.modal.js";
import {
    createRequestForm,
    currentMethod,
    currentQuery,
    currentSelectedID,
    currentServer,
    currentService,
    currentServices,
    loadServer,
    saveRequest,
    setCurrentQuery,
    setRequestTitle,
} from "./server.js";
import {hideStreamControl, initStreamControl, query, response, showQueryError,} from "./request.js";
import {workspaceExport, workspaceImport} from "./workspace.export.js";

let currentSettings = undefined;
let treeRootNodes = new Set();
let modalTreeExpandedNodes = new Set();
let dataIdToNode = new Map();
let treeNodeNames = new Map();

$(document).ready(function () {
    document.addEventListener("astilectron-ready", function () {
        astilectron.onMessage(function (message) {
            console.log("server message: " + JSON.stringify(message, null, 1));
            if (isNull(message.name)) {
                return;
            }
            switch (message.name) {
                case "server.load":
                    loadServer({id: message.payload.id});
                    break;
                case "menu.export.file":
                    workspaceExport();
                    break;
                case "menu.import.file":
                    workspaceImport();
                    break;
                case "menu.settings":
                    showSettingsModal();
                    break;
                case "menu.about":
                    showAbout(message.payload);
                    break;
                case "query.response":
                    response(message.payload);
                    break;
                case "message.info":
                    addInfoMessage(message.payload);
                    break;
                case "message.error":
                    showQueryError(message.payload)
                    break;
                case "check.updates":
                    showUpdates(message.payload);
                    break;
            }
        });

        astilectron.sendMessage({name: "application.state"}, function (message) {
            if (message.payload.status !== "ok") {
                return;
            }
            currentSettings = message.payload.data.settings;
            if (message.payload.data.state.servers === 0) {
                $("#workspaceModal").modal("show");
                return;
            }
            if (
                isNull(message.payload.data.state.startup_workspace_id) ||
                message.payload.data.state.startup_workspace_id === 0
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
                case "modal.settings.html":
                    initSettingsModal();
                    break;
            }
        });
    });

    initSidebars();
    initOffcanvas();
    initQueryPopover();
    initStreamControl();

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

    let navRequest = $("#nav-request form");
    navRequest.submit(function () {
        query();
        return false;
    });

    navRequest.keypress(function (e) {
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
        setCurrentQuery(undefined);
        setRequestTitle("");
        hideStreamControl();

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
        hideStreamControl();
        let service = $("#sidebar-services-list option:selected").val();
        if (!isNull(currentQuery)) {
            currentServer.breadcrumb.pop();
        }
        setCurrentQuery(undefined);
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

    $("#offcanvas-tree-search").on("input", function () {
        let v = $(this).val();
        if (v.length > 0) {
            $("#tree").treeview("search", v, {revealResults: false});
        } else {
            $("#tree").treeview("clearSearch");
        }
    });

    const {app} = require("electron").remote;
    app.on("second-instance", (event, commandLine) => {
        if (commandLine.length <= 4) {
            return;
        }
        for (const arg of commandLine.slice(4)) {
            const result = arg.split(/[=\s]/);
            if (result.length === 2) {
                switch (result[0]) {
                    case "--workspace-id":
                    case "-workspace-id":
                        loadServer({id: parseInt(result[1], 10)});
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
    offcanvasTree.addEventListener("show.bs.offcanvas", function () {
        astilectron.sendMessage(
            {name: "workspace.get", payload: {selected_id: currentSelectedID}},
            function (message) {
                $("#offcanvas-tree-search").val("");
                showTree(message.payload.data);
            }
        );
    });
}

function initSidebars() {
    const resizerLeft = document.querySelector("#resizerLeft");
    const sidebarLeft = document.querySelector("#sidebarLeft");

    resizerLeft.addEventListener("mousedown", () => {
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
        let row = $(elm).closest(".metadata-row");
        row.find(".metadata-key").val("");
        row.find(".metadata-value").val("");
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

function showUpdates(data) {
    let modal = $("#updatesModal");
    if (!isNull(data)) {
        modal.find(".new-version").show();
        modal.find(".release-url").attr("href", data.url);
        modal.find("modal-title").html("Update Available")
        modal.find(".no-updates").hide();
    } else {
        modal.find(".new-version").hide();
        modal.find("modal-title").html("No Updates Available")
        modal.find(".no-updates").show();
    }
    modal.modal("show");
}

function isNull(v) {
    return v === undefined || v === null;
}

function setCurrentSettings(settinigs) {
    currentSettings = settinigs;
}

function loadFile(elm) {
    const {dialog} = require("electron").remote;
    let files = dialog.showOpenDialogSync({
        properties: ["openFile", "showHiddenFiles"],
    });
    if (files === undefined) {
        return;
    }

    const fs = require("fs");
    fs.readFile(files[0], (err, contents) => {
        if (err) {
            return;
        }
        elm.val(contents);
    })
}

function addInfoMessage(data) {
    let info = $("#info-message");
    if (data.message === "") {
        info.html("").hide();
        return
    }
    info.append("<div>" + data.message + "</div>");
    info.show();
}

function showModalError(message) {
    $('#errorModal .error-description').html(message.payload.error.message);
    $('#errorModal').modal("show");
}
