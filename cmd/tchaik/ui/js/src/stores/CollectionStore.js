"use strict";

import {ChangeEmitter, CHANGE_EVENT} from "../utils/ChangeEmitter.js";

import AppDispatcher from "../dispatcher/AppDispatcher";

import CollectionConstants from "../constants/CollectionConstants.js";


const _commonFields = ["album", "albumArtist", "artist", "composer", "year", "bitRate", "discNumber"];

const _collections = {};

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
    const key = pathToKey(path);
    return _collections[key];
  }

  emitChange(path) {
    this.emit(CHANGE_EVENT, pathToKey(path));
  }
}

const _store = new CollectionStore();

_store.dispatchToken = AppDispatcher.register(function(payload) {
  const action = payload.action;
  const source = payload.source;

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
