"use strict";

import React from "react";

import NowPlayingConstants from "../constants/NowPlayingConstants.js";
import NowPlayingActions from "../actions/NowPlayingActions.js";
import NowPlayingStore from "../stores/NowPlayingStore.js";

import PlayingStatusStore from "../stores/PlayingStatusStore.js";

import VolumeStore from "../stores/VolumeStore.js";

function getPlayerState() {
  return {
    source: `/track/${NowPlayingStore.getTrack().id}`,
    playing: NowPlayingStore.getPlaying(),
    volume: VolumeStore.getVolume(),
  };
}

export default class Player extends React.Component {
  constructor(props) {
    super(props);

    this.state = getPlayerState();

    this._onPlayerEvent = this._onPlayerEvent.bind(this);
    this._onChange = this._onChange.bind(this);
    this._onNowPlayingControl = this._onNowPlayingControl.bind(this);
  }

  componentDidMount() {
    this._audio = new Audio();
    this._audio.addEventListener("error", this._onPlayerEvent);
    this._audio.addEventListener("progress", this._onPlayerEvent);
    this._audio.addEventListener("play", this._onPlayerEvent);
    this._audio.addEventListener("pause", this._onPlayerEvent);
    this._audio.addEventListener("ended", this._onPlayerEvent);
    this._audio.addEventListener("timeupdate", this._onPlayerEvent);
    this._audio.addEventListener("loadedmetadata", this._onPlayerEvent);
    this._audio.addEventListener("loadstart", this._onPlayerEvent);

    NowPlayingStore.addChangeListener(this._onChange);
    NowPlayingStore.addControlListener(this._onNowPlayingControl);
  }

  componentWillUnmount() {
    this._audio.removeEventListener("error", this._onPlayerEvent);
    this._audio.removeEventListener("progress", this._onPlayerEvent);
    this._audio.removeEventListener("play", this._onPlayerEvent);
    this._audio.removeEventListener("pause", this._onPlayerEvent);
    this._audio.removeEventListener("ended", this._onPlayerEvent);
    this._audio.removeEventListener("timeupdate", this._onPlayerEvent);
    this._audio.removeEventListener("loadedmetadata", this._onPlayerEvent);
    this._audio.removeEventListener("loadstart", this._onPlayerEvent);

    NowPlayingStore.removeChangeListener(this._onChange);
    NowPlayingStore.removeControlListener(this._onNowPlayingControl);
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

  setSrc(src) {
    this._audio.src = src;
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

  _onPlayerEvent(evt) {
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

  componentDidUpdate(prevProps, prevState) {
    if (prevState.playing !== this.state.playing) {
      if (this.state.playing) {
        this.play();
      } else {
        this.pause();
      }
    }

    if (prevState.source !== this.state.source) {
      this.setSrc(this.state.source);
      this.load();
      this.play();
    }
  }

  render() {
    return null;
  }

  _onChange() {
    this.setState(getPlayerState());
  }

  _onNowPlayingControl(type, value) {
    if (type === NowPlayingConstants.SET_CURRENT_TIME) {
      this.setCurrentTime(value);
    }
  }
}
