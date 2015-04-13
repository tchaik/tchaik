'use strict';

var keyMirror = require('keymirror');

module.exports = keyMirror({
  STORE_CURRENT_TIME: null,
  SET_CURRENT_TIME: null,
  SET_CURRENT_TRACK: null,
  SET_VOLUME: null,
  TOGGLE_VOLUME_MUTE: null,
  SET_PLAYING: null,
  ENDED: null,
});
