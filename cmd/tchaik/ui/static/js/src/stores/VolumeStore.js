'use strict';

var AppDispatcher = require('../dispatcher/AppDispatcher');
var EventEmitter = require('eventemitter3').EventEmitter;
var assign = require('object-assign');

var CtrlConstants = require('../constants/ControlConstants.js');
var VolumeConstants = require('../constants/VolumeConstants.js');

var CHANGE_EVENT = 'change';

var defaultVolume = 0.75;
var defaultVolumeMute = false;

function setVolume(v) {
  localStorage.setItem("volume", v);
}

function volume() {
  var v = localStorage.getItem("volume");
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
  var v = localStorage.getItem("volumeMute");
  if (v === null) {
    setVolumeMute(defaultVolumeMute);
    return defaultVolumeMute;
  }
  return (v === "true");
}


var VolumeStore = assign({}, EventEmitter.prototype, {
  getVolume: function() {
    return volumeMute() ? 0.0 : volume();
  },

  getVolumeMute: function() {
    return volumeMute();
  },

  emitChange: function(type) {
    this.emit(CHANGE_EVENT, type);
  },

  /**
   * @param {function} callback
   */
  addChangeListener: function(callback) {
    this.on(CHANGE_EVENT, callback);
  },

  /**
   * @param {function} callback
   */
  removeChangeListener: function(callback) {
    this.removeListener(CHANGE_EVENT, callback);
  },
});

VolumeStore.dispatchToken = AppDispatcher.register(function(payload) {
  var action = payload.action;
  var source = payload.source;

  if (source === 'SERVER_ACTION') {
    if (action.actionType === CtrlConstants.CTRL) {
      switch (action.data) {

        case CtrlConstants.TOGGLE_MUTE:
          setVolumeMute(!volumeMute());
          VolumeStore.emitChange();
          break;

        default:
          break;
      }

      if (action.data.Key) {
        switch (action.data.Key) {
          case "volume":
            setVolume(action.data.Value);
            if (action.data.Value > 0) {
              setVolumeMute(false);
            }
            VolumeStore.emitChange();
            break;

          case "mute":
            setVolumeMute(action.data.Value);
            VolumeStore.emitChange();
            break;

          default:
            break;
        }
      }
    }
  }

  if (source === 'VIEW_ACTION') {
    switch (action.actionType) {

      case VolumeConstants.SET_VOLUME:
        setVolume(action.volume);
        if (action.volume > 0) {
          setVolumeMute(false);
        }
        VolumeStore.emitChange();
        break;

      case VolumeConstants.SET_VOLUME_MUTE:
        setVolumeMute(action.volumeMute);
        VolumeStore.emitChange();
        break;

      case VolumeConstants.TOGGLE_VOLUME_MUTE:
        setVolumeMute(!volumeMute());
        VolumeStore.emitChange();
        break;

      default:
        break;
    }
  }

  return true;
});

module.exports = VolumeStore;
