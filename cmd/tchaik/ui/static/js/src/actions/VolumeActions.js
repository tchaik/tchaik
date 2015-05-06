'use strict';

var AppDispatcher = require('../dispatcher/AppDispatcher');

var VolumeConstants = require('../constants/VolumeConstants.js');

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

module.exports = VolumeActions;
