import AppDispatcher from "../dispatcher/AppDispatcher";

import WebsocketConstants from "../constants/WebsocketConstants.js";


var WebsocketActions = {

  dispatch: function(data) {
    AppDispatcher.handleServerAction({
      actionType: data.Action,
      data: data.Data,
    });
  },

  reconnect: function() {
    AppDispatcher.handleViewAction({
      actionType: WebsocketConstants.RECONNECT,
    });
  },

};

export default WebsocketActions;
