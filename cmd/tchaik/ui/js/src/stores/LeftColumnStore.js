"use strict";

import {ChangeEmitter} from "../utils/ChangeEmitter.js";

import AppDispatcher from "../dispatcher/AppDispatcher";

import LeftColumnConstants from "../constants/LeftColumnConstants.js";


const _defaultHidden = true;

function setHidden(h) {
  if (h === null) {
    h = _defaultHidden;
  }
  localStorage.setItem("leftColumnHidden", JSON.stringify(h));
}

function isHidden() {
  let h = localStorage.getItem("leftColumnHidden");
  if (h === null) {
    h = _defaultHidden;
  } else {
    h = JSON.parse(h);
  }
  return h;
}


class LeftColumnStore extends ChangeEmitter {
  getIsHidden() {
    return isHidden();
  }
}

const _store = new LeftColumnStore();

_store.dispatchToken = AppDispatcher.register(function(payload) {
  const action = payload.action;
  const source = payload.source;

  if (source === "VIEW_ACTION") {
    switch (action.actionType) {
      case LeftColumnConstants.TOGGLE_LEFTCOLUMN:
        const current = isHidden();
        setHidden(!current);
        _store.emitChange();
        break;
    }
  }
  return true;
});

export default _store;
