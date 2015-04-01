'use strict';

var AppDispatcher = require('../dispatcher/AppDispatcher');

var LeftColumnConstants = require('../constants/LeftColumnConstants.js');

var LeftColumnActions = {

  mode: function(m) {
    AppDispatcher.handleViewAction({
      actionType: LeftColumnConstants.MODE,
      mode: m
    });
  },

};

module.exports = LeftColumnActions;
