import AppDispatcher from "../dispatcher/AppDispatcher";

import WebsocketAPI from "../utils/WebsocketAPI.js";

import CollectionStore from "../stores/CollectionStore.js";

import CollectionConstants from "../constants/CollectionConstants.js";
import NowPlayingConstants from "../constants/NowPlayingConstants.js";


var CollectionActions = {

  fetch: function(path) {
    if (CollectionStore.getCollection(path)) {
      CollectionStore.emitChange(path);
      return;
    }
    WebsocketAPI.send(CollectionConstants.FETCH, {path: path});
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
      path: path,
    });
  },

  playNow: function(path) {
    AppDispatcher.handleViewAction({
      actionType: CollectionConstants.PLAY_NOW,
      path: path,
    });
  },

};

export default CollectionActions;
