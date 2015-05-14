'use strict';

import WebsocketAPI from '../utils/WebsocketAPI.js';

import AppDispatcher from '../dispatcher/AppDispatcher.js';

import SearchConstants from '../constants/SearchConstants.js';


var SearchActions = {

  search: function(input) {
    WebsocketAPI.send(SearchConstants.SEARCH, {input: input});

    AppDispatcher.handleViewAction({
      actionType: SearchConstants.SEARCH,
      input: input,
    });
  }

};

export default SearchActions;
