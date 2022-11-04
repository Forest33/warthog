function initWorkspaceModal() {
    let workspaceModal = document.getElementById('workspaceModal')
    workspaceModal.addEventListener('show.bs.modal', function (event) {
        let req = {
            name: "workspace.get",
            payload: {
                type: [WorkspaceTypeFolder]
            }
        }

        let folderID = $('#workspace-modal-folder-id').val()
        let serverID = $('#workspace-modal-server-id').val()
        if (serverID) {
            $('#workspaceModalLabel').html('Edit workspace')
            $('#workspace-modal-badge-server-id').html('ID: ' + serverID).css('visibility', 'visible')
        }

        astilectron.sendMessage(req, function (message) {
            let tree = $('#workspace-modal-tree')
            let selectedNode = undefined
            tree.treeview({
                data: message.payload.data,
                expandIcon: "fa fa-angle-down fa-fw",
                collapseIcon: "fa fa-angle-right fa-fw",
                onNodeRender: function (node) {
                    node.state = {};
                    if (modalTreeExpandedNodes.has(node.data.id)) {
                        node.state.expanded = true
                    }
                    if (folderID !== undefined && parseInt(folderID) === node.data.id) {
                        selectedNode = node
                    }
                    return onTreeNodeRender(node)
                },
                onNodeSelected: onModalWorkspaceTreeNodeSelect,
                onNodeUnselected: onModalWorkspaceTreeNodeSelect,
                onNodeExpanded: function (event, node) {
                    modalTreeExpandedNodes.add(node.data.id)
                },
                onNodeCollapsed: function (event, node) {
                    modalTreeExpandedNodes.delete(node.data.id)
                },
            });
            if (selectedNode !== undefined) {
                tree.treeview('selectNode', selectedNode.nodeId);
            }
        });

        $('#workspaceModal').on('hidden.bs.modal', function () {
            $(this).find('form').trigger('reset')
            $('#nav-workspace-modal-basic-tab').tab('show')
            document.getElementById("workspace-modal-proto-files").options.length = 0
            document.getElementById("workspace-modal-import-path").options.length = 0
            $('#workspace-modal-submit').text("Save").prop("disabled", true)
            $('#workspaceModalLabel').html('New workspace')
            $('#workspace-modal-server-id').val('')
            $('#workspace-modal-folder-id').val('')
            $('#workspace-modal-badge-server-id').css('visibility', 'hidden')
            $('#workspace-modal-form').removeClass('was-validated')
        })
    });

    $("#workspace-modal-new-folder").click(function () {
        createFolder();
    });

    $("#workspace-modal-submit").click(function (event) {
        let form = $('#workspace-modal-form')[0]
        if (!form.checkValidity()) {
            event.preventDefault()
            event.stopPropagation()
        } else {
            createWorkspace()
        }
        form.classList.add('was-validated')
    });

    $("#workspace-modal-add-proto-files").click(function () {
        addProtoFiles();
    })

    $("#workspace-modal-add-import-path").click(function () {
        addImportPath();
    })

    $("#workspace-modal-delete-proto-files").click(function () {
        $("#workspace-modal-proto-files option:selected").remove()
    })

    $("#workspace-modal-delete-import-path").click(function () {
        $("#workspace-modal-import-path option:selected").remove()
    })
}

function createFolder() {
    let title = $('#workspace-modal-folder-title').val()
    if (title === "") {
        return
    }

    let req = {
        name: "folder.create",
        payload: {
            title: title,
            type_filter: [WorkspaceTypeFolder]
        }
    }

    let selected = $('#workspace-modal-tree').treeview('getSelected')
    if (selected.length === 1) {
        req.payload.parent_id = selected[0].data.id
        modalTreeExpandedNodes.add(req.payload.parent_id)
    }

    astilectron.sendMessage(req, function (message) {
        let folderId = message.payload.data.folder.id;
        let tree = $('#workspace-modal-tree')
        let folderNode = {}

        tree.treeview(true).remove()
        tree.treeview({
            data: message.payload.data.tree,
            expandIcon: "bi bi-caret-right",
            collapseIcon: "bi bi-caret-down",
            onNodeSelected: onModalWorkspaceTreeNodeSelect,
            onNodeUnselected: onModalWorkspaceTreeNodeSelect,
            onNodeExpanded: function (event, node) {
                modalTreeExpandedNodes.add(node.data.id)
            },
            onNodeCollapsed: function (event, node) {
                modalTreeExpandedNodes.delete(node.data.id)
            },
            onNodeRender: function (node) {
                node = onTreeNodeRender(node);
                node.state = {};
                if (modalTreeExpandedNodes.has(node.data.id)) {
                    node.state.expanded = true
                }
                if (node.data.id === folderId) {
                    node.state.selected = true
                    folderNode = node
                }
                return node;
            },
        });
        document.querySelectorAll('[data-nodeid="' + folderNode.nodeId + '"]')[0].scrollIntoView()
        $('#workspace-modal-submit').text("Save to " + folderNode.data.text).prop("disabled", false)
    });
}

function createWorkspace(parent) {
    let title = $('#workspace-modal-grpc-name').val()
    let addr = $('#workspace-modal-grpc-addr').val()

    let selected = $('#workspace-modal-tree').treeview('getSelected')
    if (selected.length !== 1) {
        return;
    }
    let folderID = selected[0].data.id

    let useReflection = $('#workspace-modal-use-reflection').is(':checked')
    let noTLS = $('#workspace-modal-use-plain-text').is(':checked')
    let insecure = $('#workspace-modal-skip-verification').is(':checked')
    let rootCertificate = $('#workspace-modal-root-certificate').val()
    let clientCertificate = $('#workspace-modal-client-certificate').val()
    let clientKey = $('#workspace-modal-client-key').val()

    let protoFiles = [], importPath = []
    $("#workspace-modal-proto-files option").each(function () {
        protoFiles.push($(this).val())
    });
    $("#workspace-modal-import-path option").each(function () {
        importPath.push($(this).val())
    });

    // TODO валидация

    let req = {
        name: "server.create",
        payload: {
            folder_id: folderID,
            title: title,
            addr: addr,
            proto_files: protoFiles,
            import_path: importPath,
            use_reflection: useReflection,
            no_tls: noTLS,
            insecure: insecure,
            root_certificate: rootCertificate,
            client_certificate: clientCertificate,
            client_key: clientKey,
        }
    }

    let serverID = $('#workspace-modal-server-id').val()
    if (serverID !== undefined && serverID !== '') {
        req.name = "server.update"
        req.payload.id = parseInt(serverID)
    }

    astilectron.sendMessage(req, function (message) {
        $('#workspaceModal').modal('hide');
        if (message.payload.status !== "ok") {
            return
        }
        loadServer({id: message.payload.data.server.id})
    });
}

function addProtoFiles() {
    const {dialog} = require('electron').remote;
    let files = dialog.showOpenDialogSync({
        properties: ['openFile', 'multiSelections'],
        filters: [{name: '*.proto', extensions: ['proto']}]
    });
    if (files === undefined) {
        return
    }

    let sel = $('#workspace-modal-proto-files')
    for (let i = 0; i < files.length; i++) {
        sel.append($('<option>', {
            value: files[i],
            text: files[i]
        }));
    }
}

function addImportPath() {
    const {dialog} = require('electron').remote;
    let path = dialog.showOpenDialogSync({
        properties: ['openDirectory', 'multiSelections']
    });
    if (path === undefined) {
        return
    }

    let sel = $('#workspace-modal-import-path')
    for (let i = 0; i < path.length; i++) {
        sel.append($('<option>', {
            value: path[i],
            text: path[i]
        }));
    }
}

function onModalWorkspaceTreeNodeSelect(event, node) {
    let submit = $('#workspace-modal-submit')
    switch (event.type) {
        case "nodeSelected":
            submit.text("Save to " + node.data.text).prop("disabled", false)
            break
        case "nodeUnselected":
            submit.text("Save").prop("disabled", true)
            break
    }
}

function editServer(srv) {
    $('#workspace-modal-server-id').val(srv.id)
    $('#workspace-modal-folder-id').val(srv.parent_id)
    $('#workspace-modal-grpc-name').val(srv.text)
    $('#workspace-modal-grpc-addr').val(srv.data.addr)
    $('#workspace-modal-use-reflection').prop('checked', srv.data.use_reflection);
    $('#workspace-modal-use-plain-text').prop('checked', srv.data.no_tls);
    $('#workspace-modal-skip-verification').prop('checked', srv.data.insecure);
    $('#workspace-modal-root-certificate').val(srv.data.root_certificate)
    $('#workspace-modal-client-certificate').val(srv.data.client_certificate)
    $('#workspace-modal-client-key').val(srv.data.client_key)

    if (srv.data.proto_files !== undefined) {
        let sel = $('#workspace-modal-proto-files')
        for (const f of srv.data.proto_files) {
            sel.append($('<option>', {
                value: f,
                text: f
            }));
        }
    }
    if (srv.data.import_path !== undefined) {
        let sel = $('#workspace-modal-import-path')
        for (const f of srv.data.import_path) {
            sel.append($('<option>', {
                value: f,
                text: f
            }));
        }
    }

    $('#workspaceModal').modal('show');
}


