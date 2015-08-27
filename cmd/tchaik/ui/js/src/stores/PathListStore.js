"use strict";

import AppDispatcher from "../dispatcher/AppDispatcher";

import {ChangeEmitter} from "../utils/ChangeEmitter.js";

import PathListConstants from "../constants/PathListConstants.js";


var _pathLists = {};

function setPathList(name, list) {
  if (list !== null && list.Groups) {
    _pathLists[name] = list.Groups;
    return;
  }
  delete _pathLists[name];
}

class PathListStore extends ChangeEmitter {
  getPaths(name) {
    if (!_pathLists[name]) {
      _pathLists[name] = [];
    }
    return _pathLists[name];
  }
}

var _store = new PathListStore();

_store.dispatchToken = AppDispatcher.register(function(payload) {
  var action = payload.action;
  var source = payload.source;

  if (source === "SERVER_ACTION") {
    switch (action.actionType) {
      case PathListConstants.FETCH_PATHLIST:
        setPathList(action.data.Name, action.data.Data);
        _store.emitChange();
        break;
    }
  }

  return true;
});

export default _store;
