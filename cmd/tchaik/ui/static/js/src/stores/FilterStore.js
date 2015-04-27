'use strict';

var AppDispatcher = require('../dispatcher/AppDispatcher');
var EventEmitter = require('eventemitter3').EventEmitter;
var assign = require('object-assign');

var FilterConstants = require('../constants/FilterConstants.js');

var CHANGE_EVENT = 'change';

var _currentItems = null;
var _filters = {};
var _filterPaths = {};

function _getCurrentItem(name) {
  if (_currentItems === null) {
    var fi = localStorage.getItem("filterCurrentItems");
    if (fi) {
      _currentItems = JSON.parse(fi);
    }
  }
  if (!_currentItems || !_currentItems[name]) {
    return null;
  }
  return _currentItems[name];
}

function _setCurrentItem(name, itemName) {
  if (_currentItems === null) {
    _currentItems = {};
  }
  _currentItems[name] = itemName;
  localStorage.setItem("filterCurrentItems", JSON.stringify(_currentItems));
}

function _addFilterPaths(name, itemName, paths) {
  var filterPaths = _filterPaths[name];
  if (!filterPaths) {
    filterPaths = {};
  }
  filterPaths[itemName] = paths;
  _filterPaths[name] = filterPaths;
}

var FilterStore = assign({}, EventEmitter.prototype, {

  getCurrentItem: function(name) {
    return _getCurrentItem(name);
  },

  getItems: function(name) {
    var x = _filters[name];
    if (!x) {
      x = [];
    }
    return x;
  },
  
  getPaths: function(name, itemName) {
    var x = _filterPaths[name];
    if (!x) {
      return [];
    }
    var paths = x[itemName];
    if (!paths) {
      paths = [];
    }
    return paths;
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

FilterStore.dispatchToken = AppDispatcher.register(function(payload) {
  var action = payload.action;
  var source = payload.source;

  if (source === 'SERVER_ACTION') {
    switch (action.actionType) {
      case FilterConstants.FILTER_LIST:
        _filters[action.data.Name] = action.data.Items;
        FilterStore.emitChange();
        break;
      case FilterConstants.FILTER_PATHS:
        var path = action.data.Path; // [name, item]
        _addFilterPaths(path[0], path[1], action.data.Paths);
        FilterStore.emitChange();
        break;
    }
  } else if (source === 'VIEW_ACTION') {
    switch (action.actionType) {
      case FilterConstants.SET_FILTER_ITEM:
        _setCurrentItem(action.name, action.itemName);
        FilterStore.emitChange();
        break;
    }
  }

  return true;
});

module.exports = FilterStore;
