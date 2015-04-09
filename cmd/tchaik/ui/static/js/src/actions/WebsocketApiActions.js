'use strict';

var AppDispatcher = require('../dispatcher/AppDispatcher');

var WebsocketApiConstants = require('../constants/WebsocketApiConstants.js');

var WebsocketApiActions = {

  dispatch: function(data) {
    AppDispatcher.handleServerAction({
      actionType: data.Action,
      data: data.Data,
    });
  },

  reconnect: function() {
    AppDispatcher.handleViewAction({
      actionType: WebsocketApiConstants.RECONNECT
    });
  },

};

module.exports = WebsocketApiActions;
