'use strict';

var AppDispatcher = require('../dispatcher/AppDispatcher');
var EventEmitter = require('eventemitter3').EventEmitter;
var assign = require('object-assign');

var NowPlayingConstants = require('../constants/NowPlayingConstants.js');
var PlaylistConstants = require('../constants/PlaylistConstants.js');

var PlaylistStore = require('./PlaylistStore.js');

var CHANGE_EVENT = 'change';

var defaultVolume = 0.75;

var currentPlaying = null;

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

var NowPlayingStore = assign({}, EventEmitter.prototype, {

  getCurrentTime: function() {
    return currentTime();
  },

  getPlaying: function() {
    return playing();
  },

  getVolume: function() {
    return volume();
  },

  getCurrent: function() {
    return currentTrack();
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
  }

});

NowPlayingStore.dispatchToken = AppDispatcher.register(function(payload) {
  var action = payload.action;
  var source = payload.source;

  if (source === 'VIEW_ACTION') {
    switch (action.actionType) {

      case PlaylistConstants.PREV:
        AppDispatcher.waitFor([
          PlaylistStore.dispatchToken,
        ]);
        setCurrentTrack(PlaylistStore.getCurrentTrack());
        NowPlayingStore.emitChange();
        break;

      case NowPlayingConstants.ENDED:
        /* falls through */
      case PlaylistConstants.NEXT:
        AppDispatcher.waitFor([
          PlaylistStore.dispatchToken,
        ]);
        setCurrentTrack(PlaylistStore.getCurrentTrack());
        NowPlayingStore.emitChange();
        break;

      case NowPlayingConstants.SET_PLAYING:
        setPlaying(action.playing);
        NowPlayingStore.emitChange();
        break;

      case NowPlayingConstants.SET_CURRENT_TIME:
        setCurrentTime(action.currentTime);
        break;

      case NowPlayingConstants.SET_VOLUME:
        setVolume(action.volume);
        NowPlayingStore.emitChange();
        break;

      case NowPlayingConstants.SET_CURRENT_TRACK:
        setCurrentTrack(action.track);
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
