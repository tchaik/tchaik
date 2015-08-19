"use strict";

import AppDispatcher from "../dispatcher/AppDispatcher";

import {ChangeEmitter} from "../utils/ChangeEmitter.js";

import ChecklistConstants from "../constants/ChecklistConstants.js";


var _recent = [];

function setChecklist(recent) {
  if (recent !== null && recent.Groups) {
    _recent = recent.Groups;
    return;
  }
  _recent = [];
}

class ChecklistStore extends ChangeEmitter {
  getPaths() {
    return _recent;
  }
}

var _store = new ChecklistStore();

_store.dispatchToken = AppDispatcher.register(function(payload) {
  var action = payload.action;
  var source = payload.source;

  if (source === "SERVER_ACTION") {
    switch (action.actionType) {
      case ChecklistConstants.FETCH_CHECKLIST:
        setChecklist(action.data);
        _store.emitChange();
        break;
    }
  }

  return true;
});

export default _store;
