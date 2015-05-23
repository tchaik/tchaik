import {ChangeEmitter} from "../utils/ChangeEmitter.js";
import AppDispatcher from "../dispatcher/AppDispatcher.js";
import ContainerConstants from "../constants/ContainerConstants.js";

var _defaultMode = ContainerConstants.ALL;

function setMode(m) {
  if (m === null) {
    m = _defaultMode;
  }
  localStorage.setItem("mode", m);
}

function mode() {
  var m = localStorage.getItem("mode");
  if (m === null) {
    m = _defaultMode;
  }
  return m;
}

class ContainerStore extends ChangeEmitter {
  getMode() {
    return mode();
  }
}

var _containerStore = new ContainerStore();

_containerStore.dispatchToken = AppDispatcher.register(function(payload) {
  var action = payload.action;
  var source = payload.source;

  if (source === "VIEW_ACTION") {
    switch (action.actionType) {
      case ContainerConstants.MODE:
        setMode(action.mode);
        _containerStore.emitChange();
        break;
    }
  }
  return true;
});

export default _containerStore;
