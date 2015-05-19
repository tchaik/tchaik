'use strict';

import AppDispatcher from '../dispatcher/AppDispatcher';

import LeftColumnConstants from '../constants/LeftColumnConstants.js';


var LeftColumnActions = {

  mode: function(m) {
    AppDispatcher.handleViewAction({
      actionType: LeftColumnConstants.MODE,
      mode: m
    });
  },

  toggleVisibility: function() {
    AppDispatcher.handleViewAction({
      actionType: LeftColumnConstants.TOGGLE_VISIBILITY,
    });
  },
};

export default LeftColumnActions;
