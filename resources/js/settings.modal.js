import { currentSettings, isNull, setCurrentSettings } from "./index.js";

export { initSettingsModal, showSettingsModal };

function initSettingsModal() {
  $("#settings-modal-submit").click(function (event) {
    let form = $("#settings-modal-form")[0];
    if (!form.checkValidity()) {
      event.preventDefault();
      event.stopPropagation();
    } else {
      saveSettings();
    }
    form.classList.add("was-validated");
  });

  $("#settingsModal").on("hidden.bs.modal", function () {
    $(this).find("form").removeClass("was-validated");
    $("#settings-modal-form-single-instance-restart").css(
      "visibility",
      "hidden"
    );
  });

  $("#settings-modal-form-single-instance").on("change", function () {
    if (currentSettings.single_instance !== $(this).is(":checked")) {
      $("#settings-modal-form-single-instance-restart").css(
        "visibility",
        "visible"
      );
    } else {
      $("#settings-modal-form-single-instance-restart").css(
        "visibility",
        "hidden"
      );
    }
  });
}

function saveSettings() {
  let req = {
    name: "settings.update",
    payload: {
      connect_timeout: parseInt(
        $("#settings-modal-form-connect-timeout").val(),
        10
      ),
      request_timeout: parseInt(
        $("#settings-modal-form-request-timeout").val(),
        10
      ),
      k8s_request_timeout: parseInt(
          $("#settings-modal-form-k8s-request-timeout").val(),
          10
      ),
      max_loop_depth: parseInt(
        $("#settings-modal-form-max-loop-depth").val(),
        10
      ),
      non_blocking_connection: $(
        "#settings-modal-form-non-blocking-connection"
      ).is(":checked"),
      sort_methods_by_name: $("#settings-modal-form-sort-methods-by-name").is(
        ":checked"
      ),
      single_instance: $("#settings-modal-form-single-instance").is(":checked"),
    },
  };
  astilectron.sendMessage(req, function (message) {
    $("#settingsModal").modal("hide");
    if (message.payload.status !== "ok") {
      return;
    }
    setCurrentSettings(message.payload.data);
  });
}

function showSettingsModal() {
  if (isNull(currentSettings)) {
    return;
  }
  console.log(currentSettings);
  $("#settings-modal-form-connect-timeout").val(
    currentSettings.connect_timeout
  );
  $("#settings-modal-form-request-timeout").val(
    currentSettings.request_timeout
  );
  $("#settings-modal-form-k8s-request-timeout").val(
      currentSettings.k8s_request_timeout
  );
  $("#settings-modal-form-max-loop-depth").val(currentSettings.max_loop_depth);
  $("#settings-modal-form-non-blocking-connection").prop(
    "checked",
    currentSettings.non_blocking_connection
  );
  $("#settings-modal-form-sort-methods-by-name").prop(
    "checked",
    currentSettings.sort_methods_by_name
  );
  $("#settings-modal-form-single-instance").prop(
    "checked",
    currentSettings.single_instance
  );
  $("#settingsModal").modal("show");
}
