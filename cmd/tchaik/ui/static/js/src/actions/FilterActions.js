'use strict';

var AppDispatcher = require('../dispatcher/AppDispatcher');

var WebsocketAPI = require('../utils/WebsocketAPI.js');

var FilterConstants = require('../constants/FilterConstants.js');

var FilterActions = {

  fetchList: function(name) {
    WebsocketAPI.send(FilterConstants.FILTER_LIST, name);
  },

  fetchPaths: function(name, itemName) {
    WebsocketAPI.send(FilterConstants.FILTER_PATHS, {
      'name': name,
      'path': [itemName]
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
