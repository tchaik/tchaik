'use strict';

var EventEmitter = require('eventemitter3').EventEmitter;
var assign = require('object-assign');

var AppDispatcher = require('../dispatcher/AppDispatcher.js');

var WebsocketAPIActions = require('../actions/WebsocketAPIActions.js');
var WebsocketAPIConstants = require('../constants/WebsocketAPIConstants.js');

var CHANGE_EVENT = 'status';

var _host = null;

var _websocket = {
  open: false,
  queue: [],
  sock: null
};

function init(host) {
  if (_host === null) {
    _host = host;
  }

  try {
    _websocket.sock = new WebSocket(host);
  } catch (exception) {
    console.log("Error created websocket");
    console.log(exception);
    return;
  }

  _websocket.sock.onmessage = onMessage;
  _websocket.sock.onerror = onError;
  _websocket.sock.onopen = onOpen;
  _websocket.sock.onclose = onClose;
}

function onMessage(obj) {
  var msg = JSON.parse(obj.data);
  WebsocketAPIActions.dispatch(msg);
}

function onError(err) {
  console.error(err);
}

function onOpen() {
  _websocket.open = true;
  WebsocketAPI.emitChange();
  _websocket.queue.map(WebsocketAPI.send);
  _websocket.queue = [];
}

function onClose() {
  _websocket.open = false;
  WebsocketAPI.emitChange();
}

var WebsocketAPI = assign({}, EventEmitter.prototype, {

  init: function(host) {
    init(host);
  },

  getStatus: function() {
    return {
      'open': _websocket.open
    };
  },

  emitChange: function() {
    this.emit(CHANGE_EVENT);
  },

  send: function(action) {
    if (!_websocket.open) {
      _websocket.queue.push(action);
      return;
    }
    _websocket.sock.send(JSON.stringify(action));
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
  }

});

WebsocketAPI.dispatchToken = AppDispatcher.register(function(payload) {
  var action = payload.action;
  var source = payload.source;

  if (source === 'VIEW_ACTION') {
    switch (action.actionType) {
      case WebsocketAPIConstants.RECONNECT:
        init(_host);
        break;
    }
  }

  return true;
});

module.exports = WebsocketAPI;
