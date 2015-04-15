'use strict';

var AppDispatcher = require('../dispatcher/AppDispatcher.js');

var APIKeyConstants = require('../constants/APIKeyConstants.js');

var APIKeyActions = {

  setKey: function(key) {
    AppDispatcher.handleViewAction({
      actionType: APIKeyConstants.SET_KEY,
      key: key,
    });
  },

};

module.exports = APIKeyActions;
