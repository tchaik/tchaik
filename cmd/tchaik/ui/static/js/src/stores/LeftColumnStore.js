'use strict';

var AppDispatcher = require('../dispatcher/AppDispatcher');
var EventEmitter = require('eventemitter3').EventEmitter;
var assign = require('object-assign');

var LeftColumnConstants = require('../constants/LeftColumnConstants.js');

var CHANGE_EVENT = 'change';

var _defaultMode = "All";

function setMode(m) {
  if (m === null) {
    m = _defaultMode;
  }
  localStorage.setItem("mode", m);
}

function mode() {
  var m = localStorage.getItem("mode");
  if (m === null) {
    m = _defaultMode;
  }
  return m;
}

var LeftColumnStore = assign({}, EventEmitter.prototype, {

  getMode: function() {
    return mode();
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

LeftColumnStore.dispatchToken = AppDispatcher.register(function(payload) {
  var action = payload.action;
  var source = payload.source;

  if (source === 'VIEW_ACTION') {
    switch (action.actionType) {
      case LeftColumnConstants.MODE:
        setMode(action.mode);
        LeftColumnStore.emitChange();
        break;
    }
  }
  return true;
});

module.exports = LeftColumnStore;
