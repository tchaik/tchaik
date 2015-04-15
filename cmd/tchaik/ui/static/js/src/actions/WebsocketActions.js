'use strict';

var AppDispatcher = require('../dispatcher/AppDispatcher');

var WebsocketConstants = require('../constants/WebsocketConstants.js');

var WebsocketActions = {

  dispatch: function(data) {
    AppDispatcher.handleServerAction({
      actionType: data.Action,
      data: data.Data,
    });
  },

  reconnect: function() {
    AppDispatcher.handleViewAction({
      actionType: WebsocketConstants.RECONNECT
    });
  },

};

module.exports = WebsocketActions;
