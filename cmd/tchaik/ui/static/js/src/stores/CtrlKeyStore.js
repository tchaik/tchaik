'use strict';

var AppDispatcher = require('../dispatcher/AppDispatcher');
var EventEmitter = require('eventemitter3').EventEmitter;
var assign = require('object-assign');

var WebsocketConstants = require('../constants/WebsocketConstants.js');
var WebsocketAPI = require('../utils/WebsocketAPI.js');

var ControlConstants = require('../constants/ControlConstants.js');

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

function sendKey(key) {
  WebsocketAPI.send({
    action: "KEY",
    data: key,
  });
}

var CtrlKeyStore = assign({}, EventEmitter.prototype, {

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

CtrlKeyStore.dispatchToken = AppDispatcher.register(function(payload) {
  var action = payload.action;
  var source = payload.source;

  if (source === 'VIEW_ACTION') {
    switch (action.actionType) {
      case ControlConstants.SET_KEY:
        setKey(action.key);
        sendKey(action.key);
        CtrlKeyStore.emitChange();
        break;

      case WebsocketConstants.RECONNECT:
        if (CtrlKeyStore.isKeySet()) {
          sendKey(CtrlKeyStore.getKey());
        }
        break;
    }
  }
  return true;
});

module.exports = CtrlKeyStore;
