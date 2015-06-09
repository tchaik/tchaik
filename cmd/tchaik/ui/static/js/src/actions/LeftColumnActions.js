import AppDispatcher from "../dispatcher/AppDispatcher";

import LeftColumnConstants from "../constants/LeftColumnConstants.js";


var LeftColumnActions = {
  toggle: function() {
    AppDispatcher.handleViewAction({
      actionType: LeftColumnConstants.TOGGLE_LEFTCOLUMN,
    });
  },
};

export default LeftColumnActions;
