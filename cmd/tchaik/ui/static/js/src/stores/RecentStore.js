'use strict';

import AppDispatcher from '../dispatcher/AppDispatcher';

import {Store} from './Store.js';

import RecentConstants from '../constants/RecentConstants.js';


var _recentPaths = [];

class RecentStore extends Store {
  getPaths() {
    return _recentPaths;
  }
}

var _recentStore = new RecentStore();

_recentStore.dispatchToken = AppDispatcher.register(function(payload) {
  var action = payload.action;
  var source = payload.source;

  if (source === 'SERVER_ACTION') {
    switch (action.actionType) {
      case RecentConstants.FETCH_RECENT:
        _recentPaths = action.data;
        _recentStore.emitChange();
        break;
    }
  }

  return true;
});

export default _recentStore;
