import AppDispatcher from "../dispatcher/AppDispatcher";

import VolumeConstants from "../constants/VolumeConstants.js";


var VolumeActions = {

  volume: function(v) {
    AppDispatcher.handleViewAction({
      actionType: VolumeConstants.SET_VOLUME,
      volume: v,
    });
  },

  toggleVolumeMute: function() {
    AppDispatcher.handleViewAction({
      actionType: VolumeConstants.TOGGLE_VOLUME_MUTE,
    });
  },

};

export default VolumeActions;
