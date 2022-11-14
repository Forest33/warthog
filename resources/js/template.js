export { template };

let template = {
  "request-text-input": `<div class="mb-2 form-group request-input">
                                <label class="label-name"></label>
                                <label class="label-type"></label>
                                <input type="text" class="form-control field-value">
                            </div>`,
  "request-bool-input": `<div class="mb-2 form-group request-input">
                                <div class="form-check form-switch">
                                    <input class="form-check-input field-value" type="checkbox">
                                    <label class="form-check-label label-name"></label>
                                    <label class="form-check-label label-type"></label>
                                </div>
                            </div>`,
  "request-select-input": `<div class="mb-2 form-group request-input">
                                <label class="label-name"></label>
                                <label class="label-type"></label>
                                <select class="form-select field-value"></select>
                            </div>`,
  "request-message-input": `<div class="mb-2 form-group request-input request-container-input">
                                <button type="button" class="btn btn-primary btn-sm request-message-button-add" data-input-count="0"><i class="bi bi-plus-lg"></i></button>
                                <label class="form-check-label label-name">name</label>
                                <label class="form-check-label label-type">type</label>
                                <div class="mb-2 form-group request-message-container">                                    
                                </div>
                              </div>`,
  "request-message-input-delete": `<div class="mb-2 form-group request-input request-message-input-delete">
                                        <button type="button" class="btn btn-secondary btn-sm request-message-button-delete"><i class="bi bi-x-lg"></i></button>
                                     </div>`,
  "request-oneof-select": `<div class="mb-2 form-group request-input">
                                <label class="label-name"></label>
                                <label class="label-type"></label>
                                <select class="form-select field-value"></select>
                                <div class="mb-2 form-group request-message-container">                                    
                                </div>
                             </div>`,
  "query-error": `<div class="alert alert-danger" role="alert">
                        <table>
                        <tr><td class="error-label">code: </td><td class="code"></td></tr>
                        <tr><td class="error-label">message: </td><td class="message"></td></tr>
                        </table>
                    </div>`,
  "protobuf-error": `<div class="alert proto-error" role="alert">
                        <table>
                        <tr><td class="error-label">file: </td><td class="file"></td></tr>
                        <tr><td class="error-label">line: </td><td class="line"></td></tr>
                        <tr><td class="error-label">column: </td><td class="column"></td></tr>
                        <tr><td class="error-label">message: </td><td class="message"></td></tr>
                        </table>    
                    </div>`,
};
