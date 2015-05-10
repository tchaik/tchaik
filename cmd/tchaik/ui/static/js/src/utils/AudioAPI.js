'use strict';

var NowPlayingConstants = require('../constants/NowPlayingConstants.js');
var NowPlayingActions = require('../actions/NowPlayingActions.js');
var NowPlayingStore = require('../stores/NowPlayingStore.js');

var PlayingStatusStore = require('../stores/PlayingStatusStore.js');

var VolumeStore = require('../stores/VolumeStore.js');

var _audio = new Audio();
var _src = null;
var _playing = false;

function buffered() {
  return _audio.buffered;
}

function play() {
  _audio.play();
}

function pause() {
  _audio.pause();
}

function load() {
  _audio.load();
}

function src() {
  return _src;
}

function setSrc(l) {
  _audio.src = l;
  _src = l;
}

function setCurrentTime(t) {
  _audio.currentTime = t;
}

function currentTime() {
  return _audio.currentTime;
}

function setVolume(t) {
  _audio.volume = t;
}

function volume() {
  return _audio.volume;
}

function duration() {
  return _audio.duration;
}

function init() {
  _audio.addEventListener('error', onPlayerEvent);
  _audio.addEventListener('progress', onPlayerEvent);
  _audio.addEventListener('play', onPlayerEvent);
  _audio.addEventListener('pause', onPlayerEvent);
  _audio.addEventListener('ended', onPlayerEvent);
  _audio.addEventListener('timeupdate', onPlayerEvent);
  _audio.addEventListener('loadedmetadata', onPlayerEvent);
  _audio.addEventListener('loadstart', onPlayerEvent);

  NowPlayingStore.addChangeListener(update);
  NowPlayingStore.addControlListener(_onNowPlayingControl);
  VolumeStore.addChangeListener(_onVolumeChange);

  update();
  _onVolumeChange();
}

function onPlayerEvent(evt) {
  switch (evt.type) {
  case "error":
    NowPlayingActions.setError(evt.srcElement.error);
    console.log("Error received from Audio component:");
    console.error(evt);
    break;

  case "progress":
    var range = buffered();
    if (range.length > 0) {
      NowPlayingActions.setBuffered(range.end(range.length-1));
    }
    break;

  case "play":
    NowPlayingActions.playing(true);
    break;

  case "pause":
    NowPlayingActions.playing(false);
    break;

  case "ended":
    NowPlayingActions.ended(NowPlayingStore.getSource());
    break;

  case "timeupdate":
    NowPlayingActions.currentTime(currentTime());
    break;

  case "loadedmetadata":
    NowPlayingActions.setDuration(duration());

    setCurrentTime(PlayingStatusStore.getTime());
    if (_playing) {
      play();
    }
    break;

  case "loadstart":
    NowPlayingActions.reset();
    break;

  default:
    console.warn("unhandled player event:");
    console.warn(evt);
    break;
  }
}

function _onNowPlayingControl(type, value) {
  if (type === NowPlayingConstants.SET_CURRENT_TIME) {
    setCurrentTime(value);
  }
}

function update() {
  var track = NowPlayingStore.getTrack();
  if (track) {
    var source = "/track/"+track.TrackID;
    var orig = src();
    if (orig !== source) {
      setSrc(source);
      load();
      play();
    }
  }

  var prevPlaying = _playing;
  _playing = NowPlayingStore.getPlaying();
  if (prevPlaying !== _playing) {
    if (_playing) {
      play();
    } else {
      pause();
    }
  }
}

function _onVolumeChange() {
  var v = VolumeStore.getVolume();
  if (volume() !== v) {
    setVolume(v);
  }
}


var AudioAPI = {

  init: function() {
    init();
  },

};

module.exports = AudioAPI;