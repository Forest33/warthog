export { showTree };
export {
  WorkspaceTypeFolder,
  WorkspaceTypeServer,
  WorkspaceTypeQuery,
  onTreeNodeRender,
};
import {
  currentQuery,
  currentServer,
  loadServer,
  saveRequest,
} from "./server.js";
import { dataIdToNode, isNull, treeRootNodes } from "./index.js";
import { editServer } from "./workspace.modal.js";

const WorkspaceTypeFolder = "f";
const WorkspaceTypeServer = "s";
const WorkspaceTypeQuery = "r";

function initTreeDrag() {
  let sourceNode = undefined;
  let source = undefined;

  let items = document.querySelectorAll("#tree .list-group-item");
  items.forEach(function (item) {
    item.addEventListener("dragstart", handleDragStart);
    item.addEventListener("dragover", handleDragOver);
    item.addEventListener("dragenter", handleDragEnter);
    item.addEventListener("dragleave", handleDragLeave);
    item.addEventListener("dragend", handleDragEnd);
    item.addEventListener("drop", handleDrop);
  });

  $("#tree .list-group-item")
    .hover(function () {
      $(this).find(".tree-node-dots").css("visibility", "visible");
    })
    .mouseleave(function () {
      $(this).find(".tree-node-dots").css("visibility", "hidden");
    });

  let dropdownMenuId = undefined;
  $("#tree .tree-node-dots").click(function (e) {
    $("#tree .list-group-item").removeClass("show-menu hidden");
    $("#tree-folder").remove();
    $(this).closest(".list-group-item").addClass("show-menu");
    if (!isNull(dropdownMenuId) && e.target.id !== dropdownMenuId) {
      $('[aria-labelledby="' + dropdownMenuId + '"]').removeClass("show");
    }
    dropdownMenuId = e.target.id;
    e.stopPropagation();
  });

  $("#tree .dropdown-item").click(function (e) {
    e.stopPropagation();
    $(".dropdown-menu").removeClass("show");

    let node = dataIdToNode.get(parseInt($(e.target).attr("data-id"), 10));
    switch ($(e.target).attr("data-action")) {
      case "create-folder":
        treeMenuCreateFolder(node);
        break;
      case "create-server":
        treeMenuCreateServer(node);
        break;
      case "rename-folder":
        treeMenuRenameFolder(node);
        break;
      case "edit-server":
        treeMenuEditServer(node);
        break;
      case "delete":
        treeMenuDelete(node);
        break;
    }
  });

  function handleDragStart(e) {
    source = e.target;
    sourceNode = getNode(e);
    this.style.opacity = "0.4";
  }

  function handleDragEnd() {
    this.style.opacity = "1";
    items.forEach(function (item) {
      item.classList.remove("over");
    });
  }

  function handleDragOver(e) {
    if (e.preventDefault) {
      e.preventDefault();
    }
    return false;
  }

  function handleDragEnter(e) {
    let targetNode = getNode(e);
    if (!canDrop(targetNode)) {
      return;
    }
    this.classList.add("over");
  }

  function handleDragLeave(e) {
    let targetNode = getNode(e);
    this.classList.remove("over");
  }

  function handleDrop(e) {
    let targetNode = getNode(e);
    e.stopPropagation();
    if (!canDrop(targetNode)) {
      return;
    }

    e.target.parentNode.insertBefore(source, e.target.nextSibling);

    let list = [];
    if (targetNode.data.type === sourceNode.data.type) {
      if (
        sourceNode.data.type === WorkspaceTypeServer ||
        sourceNode.data.type === WorkspaceTypeQuery
      ) {
        $(source).find(".indent").remove();
        for (let i = 0; i < targetNode.level - 1; i++) {
          $(source).prepend('<span class="indent">');
        }
        list = insertAfter(getParent(targetNode).nodes, sourceNode, targetNode);
      } else {
        $(source).find(".indent").remove();
        for (let i = 0; i < targetNode.level - 1; i++) {
          $(source).prepend('<span class="indent">');
        }
        if (targetNode.state.expanded) {
          sourceNode.data.parent_id = targetNode.data.id;
          list = [sourceNode].concat(targetNode.nodes);
        } else {
          if (isNull(targetNode.data.parent_id)) {
            list = insertAfter(treeRootNodes.values(), sourceNode, targetNode);
          } else {
            list = insertAfter(
              getParent(targetNode).nodes,
              sourceNode,
              targetNode
            );
          }
        }
      }
    } else if (
      sourceNode.data.type === WorkspaceTypeQuery &&
      targetNode.data.type === WorkspaceTypeServer
    ) {
      list = insertAfter(getParent(sourceNode).nodes, sourceNode, null);
    } else if (
      sourceNode.data.type === WorkspaceTypeServer &&
      targetNode.data.type === WorkspaceTypeFolder
    ) {
      if (targetNode.state.expanded) {
        list = insertAfter(getParent(sourceNode).nodes, sourceNode, null);
      } else {
        let parent = getParent(targetNode);
        if (!isNull(parent)) {
          list = insertAfter(parent.nodes, sourceNode, null);
        } else {
          list = insertAfter(treeRootNodes.values(), sourceNode, null);
        }
      }
    }

    if (list.length === 0) {
      return false;
    }

    let nodes = [];
    for (const item of list) {
      nodes.push({ id: item.data.id, parent_id: item.data.parent_id });
    }

    let req = {
      name: "workspace.sorting",
      payload: {
        nodes: nodes,
      },
    };

    astilectron.sendMessage(req, function (message) {
      if (message.payload.status !== "ok") {
        return;
      }
      showTree(message.payload.data);
    });

    return false;
  }

  function canDrop(targetNode) {
    if (sourceNode.data.id === targetNode.data.id) {
      return false;
    }
    if (
      sourceNode.data.type === WorkspaceTypeQuery &&
      sourceNode.data.parent_id === targetNode.data.id
    ) {
      return true;
    }
    if (
      sourceNode.data.type === WorkspaceTypeQuery &&
      targetNode.data.parent_id !== sourceNode.data.parent_id
    ) {
      return false;
    }
    if (
      sourceNode.data.type === WorkspaceTypeFolder &&
      targetNode.data.type !== WorkspaceTypeFolder
    ) {
      return false;
    }
    if (
      sourceNode.data.type === WorkspaceTypeServer &&
      targetNode.data.type === WorkspaceTypeQuery
    ) {
      return false;
    }
    return true;
  }

  function insertAfter(nodes, sourceNode, targetNode) {
    let newNodes = [];
    if (targetNode === null) {
      newNodes.push(sourceNode);
    }
    for (const n of nodes) {
      if (n.nodeId === sourceNode.nodeId) {
        continue;
      }
      newNodes.push(n);
      if (targetNode !== null && n.nodeId === targetNode.nodeId) {
        sourceNode.data.parent_id = targetNode.data.parent_id;
        newNodes.push(sourceNode);
      }
    }
    return newNodes;
  }

  function getNode(e) {
    return $("#tree").treeview("getNode", $(e.target).attr("data-nodeid"));
  }

  function getParent(node) {
    return $("#tree").treeview("getParent", node);
  }
}

function getTreeDropdown(node) {
  let menu = $(
    '<ul class="dropdown-menu" aria-labelledby="tree-dropdown-' +
      node.data.id +
      '"></ul>'
  );

  let ext = "";
  if (node.data.has_child) {
    ext = " disabled";
  }

  switch (node.data.type) {
    case WorkspaceTypeFolder:
      menu.append(
        $(
          '<li><a class="dropdown-item" data-id="' +
            node.data.id +
            '" data-action="create-folder"><i class="bi bi-file-plus"></i> Add folder</a></li>'
        )
      );
      menu.append(
        $(
          '<li><a class="dropdown-item" data-id="' +
            node.data.id +
            '" data-action="create-server"><i class="bi bi-server"></i> Add server</a></li>'
        )
      );
      menu.append(
        $(
          '<li><a class="dropdown-item" data-id="' +
            node.data.id +
            '" data-action="rename-folder"><i class="bi bi-pen"></i> Rename</a></li>'
        )
      );
      menu.append(
        $(
          '<li><a class="dropdown-item' +
            ext +
            '" data-id="' +
            node.data.id +
            '" data-action="delete"><i class="bi bi-trash"></i> Delete</a></li>'
        )
      );
      break;
    case WorkspaceTypeServer:
      menu.append(
        $(
          '<li><a class="dropdown-item" data-id="' +
            node.data.id +
            '" data-action="edit-server"><i class="bi bi-pen"></i> Edit</a></li>'
        )
      );
      menu.append(
        $(
          '<li><a class="dropdown-item' +
            ext +
            '" data-id="' +
            node.data.id +
            '" data-action="delete"><i class="bi bi-trash"></i> Delete</a></li>'
        )
      );
      break;
    case WorkspaceTypeQuery:
      menu.append(
        $(
          '<li><a class="dropdown-item" data-id="' +
            node.data.id +
            '" data-action="delete"><i class="bi bi-trash"></i> Delete</a></li>'
        )
      );
      break;
  }

  return (
    '<span class="bi bi-three-dots tree-node-dots" data-bs-toggle="dropdown" aria-expanded="false" id="tree-dropdown-' +
    node.data.id +
    '"></span>' +
    menu.prop("outerHTML")
  );
}

function treeMenuCreateFolder(node) {
  let indent = '<span class="indent"></span>';
  for (let i = 0; i < node.level; i++) {
    indent += '<span class="indent"></span>';
  }

  $("#tree-folder").remove();
  let input = `<li class="list-group-item node-tree" id="tree-folder">
                 <form id="tree-folder-form">
                 <div class="input-group mb-3" style="margin-bottom: 0!important;">
                    ${indent}
                    <input type="text" minlength="1" class="form-control external" id="tree-folder-name" placeholder="New folder name" aria-label="New folder name" aria-describedby="tree-folder-submit" required>
                    <button type="submit" class="btn btn-primary external" id="tree-folder-submit"><i class="bi bi-plus-circle external"></i></button>
                 </div>
                 </form>
                 </li>`;
  $(input).insertAfter(
    $('#tree .list-group-item[data-nodeid="' + node.nodeId + '"]')
  );

  treeFolderHandler(node, true);
}

function treeMenuRenameFolder(node) {
  let indent = "";
  for (let i = 0; i < node.level - 1; i++) {
    indent += '<span class="indent"></span>';
  }

  $("#tree-folder").remove();
  let input = `<li class="list-group-item node-tree" id="tree-folder">
                 <form id="tree-folder-form">
                 <div class="input-group mb-3" style="margin-bottom: 0!important;">
                    ${indent}
                    <input type="text" minlength="1" class="form-control external" id="tree-folder-name" value="${node.data.text}" placeholder="Folder name" aria-label="New folder name" aria-describedby="tree-folder-submit" required>
                    <button type="submit" class="btn btn-primary external" id="tree-folder-submit"><i class="bi bi-check-circle external"></i></button>
                 </div>
                 </form>
                 </li>`;
  let treeItem = $('#tree .list-group-item[data-nodeid="' + node.nodeId + '"]');
  treeItem.addClass("hidden");
  $(input).insertAfter(treeItem);

  treeFolderHandler(node, false);
}

function treeFolderHandler(node, isCreate) {
  $("#tree-folder-form")
    .submit(function (event) {
      if (!this.checkValidity()) {
        event.preventDefault();
        event.stopPropagation();
      } else {
        let req;
        if (isCreate) {
          req = {
            name: "folder.create",
            payload: {
              parent_id: node.data.id,
              title: $("#tree-folder-name").val(),
            },
          };
        } else {
          req = {
            name: "folder.update",
            payload: {
              id: node.data.id,
              title: $("#tree-folder-name").val(),
            },
          };
        }
        astilectron.sendMessage(req, function (message) {
          if (message.payload.status !== "ok") {
            return;
          }
          showTree(message.payload.data.tree);
        });
        $("#tree-folder").remove();
      }

      $(this).addClass("was-validated");

      return false;
    })
    .keypress(function (e) {
      if (e.which === 13) {
        $("#tree-folder-form").submit();
        return false;
      }
    });
  $("#tree-folder-submit").click(function () {
    $("#tree-folder-form").submit();
  });
}

function treeMenuCreateServer(node) {
  $("#workspace-modal-folder-id").val(node.data.id);
  $("#workspaceModal").modal("show");
}

function treeMenuEditServer(node) {
  editServer(node.data);
}

function treeMenuDelete(node) {
  let req = {
    name: "workspace.delete",
    payload: {
      id: node.data.id,
    },
  };
  astilectron.sendMessage(req, function (message) {
    if (message.payload.status !== "ok") {
      return;
    }
    showTree(message.payload.data);
  });
}

function treeNodeExpand(node, expand) {
  let req = {
    name: "workspace.expand",
    payload: {
      id: node.data.id,
      expand: expand,
    },
  };
  astilectron.sendMessage(req, function (message) {});
}

function showTree(data) {
  treeRootNodes.clear();
  dataIdToNode.clear();
  $("#tree").treeview({
    data: data,
    expandIcon: "bi bi-caret-right",
    collapseIcon: "bi bi-caret-down",
    onNodeExpanded: function (event, node) {
      treeNodeExpand(node, true);
    },
    onNodeCollapsed: function (event, node) {
      treeNodeExpand(node, false);
    },
    onNodeRender: function (node) {
      if (isNull(node.data.parent_id)) {
        treeRootNodes.add(node);
      }
      return onTreeNodeRender(node);
    },
    onNodeRendered: function (id, node) {
      dataIdToNode.set(node.data.id, node);
      return node;
    },
    onNodeSelected: function (event, node) {
      if (node.data.type === WorkspaceTypeServer) {
        saveRequest();
        loadServer(node.data);
      } else if (node.data.type === WorkspaceTypeQuery) {
        loadServer(node.data);
      }
    },
    onTreeRenderComplete: function () {
      initTreeDrag();
    },
  });
}

function onTreeNodeRender(node) {
  if (node.state === undefined) {
    node.state = { draggable: true };
  }
  node.state.expanded = node.data.expanded;
  node.text = node.text + getTreeDropdown(node);
  switch (node.data.type) {
    case WorkspaceTypeFolder:
      node.icon = "bi bi-folder";
      break;
    case WorkspaceTypeServer:
      node.icon = "bi bi-server";
      if (
        isNull(currentQuery) &&
        !isNull(currentServer) &&
        node.data.id === currentServer.id
      ) {
        node.state.selected = true;
      }
      break;
    case WorkspaceTypeQuery:
      node.icon = "bi bi-card-list";
      if (!isNull(currentQuery) && node.data.id === currentQuery.id) {
        node.state.selected = true;
      }
      break;
  }
  return node;
}
