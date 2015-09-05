"use strict";

var NowPlayingConstants = require("../constants/NowPlayingConstants.js");
var NowPlayingActions = require("../actions/NowPlayingActions.js");
var NowPlayingStore = require("../stores/NowPlayingStore.js");

var PlayingStatusStore = require("../stores/PlayingStatusStore.js");

var VolumeStore = require("../stores/VolumeStore.js");


class AudioAPI {
  constructor() {
    this._audio = new Audio();
    this._src = null;
    this._playing = false;

    this.onPlayerEvent = this.onPlayerEvent.bind(this);

    this._audio.addEventListener("error", this.onPlayerEvent);
    this._audio.addEventListener("progress", this.onPlayerEvent);
    this._audio.addEventListener("play", this.onPlayerEvent);
    this._audio.addEventListener("pause", this.onPlayerEvent);
    this._audio.addEventListener("ended", this.onPlayerEvent);
    this._audio.addEventListener("timeupdate", this.onPlayerEvent);
    this._audio.addEventListener("loadedmetadata", this.onPlayerEvent);
    this._audio.addEventListener("loadstart", this.onPlayerEvent);

    this.update = this.update.bind(this);
    this._onNowPlayingControl = this._onNowPlayingControl.bind(this);
    this._onVolumeChange = this._onVolumeChange.bind(this);

    NowPlayingStore.addChangeListener(this.update);
    NowPlayingStore.addControlListener(this._onNowPlayingControl);
    VolumeStore.addChangeListener(this._onVolumeChange);
  }

  init() {
    this.update();
    this._onVolumeChange();
  }

  buffered() {
    return this._audio.buffered;
  }

  play() {
    return this._audio.play();
  }

  pause() {
    return this._audio.pause();
  }

  load() {
    return this._audio.load();
  }

  src() {
    return this._src;
  }

  setSrc(src) {
    this._audio.src = src;
    this._src = src;
  }

  setCurrentTime(t) {
    this._audio.currentTime = t;
  }

  currentTime() {
    return this._audio.currentTime;
  }

  setVolume(v) {
    this._audio.volume = v;
  }

  volume() {
    return this._audio.volume;
  }

  duration() {
    return this._audio.duration;
  }

  onPlayerEvent(evt) {
    switch (evt.type) {
    case "error":
      NowPlayingActions.setError(evt.srcElement.error);
      console.log("Error received from Audio component:");
      console.error(evt);
      break;

    case "progress":
      var range = this.buffered();
      if (range.length > 0) {
        NowPlayingActions.setBuffered(range.end(range.length - 1));
      }
      break;

    case "play":
      NowPlayingActions.playing(true);
      break;

    case "pause":
      NowPlayingActions.playing(false);
      break;

    case "ended":
      NowPlayingActions.ended(NowPlayingStore.getSource(), NowPlayingStore.getRepeat());
      break;

    case "timeupdate":
      NowPlayingActions.currentTime(this.currentTime());
      break;

    case "loadedmetadata":
      NowPlayingActions.setDuration(this.duration());

      this.setCurrentTime(PlayingStatusStore.getTime());
      if (this._playing) {
        this.play();
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

  _onNowPlayingControl(type, value) {
    if (type === NowPlayingConstants.SET_CURRENT_TIME) {
      this.setCurrentTime(value);
    }
  }

  update() {
    var track = NowPlayingStore.getTrack();
    if (track) {
      var source = `/track/${track.ID}`;
      var orig = this.src();
      if (orig !== source) {
        this.setSrc(source);
        this.load();
        this.play();
      }
    }

    var prevPlaying = this._playing;
    this._playing = NowPlayingStore.getPlaying();
    if (prevPlaying !== this._playing) {
      if (this._playing) {
        this.play();
      } else {
        this.pause();
      }
    }
  }

  _onVolumeChange() {
    var v = VolumeStore.getVolume();
    if (this.volume() !== v) {
      this.setVolume(v);
    }
  }
}

export default new AudioAPI();
