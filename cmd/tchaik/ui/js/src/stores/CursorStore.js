"use strict";

import AppDispatcher from "../dispatcher/AppDispatcher";
import {ChangeEmitter} from "../utils/ChangeEmitter.js";

import CollectionStore from "./CollectionStore.js";

import CursorConstants from "../constants/CursorConstants.js";
import NowPlayingConstants from "../constants/NowPlayingConstants.js";
import ControlConstants from "../constants/ControlConstants.js";


function trackForPath(path) {
  path = path.slice(0);
  var i = path.pop();
  var t = CollectionStore.getCollection(path);
  if (t === null) {
    console.log("Could not find collection with path:" + path);
  }

  if (t.Tracks) {
    var track = t.Tracks[i];
    if (track) {
      return track;
    }
    console.log("No track found for index: " + i);
    console.log(t.Tracks);
    return null;
  }

  console.log("Collection item did not have Tracks property");
  return null;
}

class Position {
  constructor(index, path) {
    this._index = index;
    this._path = path;
  }

  path() {
    return this._path;
  }

  index() {
    return this._index;
  }

  isEmpty() {
    return this._path === null;
  }
}

class Cursor {
  constructor(current = new Position(0, null), previous= new Position(0, null), next = new Position(0, null)) {
    this._current = current;
    this._previous = previous;
    this._next = next;
  }

  current() {
    return this._current;
  }

  currentTrack() {
    if (this._current.isEmpty()) {
      return null;
    }
    return trackForPath(this._current.path());
  }

  forward() {
    if (this._next.isEmpty()) {
      return;
    }
    this._previous = this._current;
    this._current = this._next;
    this._next = new Position(0, null);
  }

  backward() {
    if (this._previous.isEmpty()) {
      return;
    }
    this._next = this._current;
    this._current = this._previous;
    this._previous = new Position(0, null);
  }

  set(index, path) {
    this._current = new Position(index, path);
    this._previous = new Position(0, null);
    this._next = new Position(0, null);
  }

  canForward() {
    return !this._next.isEmpty();
  }

  canBackward() {
    return !this._previous.isEmpty();
  }
}


var _cursor = new Cursor();

class CursorStore extends ChangeEmitter {

  getCurrentPosition() {
    return _cursor.position();
  }

  getCurrent() {
    return _cursor.current();
  }

  getCurrentTrack() {
    return _cursor.currentTrack();
  }

  canPrev() {
    return _cursor.canBackward();
  }

  canNext() {
    return _cursor.canForward();
  }

}

function positionFromData(position) {
  return new Position(position.index, position.path);
}

var _store = new CursorStore();

_store.dispatchToken = AppDispatcher.register(function(payload) {
  var action = payload.action;
  var source = payload.source;

  if (source === "SERVER_ACTION") {
    if (action.actionType === ControlConstants.CTRL) {
      switch (action.data.action) {
        case ControlConstants.NEXT:
          // TODO: Implement this
          break;

        case ControlConstants.PREV:
          // TODO: Implement this
          break;
      }
    }

    if (action.actionType === CursorConstants.CURSOR) {
      let current = positionFromData(action.data.current);
      let previous = positionFromData(action.data.previous);
      let next = positionFromData(action.data.next);

      _cursor = new Cursor(current, previous, next);
      _store.emitChange();
    }
  }

  if (source === "VIEW_ACTION") {
    switch (action.actionType) {

      case NowPlayingConstants.ENDED:
        if (action.repeat === true) {
          break;
        }
        if (action.source !== "cursor") {
          break;
        }
        /* falls through */
      case CursorConstants.NEXT:
        _cursor.forward();
        _store.emitChange();
        break;

      case CursorConstants.PREV:
        _cursor.backward();
        _store.emitChange();
        break;

      case CursorConstants.SET:
        _cursor.set(action.itemIndex, action.path);
        _store.emitChange();
        break;

      default:
        break;
    }
  }

  return true;
});

export default _store;
