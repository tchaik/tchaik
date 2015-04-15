'use strict';

var AppDispatcher = require('../dispatcher/AppDispatcher');

var WebsocketAPIConstants = require('../constants/WebsocketAPIConstants.js');

var WebsocketAPIActions = {

  dispatch: function(data) {
    AppDispatcher.handleServerAction({
      actionType: data.Action,
      data: data.Data,
    });
  },

  reconnect: function() {
    AppDispatcher.handleViewAction({
      actionType: WebsocketAPIConstants.RECONNECT
    });
  },

};

module.exports = WebsocketAPIActions;
