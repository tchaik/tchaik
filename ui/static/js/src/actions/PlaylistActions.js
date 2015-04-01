'use strict';

var AppDispatcher = require('../dispatcher/AppDispatcher');

var PlaylistConstants = require('../constants/PlaylistConstants.js');
var NowPlayingConstants = require('../constants/NowPlayingConstants.js');
var WebsocketApi = require('../api/WebsocketApi.js');

var PlaylistActions = {

  fetch: function(path) {
    WebsocketApi.send({
      path: path,
      action: PlaylistConstants.FETCH,
    });
  },

  remove: function(itemIndex, path) {
    AppDispatcher.handleViewAction({
      actionType: PlaylistConstants.REMOVE,
      itemIndex: itemIndex,
      path: path,
    });
  },

  next: function() {
    AppDispatcher.handleViewAction({
      actionType: PlaylistConstants.NEXT,
    });
  },

  prev: function() {
    AppDispatcher.handleViewAction({
      actionType: PlaylistConstants.PREV,
    });
  },

  play: function(itemIndex, path, data) {
    AppDispatcher.handleViewAction({
      actionType: PlaylistConstants.PLAY_ITEM,
      itemIndex: itemIndex,
      path: path,
    });

    AppDispatcher.handleViewAction({
      actionType: NowPlayingConstants.SET_CURRENT_TRACK,
      track: data,
    });
  },

};

module.exports = PlaylistActions;
