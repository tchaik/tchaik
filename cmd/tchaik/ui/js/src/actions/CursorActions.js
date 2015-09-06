import AppDispatcher from "../dispatcher/AppDispatcher";

import CursorConstants from "../constants/CursorConstants.js";
import WebsocketAPI from "../utils/WebsocketAPI.js";

var cursorName = "Default";


var CursorActions = {

  fetch: function() {
    WebsocketAPI.send(CursorConstants.CURSOR, {
      action: CursorConstants.FETCH,
      name: cursorName,
    });
  },

  set: function(index, path) {
    WebsocketAPI.send(CursorConstants.CURSOR, {
      action: CursorConstants.SET,
      name: cursorName,
      index: index,
      path: path,
    });

    AppDispatcher.handleViewAction({
      actionType: CursorConstants.SET,
      index: index,
      path: path,
    });
  },

  next: function() {
    WebsocketAPI.send(CursorConstants.CURSOR, {
      action: CursorConstants.NEXT,
      name: cursorName,
    });

    AppDispatcher.handleViewAction({
      actionType: CursorConstants.NEXT,
    });
  },

  prev: function() {
    WebsocketAPI.send(CursorConstants.CURSOR, {
      action: CursorConstants.PREV,
      name: cursorName,
    });

    AppDispatcher.handleViewAction({
      actionType: CursorConstants.PREV,
    });
  },

};

export default CursorActions;
