export {
    initAuth,
    setServerAuth,
    getServerAuth,
    validateAuthJWTPayload
}

import {isNull, loadFile} from "./index.js";
import {currentServer} from "./server.js";

let jwtStandardClaims = {
    exp: 0,
    nbf: 0,
    iss: "",
    sub: "",
    aud: "",
    iat: "",
    jti: "",
}

function initAuth() {
    $("#authentication-type").on("change", function () {
        let authType = $(this).val();
        $(".authentication").hide();
        if (authType === "none") {
            return
        }
        $(".authentication-" + authType).show();
        $(".auth .symmetric").hide();
        $(".auth .asymmetric").hide();
        if (authType === "jwt") {
            if (!isNull(currentServer) && !isNull(currentServer.data.auth) && currentServer.data.auth.type === authType) {
                $("#authentication-jwt-algorithm").val(currentServer.data.auth.algorithm).trigger('change');
            } else {
                $("#authentication-jwt-algorithm").prop("selectedIndex", 0).trigger('change');
            }
        }
    });

    $("#authentication-jwt-algorithm").on("change", function () {
        let symmetric = $(".auth .symmetric");
        let asymmetric = $(".auth .asymmetric");
        symmetric.hide();
        asymmetric.hide();
        if ($(this).val() === "HS256" || $(this).val() === "HS384" || $(this).val() === "HS512") {
            symmetric.show();
        } else {
            asymmetric.show();
        }
    });

    $("#authentication-jwt-private-key-select-file").click(function () {
        loadFile($("#authentication-jwt-private-key"));
    });
    $("#authentication-google-token-select-file").click(function () {
        loadFile($("#authentication-google-token"));
    });

    $("#authentication-jwt-standard-claims").click(function () {
        jwtStandardClaims.nbf = Math.floor(Date.now() / 1000);
        jwtStandardClaims.exp = jwtStandardClaims.nbf + 86400;
        $("#authentication-jwt-payload").val(JSON.stringify(jwtStandardClaims, null, 2));
    });

    $(document).on('input propertychange', "textarea[id='authentication-jwt-payload']", function () {
        if ($("#authentication-jwt-payload").hasClass("error")) {
            validateAuthJWTPayload();
        }
    });
}

function setServerAuth(auth) {
    let authType = $('#authentication-type');

    if (isNull(auth) || auth === {} || isNull(auth.type) || auth.type === "none") {
        authType.val("none").trigger('change');
        return;
    }

    authType.val(auth.type).trigger('change');

    if (!isNull(auth.login)) {
        $("#authentication-basic-login").val(auth.login);
    }
    if (!isNull(auth.password)) {
        $("#authentication-basic-password").val(auth.password);
    }
    if (!isNull(auth.token)) {
        $("#authentication-bearer-token").val(auth.token);
    }
    if (!isNull(auth.algorithm)) {
        $("#authentication-jwt-algorithm").val(auth.algorithm).trigger('change');
    }
    if (!isNull(auth.secret)) {
        $("#authentication-jwt-secret").val(auth.secret);
    }
    if (!isNull(auth.header_prefix)) {
        $("#authentication-jwt-header-prefix").val(auth.header_prefix);
    }
    if (!isNull(auth.private_key)) {
        $("#authentication-jwt-private-key").val(auth.private_key);
    }
    if (!isNull(auth.secret_base64)) {
        $("#authentication-jwt-secret-base64").prop("checked", auth.secret_base64);
    }
    if (!isNull(auth.payload)) {
        $("#authentication-jwt-payload").val(JSON.stringify(auth.payload, null, 2));
    }
    if (!isNull(auth.google_scopes)) {
        $("#authentication-google-scopes").val(auth.google_scopes);
    }
    if (!isNull(auth.google_token)) {
        $("#authentication-google-token").val(auth.google_token);
    }
}

function getServerAuth() {
    let authType = $("#authentication-type option:selected").val();
    if (authType === "none") {
        return {};
    }

    let resp = {};
    switch (authType) {
        case "basic":
            resp = {
                "login": $("#authentication-basic-login").val(),
                "password": $("#authentication-basic-password").val(),
            };
            break;
        case "bearer":
            resp = {
                "token": $("#authentication-bearer-token").val(),
            };
            break;
        case "jwt":
            resp = {
                "algorithm": $("#authentication-jwt-algorithm option:selected").val(),
                "secret": $("#authentication-jwt-secret").val(),
                "secret_base64": $("#authentication-jwt-secret-base64").is(":checked"),
                "header_prefix": $("#authentication-jwt-header-prefix").val(),
                "private_key": $("#authentication-jwt-private-key").val(),
                "payload": $("#authentication-jwt-payload").val(),
            };
            break;
        case "google":
            resp = {
                "google_scopes": $("#authentication-google-scopes").val(),
                "google_token": $("#authentication-google-token").val(),
            };
            break;
    }

    resp.type = authType;

    return resp;
}


function validateAuthJWTPayload() {
    if ($("#authentication-type option:selected").val() !== "jwt") {
        return true;
    }

    let input = $("#authentication-jwt-payload");
    input.removeClass("error");
    let payload = input.val();
    if (isNull(payload) || payload === "") {
        return true;
    }

    try {
        JSON.parse(payload);
    } catch (e) {
        input.addClass("error");
        return false;
    }

    return true;
}



