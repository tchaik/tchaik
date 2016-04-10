"use strict";

import {ChangeEmitter} from "../utils/ChangeEmitter.js";

import AppDispatcher from "../dispatcher/AppDispatcher";

import RightColumnConstants from "../constants/RightColumnConstants.js";


const _defaultHidden = true;

function setHidden(h) {
  if (h === null) {
    h = _defaultHidden;
  }
  localStorage.setItem("rightColumnHidden", JSON.stringify(h));
}

function isHidden() {
  let h = localStorage.getItem("rightColumnHidden");
  if (h === null) {
    h = _defaultHidden;
  } else {
    h = JSON.parse(h);
  }
  return h;
}


class RightColumnStore extends ChangeEmitter {
  getIsHidden() {
    return isHidden();
  }
}

const _store = new RightColumnStore();

_store.dispatchToken = AppDispatcher.register(function(payload) {
  const action = payload.action;
  const source = payload.source;

  if (source === "VIEW_ACTION") {
    switch (action.actionType) {
      case RightColumnConstants.TOGGLE_RIGHTCOLUMN:
        const current = isHidden();
        setHidden(!current);
        _store.emitChange();
        break;
    }
  }
  return true;
});

export default _store;
