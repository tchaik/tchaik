'use strict';

var AppDispatcher = require('../dispatcher/AppDispatcher');

var WebsocketAPI = require('../utils/WebsocketAPI.js');

var CollectionConstants = require('../constants/CollectionConstants.js');
var NowPlayingConstants = require('../constants/NowPlayingConstants.js');

var CollectionActions = {

  fetch: function(path) {
    WebsocketAPI.send({
      path: path,
      action: CollectionConstants.FETCH,
    });
  },

  expandPath: function(path, expand) {
    AppDispatcher.handleViewAction({
      actionType: CollectionConstants.EXPAND_PATH,
      path: path,
      expand: expand,
    });
  },

  setCurrentTrack: function(track) {
    AppDispatcher.handleViewAction({
      actionType: NowPlayingConstants.SET_CURRENT_TRACK,
      track: track,
      source: "collection",
    });
  },

  appendToPlaylist: function(path) {
    AppDispatcher.handleViewAction({
      actionType: CollectionConstants.APPEND_TO_PLAYLIST,
      path: path
    });
  },

  playNow: function(path) {
    AppDispatcher.handleViewAction({
      actionType: CollectionConstants.PLAY_NOW,
      path: path
    });
  }

};

module.exports = CollectionActions;
