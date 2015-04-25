'use strict';

var AppDispatcher = require('../dispatcher/AppDispatcher');

var WebsocketAPI = require('../utils/WebsocketAPI.js');

var FilterConstants = require('../constants/FilterConstants.js');

var FilterActions = {

  fetchList: function(name) {
    WebsocketAPI.send({
      data: name,
      action: FilterConstants.FILTER_LIST,
    });
  },

  fetchPaths: function(name, itemName) {
    WebsocketAPI.send({
      data: name,
      path: [itemName],
      action: FilterConstants.FILTER_PATHS,
    });
  },

  setItem: function(name, itemName) {
    AppDispatcher.handleViewAction({
      actionType: FilterConstants.SET_FILTER_ITEM,
      name: name,
      itemName: itemName,
    });
  }

};

module.exports = FilterActions;
