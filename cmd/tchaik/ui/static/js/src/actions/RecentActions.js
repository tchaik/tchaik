'use strict';

import WebsocketAPI from '../utils/WebsocketAPI.js';

import RecentConstants from '../constants/RecentConstants.js';


var RecentActions = {

  fetch: function() {
    WebsocketAPI.send(RecentConstants.FETCH_RECENT);
  }

};

export default RecentActions;
