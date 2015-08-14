"use strict";

import AppDispatcher from "../dispatcher/AppDispatcher";

import {ChangeEmitter} from "../utils/ChangeEmitter.js";

import FavouriteConstants from "../constants/FavouriteConstants.js";


var _recent = [];

function setFavourite(recent) {
  if (recent !== null && recent.Groups) {
    _recent = recent.Groups;
    return;
  }
  _recent = [];
}

class FavouriteStore extends ChangeEmitter {
  getPaths() {
    return _recent;
  }
}

var _store = new FavouriteStore();

_store.dispatchToken = AppDispatcher.register(function(payload) {
  var action = payload.action;
  var source = payload.source;

  if (source === "SERVER_ACTION") {
    switch (action.actionType) {
      case FavouriteConstants.FETCH_FAVOURITE:
        setFavourite(action.data);
        _store.emitChange();
        break;
    }
  }

  return true;
});

export default _store;
