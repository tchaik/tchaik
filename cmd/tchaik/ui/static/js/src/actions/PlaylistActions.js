import AppDispatcher from "../dispatcher/AppDispatcher";

import PlaylistConstants from "../constants/PlaylistConstants.js";
import NowPlayingConstants from "../constants/NowPlayingConstants.js";
import WebsocketAPI from "../utils/WebsocketAPI.js";


var PlaylistActions = {

  fetch: function(path) {
    WebsocketAPI.send(PlaylistConstants.FETCH, {path: path});
  },

  addItem: function(path) {
    AppDispatcher.handleViewAction({
      actionType: PlaylistConstants.ADD_ITEM,
      path: path,
    });
  },

  addItemPlayNow: function(path) {
    AppDispatcher.handleViewAction({
      actionType: PlaylistConstants.ADD_ITEM_PLAY_NOW,
      path: path,
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
      source: "playlist",
    });
  },

};

export default PlaylistActions;
