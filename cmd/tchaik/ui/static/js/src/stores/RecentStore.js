"use strict";

import AppDispatcher from "../dispatcher/AppDispatcher";

import {ChangeEmitter} from "../utils/ChangeEmitter.js";

import RecentConstants from "../constants/RecentConstants.js";


var _recent = [];

function setRecent(recent) {
  if (recent !== null && recent.Groups) {
    _recent = recent.Groups;
    return;
  }
  _recent = [];
}

class RecentStore extends ChangeEmitter {
  getPaths() {
    return _recent;
  }
}

var _store = new RecentStore();

_store.dispatchToken = AppDispatcher.register(function(payload) {
  var action = payload.action;
  var source = payload.source;

  if (source === "SERVER_ACTION") {
    switch (action.actionType) {
      case RecentConstants.FETCH_RECENT:
        setRecent(action.data);
        _store.emitChange();
        break;
    }
  }

  return true;
});

export default _store;
