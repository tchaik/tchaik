import {ChangeEmitter} from "../utils/ChangeEmitter.js";
import AppDispatcher from "../dispatcher/AppDispatcher.js";

import ContainerConstants from "../constants/ContainerConstants.js";
import SearchConstants from "../constants/SearchConstants.js";


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

var _titles = {};
_titles[ContainerConstants.ALL] = "Library";
_titles[ContainerConstants.RETRO] = "";

class ContainerStore extends ChangeEmitter {
  getMode() {
    return mode();
  }

  getTitle() {
    var m = mode();
    if (_titles.hasOwnProperty(m)) {
      return _titles[m];
    }
    return m.toLowerCase();
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

      case SearchConstants.SEARCH:
        setMode(ContainerConstants.SEARCH);
        _containerStore.emitChange();
        break;
    }
  }
  return true;
});

export default _containerStore;
