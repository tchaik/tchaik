'use strict';

import AppDispatcher from '../dispatcher/AppDispatcher.js';

import ControlConstants from '../constants/ControlConstants.js';


var PlayerKeyActions = {

  setKey: function(key) {
    AppDispatcher.handleViewAction({
      actionType: ControlConstants.SET_KEY,
      key: key,
    });
  },

  setPushKey: function(key) {
    AppDispatcher.handleViewAction({
      actionType: ControlConstants.SET_PUSH_KEY,
      key: key,
    });
  }

};

export default PlayerKeyActions;
