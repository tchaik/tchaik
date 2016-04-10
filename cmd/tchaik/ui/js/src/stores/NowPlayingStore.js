"use strict";

import {ChangeEmitter} from "../utils/ChangeEmitter.js";
import AppDispatcher from "../dispatcher/AppDispatcher";

import NowPlayingConstants from "../constants/NowPlayingConstants.js";
import PlaylistConstants from "../constants/PlaylistConstants.js";
import CursorConstants from "../constants/CursorConstants.js";

import CursorStore from "./CursorStore.js";
import PlaylistStore from "./PlaylistStore.js";

import CtrlConstants from "../constants/ControlConstants.js";


var CONTROL_EVENT = "control";

var currentPlaying = null;
var currentRepeat = null;
var _currentTrack = null;

function setCurrentTrackSource(source) {
  localStorage.setItem("currentTrackSource", source);
}

function currentTrackSource() {
  var s = localStorage.getItem("currentTrackSource");
  if (s === null) {
    return null;
  }
  return s;
}

function setCurrentTrack(track) {
  let change = true;
  if (_currentTrack !== null) {
    change = (_currentTrack.id !== track.id);
  }
  localStorage.setItem("currentTrack", JSON.stringify(track));
  _currentTrack = track;
  return change;
}

function _playing() {
  var v = localStorage.getItem("playing");
  if (v === null) {
    return false;
  }
  return JSON.parse(v);
}

function _repeat() {
  var v = localStorage.getItem("repeat");
  if (v === null) {
    return false;
  }
  return JSON.parse(v);
}

function playing() {
  if (currentPlaying === null) {
    currentPlaying = _playing();
  }
  return currentPlaying;
}

function repeat() {
  if (currentRepeat === null) {
    currentRepeat = _repeat();
  }
  return currentRepeat;
}

function setPlaying(v) {
  currentPlaying = v;
  localStorage.setItem("playing", JSON.stringify(v));
}

function setRepeat(v) {
  currentRepeat = v;
  localStorage.setItem("repeat", JSON.stringify(v));
}

function currentTrack() {
  if (_currentTrack === null) {
    var c = localStorage.getItem("currentTrack");
    if (c === null) {
      return null;
    }
    _currentTrack = JSON.parse(c);
  }
  return _currentTrack;
}


class NowPlayingStore extends ChangeEmitter {
  getPlaying() {
    return playing();
  }

  getRepeat() {
    return repeat();
  }

  getTrack() {
    return currentTrack();
  }

  getSource() {
    return currentTrackSource();
  }

  emitControl(type, value) {
    this.emit(CONTROL_EVENT, type, value);
  }

  /**
   * @param {function} callback
   */
  addControlListener(callback) {
    this.on(CONTROL_EVENT, callback);
  }

  /**
   * @param {function} callback
   */
  removeControlListener(callback) {
    this.removeListener(CONTROL_EVENT, callback);
  }
}

var _nowPlayingStore = new NowPlayingStore();

_nowPlayingStore.dispatchToken = AppDispatcher.register(function(payload) {
  var action = payload.action;
  var source = payload.source;

  if (source === "SERVER_ACTION") {
    if (action.actionType === CtrlConstants.CTRL) {
      switch (action.data.action) {

        case CtrlConstants.PLAY:
          setPlaying(true);
          _nowPlayingStore.emitChange();
          break;

        case CtrlConstants.PAUSE:
          setPlaying(false);
          _nowPlayingStore.emitChange();
          break;

        case CtrlConstants.NEXT:
          /* fallsthrough */
        case CtrlConstants.PREV:
          AppDispatcher.waitFor([
            CursorStore.dispatchToken,
          ]);
          if (setCurrentTrack(CursorStore.getCurrentTrack())) {
            _nowPlayingStore.emitChange();
          }
          break;

        case CtrlConstants.TOGGLE_PLAY_PAUSE:
          setPlaying(!_nowPlayingStore.getPlaying());
          _nowPlayingStore.emitChange();
          break;

        case CtrlConstants.TOGGLE_REPEAT:
          setRepeat(!_nowPlayingStore.getRepeat());
          _nowPlayingStore.emitChange();
          break;

        case "repeat":
          setRepeat(action.data.Value);
          _nowPlayingStore.emitChange();
          break;

        case "time":
          _nowPlayingStore.emitControl(NowPlayingConstants.SET_CURRENT_TIME, action.data.Value);
          break;

        default:
          console.log("Unknown action:", action.data.action);
          break;
      }
    }

    if (action.actionType === CursorConstants.CURSOR) {
      AppDispatcher.waitFor([
        CursorStore.dispatchToken,
      ]);
      setCurrentTrackSource("cursor");
      if (setCurrentTrack(CursorStore.getCurrentTrack())) {
        _nowPlayingStore.emitChange();
      }
    }
  }

  if (source === "VIEW_ACTION") {
    switch (action.actionType) {

      case NowPlayingConstants.ENDED:
        if (action.repeat === true) {
          setRepeat(false);
          _nowPlayingStore.emitControl(NowPlayingConstants.SET_CURRENT_TIME, 0);
          setPlaying(true);
          _nowPlayingStore.emitChange();
          break;
        }

        if (action.source !== "cursor") {
          break;
        }
        /* falls through */
      case CursorConstants.PREV:
        /* falls through */
      case CursorConstants.NEXT:
        AppDispatcher.waitFor([
          CursorStore.dispatchToken,
        ]);
        setCurrentTrack(CursorStore.getCurrentTrack());
        _nowPlayingStore.emitChange();
        break;

      case NowPlayingConstants.SET_PLAYING:
        setPlaying(action.playing);
        _nowPlayingStore.emitChange();
        break;

      case NowPlayingConstants.SET_REPEAT:
        setRepeat(action.repeat);
        _nowPlayingStore.emitChange();
        break;

      case NowPlayingConstants.SET_CURRENT_TIME:
        _nowPlayingStore.emitControl(NowPlayingConstants.SET_CURRENT_TIME, action.currentTime);
        break;

      case CursorConstants.SET:
        AppDispatcher.waitFor([
          CursorStore.dispatchToken,
        ]);
        setCurrentTrackSource("cursor");
        let change = setCurrentTrack(CursorStore.getCurrentTrack());
        if (!playing()) {
          change = true;
          setPlaying(true);
        }
        if (change) {
          _nowPlayingStore.emitChange();
        }
        break;

      case NowPlayingConstants.SET_CURRENT_TRACK:
        setCurrentTrack(action.track);
        setCurrentTrackSource("collection");
        _nowPlayingStore.emitChange();
        break;

      case PlaylistConstants.ADD_ITEM_PLAY_NOW:
        AppDispatcher.waitFor([
          PlaylistStore.dispatchToken,
        ]);
        console.warn("Not implemented.");
        _nowPlayingStore.emitChange();
        break;

      default:
        break;
    }
  }

  return true;
});

export default _nowPlayingStore;
