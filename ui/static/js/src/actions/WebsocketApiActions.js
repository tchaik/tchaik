'use strict';

var AppDispatcher = require('../dispatcher/AppDispatcher');

var WebsocketApiActions = {

  dispatch: function(data) {
    AppDispatcher.handleServerAction({
      actionType: data.Action,
      data: data.Data,
    });
  },

};

module.exports = WebsocketApiActions;
