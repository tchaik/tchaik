'use strict';

var AppDispatcher = require('../dispatcher/AppDispatcher');
var EventEmitter = require('eventemitter3').EventEmitter;
var assign = require('object-assign');

var NowPlayingConstants = require('../constants/NowPlayingConstants.js');
var PlaylistConstants = require('../constants/PlaylistConstants.js');

var PlaylistStore = require('./PlaylistStore.js');

var CtrlConstants = require('../constants/ControlConstants.js');

var CHANGE_EVENT = 'change';
var CONTROL_EVENT = 'control';

var defaultVolume = 0.75;
var defaultVolumeMute = false;

var currentPlaying = null;

var _defaultTrackState = {
  buffered: 0.0,
  duration: 0.0,
};
var _trackState = _defaultTrackState;

function setCurrentTime(time) {
  localStorage.setItem("currentTime", time);
}

function currentTime() {
  var t = localStorage.getItem("currentTime");
  if (t === null) {
    return 0;
  }
  return parseFloat(t);
}

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
  setCurrentTime(0);
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

var NowPlayingStore = assign({}, EventEmitter.prototype, {

  getTime: function() {
    return currentTime();
  },

  getBuffered: function() {
    return _trackState.buffered;
  },

  getDuration: function() {
    return _trackState.duration;
  },

  getPlaying: function() {
    return playing();
  },

  getVolume: function() {
    return volumeMute() ? 0.0 : volume();
  },

  getVolumeMute: function() {
    return volumeMute();
  },

  getTrack: function() {
    return currentTrack();
  },

  getSource: function() {
    return currentTrackSource();
  },

  emitChange: function(type) {
    this.emit(CHANGE_EVENT, type);
  },

  emitControl: function(type, value) {
    this.emit(CONTROL_EVENT, type, value);
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

  /**
   * @param {function} callback
   */
  addControlListener: function(callback) {
    this.on(CONTROL_EVENT, callback);
  },

  /**
   * @param {function} callback
   */
  removeControlListener: function(callback) {
    this.removeListener(CONTROL_EVENT, callback);
  },

});

function handlePrevAction() {
  AppDispatcher.waitFor([
    PlaylistStore.dispatchToken,
  ]);
  setCurrentTrack(PlaylistStore.getCurrentTrack());
  NowPlayingStore.emitChange();
}

function handleNextAction() {
  AppDispatcher.waitFor([
    PlaylistStore.dispatchToken,
  ]);
  setCurrentTrack(PlaylistStore.getCurrentTrack());
  NowPlayingStore.emitChange();
}

NowPlayingStore.dispatchToken = AppDispatcher.register(function(payload) {
  var action = payload.action;
  var source = payload.source;

  if (source === 'SERVER_ACTION') {
    if (action.actionType === CtrlConstants.CTRL) {
      switch (action.data) {

        case CtrlConstants.PLAY:
          setPlaying(true);
          NowPlayingStore.emitChange();
          break;

        case CtrlConstants.PAUSE:
          setPlaying(false);
          NowPlayingStore.emitChange();
          break;

        case CtrlConstants.NEXT:
          handleNextAction();
          break;

        case CtrlConstants.PREV:
          handlePrevAction();
          break;

        case CtrlConstants.TOGGLE_PLAY_PAUSE:
          setPlaying(!NowPlayingStore.getPlaying());
          NowPlayingStore.emitChange();
          break;

        case CtrlConstants.TOGGLE_MUTE:
          setVolumeMute(!volumeMute());
          NowPlayingStore.emitChange();
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
            NowPlayingStore.emitChange();
            break;

          case "mute":
            setVolumeMute(action.data.Value);
            NowPlayingStore.emitChange();
            break;

          case "time":
            NowPlayingStore.emitControl(NowPlayingConstants.SET_CURRENT_TIME, action.data.Value);
            break;

          default:
            console.log("Unknown key:", action.data.Key);
            break;
        }
      }
    }
  }

  if (source === 'VIEW_ACTION') {
    switch (action.actionType) {

      case NowPlayingConstants.RESET:
        _trackState = {
          buffered: 0,
          duration: 0,
        };
        NowPlayingStore.emitChange();
        break;

      case NowPlayingConstants.SET_DURATION:
        _trackState.duration = action.duration;
        NowPlayingStore.emitChange();
        break;

      case NowPlayingConstants.SET_BUFFERED:
        _trackState.buffered = action.buffered;
        NowPlayingStore.emitChange();
        break;

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
        NowPlayingStore.emitChange();
        break;

      case NowPlayingConstants.STORE_CURRENT_TIME:
        setCurrentTime(action.currentTime);
        NowPlayingStore.emitChange();
        break;

      case NowPlayingConstants.SET_CURRENT_TIME:
        NowPlayingStore.emitControl(NowPlayingConstants.SET_CURRENT_TIME, action.currentTime);
        break;

      case NowPlayingConstants.SET_VOLUME:
        setVolume(action.volume);
        if (action.volume > 0) {
          setVolumeMute(false);
        }
        NowPlayingStore.emitChange();
        break;

      case NowPlayingConstants.SET_VOLUME_MUTE:
        setVolumeMute(action.volumeMute);
        NowPlayingStore.emitChange();
        break;

      case NowPlayingConstants.TOGGLE_VOLUME_MUTE:
        setVolumeMute(!volumeMute());
        NowPlayingStore.emitChange();
        break;

      case NowPlayingConstants.SET_CURRENT_TRACK:
        setCurrentTrack(action.track);
        setCurrentTrackSource(action.source);
        NowPlayingStore.emitChange();
        break;

      case PlaylistConstants.PLAY_NOW:
        AppDispatcher.waitFor([
          PlaylistStore.dispatchToken,
        ]);
        setCurrentTrack(PlaylistStore.getCurrentTrack());
        NowPlayingStore.emitChange();
        break;

      default:
        break;
    }
  }

  return true;
});

module.exports = NowPlayingStore;
