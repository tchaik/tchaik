'use strict';

var AppDispatcher = require('../dispatcher/AppDispatcher');

var NowPlayingConstants = require('../constants/NowPlayingConstants.js');

var NowPlayingActions = {

  remove: function(index) {
    AppDispatcher.handleViewAction({
      actionType: NowPlayingConstants.REMOVE,
      index: index,
    });
  },

  currentTime: function(time) {
    AppDispatcher.handleViewAction({
      actionType: NowPlayingConstants.STORE_CURRENT_TIME,
      currentTime: time,
    });
  },

  setCurrentTime: function(time) {
    AppDispatcher.handleViewAction({
      actionType: NowPlayingConstants.SET_CURRENT_TIME,
      currentTime: time,
    });
  },

  ended: function(source) {
    AppDispatcher.handleViewAction({
      actionType: NowPlayingConstants.ENDED,
      source: source,
    });
  },

  volume: function(v) {
    AppDispatcher.handleViewAction({
      actionType: NowPlayingConstants.SET_VOLUME,
      volume: v,
    });
  },

  toggleVolumeMute: function() {
    AppDispatcher.handleViewAction({
      actionType: NowPlayingConstants.TOGGLE_VOLUME_MUTE,
    });
  },

  playing: function(v) {
    AppDispatcher.handleViewAction({
      actionType: NowPlayingConstants.SET_PLAYING,
      playing: v,
    });
  },

};

module.exports = NowPlayingActions;
