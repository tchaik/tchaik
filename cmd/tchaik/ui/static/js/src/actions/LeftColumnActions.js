import AppDispatcher from "../dispatcher/AppDispatcher";

import LeftColumnConstants from "../constants/LeftColumnConstants.js";


var LeftColumnActions = {
  toggleVisibility: function() {
    AppDispatcher.handleViewAction({
      actionType: LeftColumnConstants.TOGGLE_VISIBILITY,
    });
  },
};

export default LeftColumnActions;
