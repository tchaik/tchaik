'use strict';

var AppDispatcher = require('../dispatcher/AppDispatcher.js');

var ApiKeyConstants = require('../constants/ApiKeyConstants.js');

var ApiKeyActions = {

  setKey: function(key) {
    AppDispatcher.handleViewAction({
      actionType: ApiKeyConstants.SET_KEY,
      key: key,
    });
  },

};

module.exports = ApiKeyActions;
