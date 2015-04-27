'use strict';

var AppDispatcher = require('../dispatcher/AppDispatcher');
var EventEmitter = require('eventemitter3').EventEmitter;
var assign = require('object-assign');

var RecentConstants = require('../constants/RecentConstants.js');

var CHANGE_EVENT = 'change';

var _recentPaths = [];

var RecentStore = assign({}, EventEmitter.prototype, {

  getPaths: function() {
    return _recentPaths;
  },

  emitChange: function() {
    this.emit(CHANGE_EVENT);
  },

  /**
   * @param {function} callback
   */
  addChangeListener: function(callback) {
    this.on(CHANGE_EVENT, callback);
  },

  /**
   * @param {function} callback
   */
  removeChangeListener: function(callback) {
    this.removeListener(CHANGE_EVENT, callback);
  },

});

RecentStore.dispatchToken = AppDispatcher.register(function(payload) {
  var action = payload.action;
  var source = payload.source;

  if (source === 'SERVER_ACTION') {
    switch (action.actionType) {
      case RecentConstants.FETCH_RECENT:
        console.log(action);
        _recentPaths = action.data;
        RecentStore.emitChange();
        break;
    }
  }

  return true;
});

module.exports = RecentStore;
