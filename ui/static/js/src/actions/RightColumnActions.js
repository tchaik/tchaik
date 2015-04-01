'use strict';

var AppDispatcher = require('../dispatcher/AppDispatcher');

var RightColumnConstants = require('../constants/RightColumnConstants.js');

var RightColumnActions = {

  layout: function() {
    AppDispatcher.handleViewAction({
      actionType: RightColumnConstants.LAYOUT,
    });
  },

};

module.exports = RightColumnActions;
