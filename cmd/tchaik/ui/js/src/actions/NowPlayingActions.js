import AppDispatcher from "../dispatcher/AppDispatcher";

import WebsocketAPI from "../utils/WebsocketAPI.js";

import NowPlayingStore from "../stores/NowPlayingStore.js";
import NowPlayingConstants from "../constants/NowPlayingConstants.js";
import CursorConstants from "../constants/CursorConstants.js";


var NowPlayingActions = {

  reset: function() {
    AppDispatcher.handleViewAction({
      actionType: NowPlayingConstants.RESET,
    });
  },

  setError: function(err) {
    AppDispatcher.handleViewAction({
      actionType: NowPlayingConstants.SET_ERROR,
      error: err,
    });
  },

  setDuration: function(duration) {
    AppDispatcher.handleViewAction({
      actionType: NowPlayingConstants.SET_DURATION,
      duration: duration,
    });
  },

  setBuffered: function(buffered) {
    AppDispatcher.handleViewAction({
      actionType: NowPlayingConstants.SET_BUFFERED,
      buffered: buffered,
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

  ended: function(source, repeat) {
    WebsocketAPI.send(NowPlayingConstants.RECORD_PLAY, {path: ["T", NowPlayingStore.getTrack().id]});

    if (source === "cursor") {
      WebsocketAPI.send(CursorConstants.CURSOR, {
        action: CursorConstants.NEXT,
        name: "Default",
      });
    }

    AppDispatcher.handleViewAction({
      actionType: NowPlayingConstants.ENDED,
      source: source,
      repeat: repeat,
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

  repeat: function(v) {
    AppDispatcher.handleViewAction({
      actionType: NowPlayingConstants.SET_REPEAT,
      repeat: v,
    });
  },

};

export default NowPlayingActions;
