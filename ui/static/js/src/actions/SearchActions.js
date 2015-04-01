'use strict';

var WebsocketApi = require('../api/WebsocketApi.js');

var AppDispatcher = require('../dispatcher/AppDispatcher.js');

var SearchConstants = require('../constants/SearchConstants.js');

var SearchActions = {

  search: function(input) {
    WebsocketApi.send({
      input: input,
      action: SearchConstants.SEARCH,
    });

    AppDispatcher.handleViewAction({
      actionType: SearchConstants.SEARCH,
      input: input,
    });
  },

};

module.exports = SearchActions;
