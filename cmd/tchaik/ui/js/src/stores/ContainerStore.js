import {ChangeEmitter} from "../utils/ChangeEmitter.js";
import AppDispatcher from "../dispatcher/AppDispatcher.js";

import ContainerConstants from "../constants/ContainerConstants.js";
import SearchConstants from "../constants/SearchConstants.js";


const _defaultMode = ContainerConstants.ALL;

function setMode(m) {
  if (m === null) {
    m = _defaultMode;
  }
  localStorage.setItem("mode", m);
}

function mode() {
  let m = localStorage.getItem("mode");
  if (m === null) {
    m = _defaultMode;
  }
  return m;
}

const _titles = {};
_titles[ContainerConstants.ALL] = "Library";
_titles[ContainerConstants.RECENT] = "Recently Added";
_titles[ContainerConstants.RETRO] = "";

class ContainerStore extends ChangeEmitter {
  getMode() {
    return mode();
  }

  getTitle() {
    const m = mode();
    if (_titles.hasOwnProperty(m)) {
      return _titles[m];
    }
    return m.toLowerCase();
  }
}

const _containerStore = new ContainerStore();

_containerStore.dispatchToken = AppDispatcher.register(function(payload) {
  const action = payload.action;
  const source = payload.source;

  if (source === "VIEW_ACTION") {
    switch (action.actionType) {
      case ContainerConstants.MODE:
        setMode(action.mode);
        _containerStore.emitChange();
        break;

      case SearchConstants.SEARCH:
        if (_containerStore.getMode() !== ContainerConstants.SEARCH) {
          setMode(ContainerConstants.SEARCH);
          _containerStore.emitChange();
        }
        break;
    }
  }
  return true;
});

export default _containerStore;
