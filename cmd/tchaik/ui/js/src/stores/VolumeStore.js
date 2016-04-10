"use strict";

import {ChangeEmitter} from "../utils/ChangeEmitter.js";

import AppDispatcher from "../dispatcher/AppDispatcher";

import CtrlConstants from "../constants/ControlConstants.js";
import VolumeConstants from "../constants/VolumeConstants.js";


const defaultVolume = 0.75;
const defaultVolumeMute = false;

function _setVolume(v) {
  localStorage.setItem("volume", v);
}

function volume() {
  const v = localStorage.getItem("volume");
  if (v === null) {
    setVolume(defaultVolume);
    return defaultVolume;
  }
  return parseFloat(v);
}

function setVolumeMute(v) {
  localStorage.setItem("volumeMute", v);
}

function volumeMute() {
  const v = localStorage.getItem("volumeMute");
  if (v === null) {
    setVolumeMute(defaultVolumeMute);
    return defaultVolumeMute;
  }
  return (v === "true");
}

function setVolume(v) {
  _setVolume(v);
  if (v > 0) {
    setVolumeMute(false);
  }
}


class VolumeStore extends ChangeEmitter {
  getVolume() {
    return volumeMute() ? 0.0 : volume();
  }

  getVolumeMute() {
    return volumeMute();
  }
}

const _volumeStore = new VolumeStore();

_volumeStore.dispatchToken = AppDispatcher.register(function(payload) {
  const action = payload.action;
  const source = payload.source;

  if (source === "SERVER_ACTION") {
    if (action.actionType === CtrlConstants.CTRL) {
      switch (action.data.action) {

        case CtrlConstants.TOGGLE_MUTE:
          setVolumeMute(!volumeMute());
          _volumeStore.emitChange();
          break;

        case "volume":
          setVolume(action.data.Value);
          _volumeStore.emitChange();
          break;

        case "mute":
          setVolumeMute(action.data.Value);
          _volumeStore.emitChange();
          break;

        default:
          break;
      }
    }
  }

  if (source === "VIEW_ACTION") {
    switch (action.actionType) {

      case VolumeConstants.SET_VOLUME:
        setVolume(action.volume);
        _volumeStore.emitChange();
        break;

      case VolumeConstants.TOGGLE_VOLUME_MUTE:
        setVolumeMute(!volumeMute());
        _volumeStore.emitChange();
        break;

      default:
        break;
    }
  }

  return true;
});

export default _volumeStore;
