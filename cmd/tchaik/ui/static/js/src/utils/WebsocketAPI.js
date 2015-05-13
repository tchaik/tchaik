'use strict';

import EventEmitter from 'eventemitter3';

import {Store} from '../stores/Store.js';

import AppDispatcher from '../dispatcher/AppDispatcher.js';

import WebsocketActions from '../actions/WebsocketActions.js';
import WebsocketConstants from '../constants/WebsocketConstants.js';

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
  WebsocketActions.dispatch(msg);
}

function onError(err) {
  console.error(err);
}

class WebsocketAPI extends Store {
  init(host) {
    init(host);
  }

  getStatus() {
    return {
      'open': _websocket.open
    };
  }

  send(action, data) {
    var payload = {
      action: action,
      data: data
    };
    if (!_websocket.open) {
      _websocket.queue.push(payload);
      return;
    }
    _websocket.sock.send(JSON.stringify(payload));
  }
}

var _websocketAPI = new WebsocketAPI();

function onOpen() {
  _websocket.open = true;
  _websocketAPI.emitChange();
  _websocket.queue.map(function(payload) {
    _websocketAPI.send(payload.action, payload.data);
  });
  _websocket.queue = [];
}


function onClose() {
  _websocket.open = false;
  _websocketAPI.emitChange();
}

_websocketAPI.dispatchToken = AppDispatcher.register(function(payload) {
  var action = payload.action;
  var source = payload.source;

  if (source === 'VIEW_ACTION') {
    switch (action.actionType) {
      case WebsocketConstants.RECONNECT:
        init(_host);
        break;
    }
  }

  return true;
});

export default _websocketAPI;