import AppDispatcher from "../dispatcher/AppDispatcher";

import WebsocketAPI from "../utils/WebsocketAPI.js";

import FilterConstants from "../constants/FilterConstants.js";


var FilterActions = {

  fetchList: function(name) {
    WebsocketAPI.send(FilterConstants.FILTER_LIST, {name: name});
  },

  fetchPaths: function(name, itemName) {
    WebsocketAPI.send(FilterConstants.FILTER_PATHS, {
      name: name,
      path: [itemName],
    });
  },

  setItem: function(name, itemName) {
    AppDispatcher.handleViewAction({
      actionType: FilterConstants.SET_FILTER_ITEM,
      name: name,
      itemName: itemName,
    });
  },

};

export default FilterActions;
