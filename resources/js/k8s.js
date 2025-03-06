export {
    initK8S,
    setServerK8S,
    getServerK8S,
    syncK8SLocalPort,
}

import {isNull, loadFile} from "./index.js";

function initK8S() {
    $("#workspace-modal-k8s-enabled").change(function () {
        if ($(this).is(":checked")) {
            $("#nav-workspace-k8s .k8s").attr("disabled", false);
            syncK8SLocalPort();
        } else {
            $("#nav-workspace-k8s .k8s").attr("disabled", true);
        }
    });

    $("#workspace-modal-k8s-gcs-enabled").change(function () {
        if ($(this).is(":checked")) {
            $("#nav-workspace-k8s .gcs").show();
            $("#nav-workspace-k8s .gcs-control").attr("required");
        } else {
            $("#nav-workspace-k8s .gcs").hide();
            $("#nav-workspace-k8s .gcs-control").removeAttr("required");
        }
    });

    $("#k8s-config-select-file").click(function () {
        addK8SConfigFile();
    });

    $("#k8s-token-select-file").click(function () {
        loadFile($("#k8s-token"));
    });

    $("#k8s-pod-name").keyup(function () {
        let selector = $("#k8s-pod-name-selector");
        selector.attr("required");
        if ($(this).val() !== "") {
            selector.removeAttr("required");
        }
    });

    $("#k8s-pod-name-selector").keyup(function () {
        let name = $("#k8s-pod-name");
        name.attr("required");
        if ($(this).val() !== "") {
            name.removeAttr("required");
        }
    });

    $("#k8s-local-port").change(function () {
        if ($("#workspace-modal-k8s-enabled").is(":checked")) {
            syncK8SLocalPort();
        }
    })
}

function setServerK8S(k8s) {
    let enabled = $("#workspace-modal-k8s-enabled");

    if (isNull(k8s) || k8s === {} || isNull(k8s.enabled) || !k8s.enabled) {
        enabled.prop("checked", false).trigger('change');
        return;
    }

    enabled.prop("checked", true).trigger('change');

    if (!isNull(k8s.namespace)) {
        $("#k8s-namespace").val(k8s.namespace);
    }
    if (!isNull(k8s.local_port)) {
        $("#k8s-local-port").val(k8s.local_port);
    }
    if (!isNull(k8s.pod_port)) {
        $("#k8s-pod-port").val(k8s.pod_port);
    }
    if (!isNull(k8s.pod_name) && k8s.pod_name !== "") {
        $("#k8s-pod-name").val(k8s.pod_name);
        $("#k8s-pod-name-selector").removeAttr("required");
    }
    if (!isNull(k8s.pod_name_selector) && k8s.pod_name_selector !== "") {
        $("#k8s-pod-name-selector").val(k8s.pod_name_selector);
        $("#k8s-pod-name").removeAttr("required");
    }

    if (!isNull(k8s.client_config) && k8s.client_config !== {}) {
        if (!isNull(k8s.client_config.config_file)) {
            $("#k8s-config-file-path").val(k8s.client_config.config_file);
        }
        if (!isNull(k8s.client_config.bearer_token)) {
            $("#k8s-token").val(k8s.client_config.bearer_token);
        }

        if (!isNull(k8s.client_config.auth) && k8s.client_config.auth !== {}) {
            $("#workspace-modal-k8s-gcs-enabled").prop("checked", k8s.client_config.auth.enabled).trigger("change");
            if (!isNull(k8s.client_config.auth.project)) {
                $("#k8s-gcs-project").val(k8s.client_config.auth.project);
            }
            if (!isNull(k8s.client_config.auth.location)) {
                $("#k8s-gcs-location").val(k8s.client_config.auth.location);
            }
            if (!isNull(k8s.client_config.auth.cluster)) {
                $("#k8s-gcs-cluster").val(k8s.client_config.auth.cluster);
            }
        } else {
            $("workspace-modal-k8s-gcs-enabled").prop("checked", false).trigger('change');
        }
    }

    $("#k8s-local-port").change(function () {
        if ($("#workspace-modal-k8s-enabled").is(":checked")) {
            syncK8SLocalPort();
        }
    });
}

function getServerK8S() {
    if (!$("#workspace-modal-k8s-enabled").is(":checked")) {
        return {};
    }

    let k8s = {
        enabled: true,
        namespace: $("#k8s-namespace").val(),
        local_port: $("#k8s-local-port").val(),
        pod_port: $("#k8s-pod-port").val(),
        pod_name: $("#k8s-pod-name").val(),
        pod_name_selector: $("#k8s-pod-name-selector").val(),
        client_config: {
            config_file: $("#k8s-config-file-path").val(),
            bearer_token: $("#k8s-token").val(),
        }
    };

    if ($("#workspace-modal-k8s-gcs-enabled").is(":checked")) {
        k8s.client_config.auth = {
            enabled: true,
            project: $("#k8s-gcs-project").val(),
            location: $("#k8s-gcs-location").val(),
            cluster: $("#k8s-gcs-cluster").val(),
        };
    }

    return k8s;
}

function addK8SConfigFile() {
    const {dialog} = require("electron").remote;
    let files = dialog.showOpenDialogSync({
        properties: ["openFile", "showHiddenFiles"],
    });
    if (files === undefined) {
        return;
    }
    $("#k8s-config-file-path").val(files[0]);
}

function syncK8SLocalPort() {
    let addr = $("#workspace-modal-grpc-addr").val().split(':');
    let k8sPort = $("#k8s-local-port").val();
    let basicPort = '';
    if (addr.length === 2) {
        basicPort = addr[1];
    }
    if (basicPort !== '') {
        $("#k8s-local-port").val(basicPort);
    } else if (k8sPort !== '') {
        $("#workspace-modal-grpc-addr").val(addr[0] + ':' + k8sPort);
    }
}