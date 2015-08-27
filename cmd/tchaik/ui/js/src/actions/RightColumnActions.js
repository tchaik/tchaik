"use strict";

import AppDispatcher from "../dispatcher/AppDispatcher";

import RightColumnConstants from "../constants/RightColumnConstants.js";


var RightColumnActions = {
  toggle: function() {
    AppDispatcher.handleViewAction({
      actionType: RightColumnConstants.TOGGLE_RIGHTCOLUMN,
    });
  },
};

export default RightColumnActions;
