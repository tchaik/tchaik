'use strict';

var AppDispatcher = require('../dispatcher/AppDispatcher');
var EventEmitter = require('eventemitter3').EventEmitter;
var assign = require('object-assign');

var NowPlayingConstants = require('../constants/NowPlayingConstants.js');
var PlaylistConstants = require('../constants/PlaylistConstants.js');

var NowPlayingStore = require('../stores/NowPlayingStore.js');

// var CtrlConstants = require('../constants/ControlConstants.js');

var CHANGE_EVENT = 'change';

var _defaultTrackState = {
  buffered: 0.0,
  duration: 0.0,
};
var _trackState = _defaultTrackState;

var _currentTime = null;

function setCurrentTime(time) {
  _currentTime = time;
  localStorage.setItem("currentTime", time);
}

function currentTime() {
  if (_currentTime !== null) {
    return _currentTime;
  }

  _currentTime = 0;
  var t = localStorage.getItem("currentTime");
  if (t !== null) {
    _currentTime = parseFloat(t);
  }
  return _currentTime;
}


var PlayingStatusStore = assign({}, EventEmitter.prototype, {

  getTime: function() {
    return currentTime();
  },

  getBuffered: function() {
    return _trackState.buffered;
  },

  getDuration: function() {
    return _trackState.duration;
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


PlayingStatusStore.dispatchToken = AppDispatcher.register(function(payload) {
  var action = payload.action;
  var source = payload.source;

  if (source === 'VIEW_ACTION') {
    switch (action.actionType) {

    case NowPlayingConstants.ENDED:
      if (action.source !== "playlist") {
        break;
      }
      /* falls through */
    case PlaylistConstants.PREV:
      /* falls through */
    case PlaylistConstants.NEXT:
      /* falls through */
    case PlaylistConstants.PLAY_NOW:
      /* falls through */
    case NowPlayingConstants.SET_CURRENT_TRACK:
        AppDispatcher.waitFor([
          NowPlayingStore.dispatchToken,
        ]);
        setCurrentTime(0);
        PlayingStatusStore.emitChange();
        break;

      case NowPlayingConstants.RESET:
        _trackState = {
          buffered: 0,
          duration: 0,
        };
        PlayingStatusStore.emitChange();
        break;

      case NowPlayingConstants.SET_DURATION:
        _trackState.duration = action.duration;
        PlayingStatusStore.emitChange();
        break;

      case NowPlayingConstants.SET_BUFFERED:
        _trackState.buffered = action.buffered;
        PlayingStatusStore.emitChange();
        break;

      case NowPlayingConstants.STORE_CURRENT_TIME:
        setCurrentTime(action.currentTime);
        PlayingStatusStore.emitChange();
        break;

      default:
        break;
    }
  }

  return true;
});

module.exports = PlayingStatusStore;
