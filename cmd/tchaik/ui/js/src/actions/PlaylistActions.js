import AppDispatcher from "../dispatcher/AppDispatcher";

import PlaylistConstants from "../constants/PlaylistConstants.js";
import WebsocketAPI from "../utils/WebsocketAPI.js";

var playlistName = "Default";


var PlaylistActions = {

  fetch: function() {
    WebsocketAPI.send(PlaylistConstants.PLAYLIST, {
      action: PlaylistConstants.FETCH,
      name: playlistName,
    });
  },

  addItem: function(path) {
    WebsocketAPI.send(PlaylistConstants.PLAYLIST, {
      action: PlaylistConstants.ADD_ITEM,
      name: playlistName,
      path: path,
    });
  },

  addItemPlayNow: function(path) {
    WebsocketAPI.send(PlaylistConstants.PLAYLIST, {
      action: PlaylistConstants.ADD_ITEM,
      name: playlistName,
      path: path,
    });
  },

  remove: function(itemIndex, path) {
    WebsocketAPI.send(PlaylistConstants.PLAYLIST, {
      name: playlistName,
      action: PlaylistConstants.REMOVE,
      index: itemIndex,
      path: path,
    });
  },

  clear: function() {
    AppDispatcher.handleViewAction({
      actionType: PlaylistConstants.CLEAR_PLAYLIST,
    });
  },

};

export default PlaylistActions;
