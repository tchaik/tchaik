'use strict';

var AppDispatcher = require('../dispatcher/AppDispatcher');
var EventEmitter = require('eventemitter3').EventEmitter;
var assign = require('object-assign');

var CollectionConstants = require('../constants/CollectionConstants.js');

var CHANGE_EVENT = 'change';

var _commonFields = ["Album", "AlbumArtist", "Artist", "Composer", "Year"];

var _collections = {};

function addItem(path, item) {
  if (item.Tracks) { // fill in common fields if they are set
    item.Tracks.forEach(function(track) {
      _commonFields.forEach(function(fld) {
        if (item[fld]) {
          track[fld] = item[fld];
        }
      });
    });
  }
  _collections[CollectionStore.pathToKey(path)] = item;
}

var _expandedPaths = null;

function _getExpandedPaths() {
  var p = localStorage.getItem("expandedPaths");
  if (p === null) {
    return {};
  }
  return JSON.parse(p);
}

function getExpandedPaths() {
  if (_expandedPaths === null) {
    _expandedPaths = _getExpandedPaths();
  }
  return _expandedPaths;
}

function expandPath(path, expand) {
  if (expand) {
    _expandedPaths[path] = true;
  } else {
    delete _expandedPaths[path];
  }
  localStorage.setItem("expandedPaths", JSON.stringify(_expandedPaths));
}

var CollectionStore = assign({}, EventEmitter.prototype, {

  // pathToKey returns a string representation of the path.  The only requirement is that
  // subpaths should be prefixes.
  pathToKey: function(path) {
    if (path) {
      return path.join(">>");
    }
    return null;
  },

  getCollection: function(path) {
    var key = CollectionStore.pathToKey(path);
    return _collections[key];
  },

  isExpanded: function(path) {
    var key = CollectionStore.pathToKey(path);
    var ep = getExpandedPaths();
    if (ep[key]) {
      return true;
    }
    return false;
  },

  emitChange: function(path) {
    this.emit(CHANGE_EVENT, CollectionStore.pathToKey(path));
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

CollectionStore.dispatchToken = AppDispatcher.register(function(payload) {
  var action = payload.action;
  var source = payload.source;

  if (source === 'SERVER_ACTION') {
    switch (action.actionType) {
      case CollectionConstants.FETCH:
        addItem(action.data.Path, action.data.Item);
        CollectionStore.emitChange(action.data.Path);
        break;
    }
  } else if (source === 'VIEW_ACTION') {
    switch (action.actionType) {
      case CollectionConstants.EXPAND_PATH:
        expandPath(CollectionStore.pathToKey(action.path), action.expand);
        CollectionStore.emitChange(action.path);
        break;
    }
  }

  return true;
});

module.exports = CollectionStore;
