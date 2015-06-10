"use strict";

import {ChangeEmitter, CHANGE_EVENT} from "../utils/ChangeEmitter.js";

var AppDispatcher = require("../dispatcher/AppDispatcher");

var CollectionConstants = require("../constants/CollectionConstants.js");

var _commonFields = ["Album", "AlbumArtist", "Artist", "Composer", "Year", "BitRate", "DiscNumber"];

var _collections = {};

// pathToKey returns a string representation of the path.  The only requirement is that
// subpaths should be prefixes.
function pathToKey(path) {
  if (path) {
    return path.join(">>");
  }
  return null;
}

function addItem(path, item) {
  if (item.Tracks) { // fill in common fields if they are set
    item.Tracks.forEach(function(track) {
      _commonFields.forEach(function(fld) {
        if (item[fld]) {
          track[fld] = item[fld];
        }
      });
      track.GroupName = item.Name;
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
        addItem(action.data.Path, action.data.Item);
        _store.emitChange(action.data.Path);
        break;
    }
  }
  return true;
});

export default _store;
