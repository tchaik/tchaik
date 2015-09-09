"use strict";

import {ChangeEmitter, CHANGE_EVENT} from "../utils/ChangeEmitter.js";

var AppDispatcher = require("../dispatcher/AppDispatcher");

var CollectionConstants = require("../constants/CollectionConstants.js");

var _commonFields = ["album", "albumArtist", "artist", "composer", "year", "bitRate", "discNumber"];

var _collections = {};

// pathSeparator is a string used to separate path components.
const pathSeparator = ":";

// pathToKey returns a string representation of the path.
function pathToKey(path) {
  if (path) {
    return path.join(pathSeparator);
  }
  return null;
}

function addItem(path, item) {
  if (item.tracks) { // fill in common fields if they are set
    item.tracks.forEach(function(track) {
      _commonFields.forEach(function(fld) {
        if (item[fld]) {
          track[fld] = item[fld];
        }
      });
      track.groupName = item.name;
    });
  }
  _collections[pathToKey(path)] = item;
}

class CollectionStore extends ChangeEmitter {
  pathToKey(path) {
    return pathToKey(path);
  }

  getCollection(path) {
    var key = pathToKey(path);
    return _collections[key];
  }

  emitChange(path) {
    this.emit(CHANGE_EVENT, pathToKey(path));
  }
}

var _store = new CollectionStore();

_store.dispatchToken = AppDispatcher.register(function(payload) {
  var action = payload.action;
  var source = payload.source;

  if (source === "SERVER_ACTION") {
    switch (action.actionType) {
      case CollectionConstants.FETCH:
        addItem(action.data.path, action.data.item);
        _store.emitChange(action.data.path);
        break;
    }
  }
  return true;
});

export default _store;
