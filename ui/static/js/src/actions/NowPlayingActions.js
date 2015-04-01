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
      actionType: NowPlayingConstants.SET_CURRENT_TIME,
      currentTime: time,
    });
  },

  restartTrack: function() {
    AppDispatcher.handleViewAction({
      actionType: NowPlayingConstants.SET_CURRENT_TIME,
      currentTime: -1,
    });
  },

  ended: function() {
    AppDispatcher.handleViewAction({
      actionType: NowPlayingConstants.ENDED,
    });
  },

  volume: function(v) {
    AppDispatcher.handleViewAction({
      actionType: NowPlayingConstants.SET_VOLUME,
      volume: v,
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
