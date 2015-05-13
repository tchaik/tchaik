'use strict';

import ChangeEmitter from '../utils/ChangeEmitter.js';

import AppDispatcher from '../dispatcher/AppDispatcher';

import LeftColumnConstants from '../constants/LeftColumnConstants.js';


var _defaultMode = "All";

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


class LeftColumnStore extends ChangeEmitter {
  getMode() {
    return mode();
  }
}

var _leftColumnStore = new LeftColumnStore();

_leftColumnStore.dispatchToken = AppDispatcher.register(function(payload) {
  var action = payload.action;
  var source = payload.source;

  if (source === 'VIEW_ACTION') {
    switch (action.actionType) {
      case LeftColumnConstants.MODE:
        setMode(action.mode);
        _leftColumnStore.emitChange();
        break;
    }
  }
  return true;
});

export default _leftColumnStore;
