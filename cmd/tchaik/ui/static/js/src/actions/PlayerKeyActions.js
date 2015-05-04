'use strict';

var AppDispatcher = require('../dispatcher/AppDispatcher.js');

var ControlConstants = require('../constants/ControlConstants.js');

var PlayerKeyActions = {

  setKey: function(key) {
    AppDispatcher.handleViewAction({
      actionType: ControlConstants.SET_KEY,
      key: key,
    });
  },

};

module.exports = PlayerKeyActions;
