'use strict';

var WebsocketApi = require('../api/WebsocketApi.js');

var AppDispatcher = require('../dispatcher/AppDispatcher.js');

var WebsocketApi = require('../api/WebsocketApi.js');

var ApiKeyConstants = require('../constants/ApiKeyConstants.js');

var ApiKeyActions = {

  setKey: function(key) {
    WebsocketApi.send({
      action: "KEY",
      data: key,
    });

    AppDispatcher.handleViewAction({
      actionType: ApiKeyConstants.SET_KEY,
      key: key,
    });
  },

};

module.exports = ApiKeyActions;
