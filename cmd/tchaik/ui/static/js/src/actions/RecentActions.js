'use strict';

var WebsocketAPI = require('../utils/WebsocketAPI.js');

var RecentConstants = require('../constants/RecentConstants.js');

var RecentActions = {

  fetch: function() {
    WebsocketAPI.send({
      action: RecentConstants.FETCH_RECENT,
    });
  },

};

module.exports = RecentActions;
