"use strict";

import {ChangeEmitter} from "../utils/ChangeEmitter.js";
import AppDispatcher from "../dispatcher/AppDispatcher";

import NowPlayingConstants from "../constants/NowPlayingConstants.js";
import PlaylistConstants from "../constants/PlaylistConstants.js";

import PlaylistStore from "./PlaylistStore.js";

import CtrlConstants from "../constants/ControlConstants.js";


var CONTROL_EVENT = "control";

var currentPlaying = null;

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
  localStorage.setItem("currentTrack", JSON.stringify(track));
}

function _playing() {
  var v = localStorage.getItem("playing");
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

function setPlaying(v) {
  currentPlaying = v;
  localStorage.setItem("playing", JSON.stringify(v));
}

function currentTrack() {
  var c = localStorage.getItem("currentTrack");
  if (c === null) {
    return null;
  }
  return JSON.parse(c);
}


class NowPlayingStore extends ChangeEmitter {
  getPlaying() {
    return playing();
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

function handlePrevAction() {
  AppDispatcher.waitFor([
    PlaylistStore.dispatchToken,
  ]);
  setCurrentTrack(PlaylistStore.getCurrentTrack());
  _nowPlayingStore.emitChange();
}

function handleNextAction() {
  AppDispatcher.waitFor([
    PlaylistStore.dispatchToken,
  ]);
  setCurrentTrack(PlaylistStore.getCurrentTrack());
  _nowPlayingStore.emitChange();
}

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
          handleNextAction();
          break;

        case CtrlConstants.PREV:
          handlePrevAction();
          break;

        case CtrlConstants.TOGGLE_PLAY_PAUSE:
          setPlaying(!_nowPlayingStore.getPlaying());
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
  }

  if (source === "VIEW_ACTION") {
    switch (action.actionType) {

      case PlaylistConstants.PREV:
        handlePrevAction();
        break;

      case NowPlayingConstants.ENDED:
        if (action.source !== "playlist") {
          break;
        }
        /* falls through */
      case PlaylistConstants.NEXT:
        handleNextAction();
        break;

      case NowPlayingConstants.SET_PLAYING:
        setPlaying(action.playing);
        _nowPlayingStore.emitChange();
        break;

      case NowPlayingConstants.SET_CURRENT_TIME:
        _nowPlayingStore.emitControl(NowPlayingConstants.SET_CURRENT_TIME, action.currentTime);
        break;

      case NowPlayingConstants.SET_CURRENT_TRACK:
        setCurrentTrack(action.track);
        setCurrentTrackSource(action.source);
        _nowPlayingStore.emitChange();
        break;

      case PlaylistConstants.ADD_ITEM_PLAY_NOW:
        AppDispatcher.waitFor([
          PlaylistStore.dispatchToken,
        ]);
        setCurrentTrack(PlaylistStore.getCurrentTrack());
        setCurrentTrackSource("playlist");
        _nowPlayingStore.emitChange();
        break;

      default:
        break;
    }
  }

  return true;
});

export default _nowPlayingStore;
