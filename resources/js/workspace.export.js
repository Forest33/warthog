export {
    workspaceExport,
    workspaceImport
};

import {showModalError} from "./index.js";

function workspaceExport() {
    const {dialog} = require("electron").remote;
    let path = dialog.showSaveDialogSync({
        defaultPath: "warthog-export.json",
        properties: ["showHiddenFiles"],
    });
    if (path === undefined) {
        return;
    }

    astilectron.sendMessage(
        {name: "workspace.export.file", payload: {path: path}},
        function (message) {
            if (message.payload.status === "ok") {
                return;
            }
            showModalError(message);
        }
    );
}

function workspaceImport() {
    const {dialog} = require("electron").remote;
    let path = dialog.showOpenDialogSync({
        defaultPath: "warthog-export.json",
        properties: ["openFile", "showHiddenFiles"],
    });
    if (path === undefined) {
        return;
    }

    astilectron.sendMessage(
        {name: "workspace.import.file", payload: {path: path[0]}},
        function (message) {
            if (message.payload.status === "ok") {
                return;
            }
            showModalError(message);
        }
    );
}