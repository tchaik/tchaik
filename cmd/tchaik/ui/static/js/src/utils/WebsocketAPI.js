'use strict';

import {Store} from '../stores/Store.js';

import AppDispatcher from '../dispatcher/AppDispatcher.js';

import WebsocketActions from '../actions/WebsocketActions.js';
import WebsocketConstants from '../constants/WebsocketConstants.js';


class WebsocketAPI extends Store {
  constructor() {
    super();

    this.host = null;
    this.open = false;
    this.queue = [];
    this.sock = null;

    this.dispatchToken = AppDispatcher.register(this._handleViewAction.bind(this));
  }

  init(host) {
    if (this.host === null) {
      this.host = host;
    }

    try {
      this.sock = new WebSocket(host);
    } catch (exception) {
      console.log("Error created websocket");
      console.log(exception);
      return;
    }

    this.sock.onmessage = this._onMessage.bind(this);
    this.sock.onerror = this._onError.bind(this);
    this.sock.onopen = this._onOpen.bind(this);
    this.sock.onclose = this._onClose.bind(this);
  }

  getStatus() {
    return {'open': this.open};
  }

  send(action, data) {
    var payload = {
      action: action,
      data: data
    };
    if (!this.open) {
      this.queue.push(payload);
      return;
    }
    this.sock.send(JSON.stringify(payload));
  }

  _handleViewAction(payload) {
    var action = payload.action;
    var source = payload.source;

    if (source === 'VIEW_ACTION') {
      switch (action.actionType) {
        case WebsocketConstants.RECONNECT:
          this.init(this.host);
          break;
      }
    }
    return true;
  }

  _onMessage(obj) {
    var msg = JSON.parse(obj.data);
    WebsocketActions.dispatch(msg);
  }

  _onError(err) {
    console.error(err);
  }

  _onOpen() {
    this.open = true;
    this.emitChange();
    this.queue.map(function(payload) {
      this.send(payload.action, payload.data);
    }.bind(this));
    this.queue = [];
  }

  _onClose() {
    this.open = false;
    this.emitChange();
  }
}

export default new WebsocketAPI();