'use strict';

var AppDispatcher = require('../dispatcher/AppDispatcher');
var EventEmitter = require('eventemitter3').EventEmitter;
var assign = require('object-assign');

var RightColumnConstants = require('../constants/RightColumnConstants.js');

var CHANGE_EVENT = 'change';


var RightColumnStore = assign({}, EventEmitter.prototype, {

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

RightColumnStore.dispatchToken = AppDispatcher.register(function(payload) {
  var action = payload.action;
  var source = payload.source;

  if (source === 'VIEW_ACTION') {
    switch (action.actionType) {
      case RightColumnConstants.LAYOUT:
        RightColumnStore.emitChange();
        break;
    }
  }
  return true;
});

module.exports = RightColumnStore;
