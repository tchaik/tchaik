"use strict";

import {ChangeEmitter} from "../utils/ChangeEmitter.js";
import AppDispatcher from "../dispatcher/AppDispatcher";

import WebsocketConstants from "../constants/WebsocketConstants.js";
import WebsocketAPI from "../utils/WebsocketAPI.js";

import ControlConstants from "../constants/ControlConstants.js";


var _playerKey = null;
var _pushKey = null;

function setKey(k) {
  _playerKey = k;
  localStorage.setItem("playerKey", k);
}

function key() {
  if (_playerKey !== null) {
    return _playerKey;
  }
  var k = localStorage.getItem("playerKey");
  _playerKey = (k) ? k : "";
  return _playerKey;
}

function sendKey(k) {
  WebsocketAPI.send("KEY", {key: k});
}

function setPushKey(k) {
  _pushKey = k;
  localStorage.setItem("pushKey", k);
}

function pushKey() {
  if (_pushKey !== null) {
    return _pushKey;
  }
  var k = localStorage.getItem("pushKey");
  _pushKey = (k) ? k : "";
  return _pushKey;
}


class PlayerKeyStore extends ChangeEmitter {
  isKeySet() {
    var k = key();
    if (k === null || k === "") {
      return false;
    }
    return true;
  }

  getKey() {
    return key();
  }

  isPushKeySet() {
    var k = pushKey();
    if (k === null || k === "") {
      return false;
    }
    return true;
  }

  getPushKey() {
    return pushKey();
  }
}

var _playerKeyStore = new PlayerKeyStore();

_playerKeyStore.dispatchToken = AppDispatcher.register(function(payload) {
  var action = payload.action;
  var source = payload.source;

  if (source === "VIEW_ACTION") {
    switch (action.actionType) {
      case ControlConstants.SET_KEY:
        setKey(action.key);
        sendKey(action.key);
        _playerKeyStore.emitChange();
        break;

      case ControlConstants.SET_PUSH_KEY:
        setPushKey(action.key);
        _playerKeyStore.emitChange();
        break;

      case WebsocketConstants.RECONNECT:
        if (_playerKeyStore.isKeySet()) {
          sendKey(_playerKeyStore.getKey());
        }
        break;
    }
  }
  return true;
});

export default _playerKeyStore;
