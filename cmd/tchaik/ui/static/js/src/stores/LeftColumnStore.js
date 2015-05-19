'use strict';

import {ChangeEmitter} from '../utils/ChangeEmitter.js';

import AppDispatcher from '../dispatcher/AppDispatcher';

import LeftColumnConstants from '../constants/LeftColumnConstants.js';


var _defaultMode = "All";
var _defaultHidden = true;

function setMode(m) {
  if (m === null) {
    m = _defaultMode;
  }
  localStorage.setItem("mode", m);
}

function setHidden(h) {
  if (h === null) {
    h = _defaultHidden;
  }

  localStorage.setItem('toolbarHidden', JSON.stringify(h));
}

function mode() {
  var m = localStorage.getItem("mode");
  if (m === null) {
    m = _defaultMode;
  }
  return m;
}

function isHidden() {
  var h = localStorage.getItem('toolbarHidden');
  if (h === null) {
    h = _defaultHidden;
  } else {
    h = JSON.parse(h);
  }

  return h;
}


class LeftColumnStore extends ChangeEmitter {
  getMode() {
    return mode();
  }

  getIsHidden() {
    return isHidden();
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
      case LeftColumnConstants.TOGGLE_VISIBILITY:
        var current = isHidden();
        setHidden(!current);
        _leftColumnStore.emitChange();
        break;
    }
  }
  return true;
});

export default _leftColumnStore;
