import AppDispatcher from "../dispatcher/AppDispatcher.js";

import ContainerConstants from "../constants/ContainerConstants.js";

export default {
  mode: function(m) {
    AppDispatcher.handleViewAction({
      actionType: ContainerConstants.MODE,
      mode: m,
    });
  },
};
