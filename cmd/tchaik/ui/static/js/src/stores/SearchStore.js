'use strict';

var AppDispatcher = require('../dispatcher/AppDispatcher');
var EventEmitter = require('eventemitter3').EventEmitter;
var assign = require('object-assign');

var SearchConstants = require('../constants/SearchConstants.js');

var CHANGE_EVENT = 'change';

var _results = [];

function setResults(results) {
  if (results === null) {
    results = [];
  }
  _results = results;
}

function setInput(input) {
  localStorage.setItem("searchInput", input);
}

function input() {
  var s = localStorage.getItem("searchInput");
  if (s === null) {
    return "";
  }
  return s;
}

var SearchResultStore = assign({}, EventEmitter.prototype, {

  getResults: function() {
    return _results;
  },

  getInput: function() {
    return input();
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

SearchResultStore.dispatchToken = AppDispatcher.register(function(payload) {
  var action = payload.action;
  var source = payload.source;

  if (source === 'SERVER_ACTION') {
    switch (action.actionType) {
      case SearchConstants.SEARCH:
        setResults(action.data);
        SearchResultStore.emitChange();
        break;
    }
  }

  if (source === 'VIEW_ACTION') {
    switch (action.actionType) {
      case SearchConstants.SEARCH:
        setInput(action.input);
        SearchResultStore.emitChange();
        break;
    }
  }

  return true;
});

module.exports = SearchResultStore;
