'use strict';

var AppDispatcher = require('../dispatcher/AppDispatcher');
var EventEmitter = require('eventemitter3').EventEmitter;
var assign = require('object-assign');

var ApiKeyConstants = require('../constants/ApiKeyConstants.js');

var CHANGE_EVENT = 'change';

var _apiKey = null;

function setKey(k) {
  _apiKey = k;
  localStorage.setItem("apiKey", k);
}

function key() {
  if (_apiKey !== null) {
    return _apiKey;
  }
  var k = localStorage.getItem("apiKey");
  _apiKey = (k) ? k : null;
  return _apiKey;
}

var ApiKeyStore = assign({}, EventEmitter.prototype, {

  isKeySet: function() {
    var k = key();
    if (k === null || k === "") {
      return false;
    }
    return true;
  },

  getKey: function() {
    return key();
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

ApiKeyStore.dispatchToken = AppDispatcher.register(function(payload) {
  var action = payload.action;
  var source = payload.source;

  if (source === 'VIEW_ACTION') {
    switch (action.actionType) {
      case ApiKeyConstants.SET_KEY:
        setKey(action.key);
        ApiKeyStore.emitChange();
        break;
    }
  }
  return true;
});

module.exports = ApiKeyStore;
