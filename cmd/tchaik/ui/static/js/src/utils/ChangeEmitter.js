"use strict";

import EventEmitter from "eventemitter3";

export var CHANGE_EVENT = "change";

export class ChangeEmitter extends EventEmitter {
  emitChange() {
    this.emit(CHANGE_EVENT);
  }

  /**
   * @param {function} callback
   */
  addChangeListener(callback) {
    this.on(CHANGE_EVENT, callback);
  }

  /**
   * @param {function} callback
   */
  removeChangeListener(callback) {
    this.removeListener(CHANGE_EVENT, callback);
  }
}
