"use strict";

import React from "react";

import NowPlayingConstants from "../constants/NowPlayingConstants.js";
import NowPlayingActions from "../actions/NowPlayingActions.js";
import NowPlayingStore from "../stores/NowPlayingStore.js";

import PlayingStatusStore from "../stores/PlayingStatusStore.js";

import VolumeStore from "../stores/VolumeStore.js";


const audioEvents = ["error", "progress", "play", "pause", "ended", "timeupdate", "loadedmetadata", "loadstart"];

class AudioPlayer extends React.Component {
  constructor(props) {
    super(props);

    this._onPlayerEvent = this._onPlayerEvent.bind(this);
    this._onNowPlayingControl = this._onNowPlayingControl.bind(this);
  }

  componentDidMount() {
    this._audio = new Audio();

    for (let e of audioEvents) {
      this._audio.addEventListener(e, this._onPlayerEvent);
    }
    NowPlayingStore.addControlListener(this._onNowPlayingControl);
  }

  componentWillUnmount() {
    for (let e of audioEvents) {
      this._audio.removeEventListener(e, this._onPlayerEvent);
    }
    NowPlayingStore.removeControlListener(this._onNowPlayingControl);
  }

  componentDidUpdate(prevProps) {
    if (prevProps.playing !== this.props.playing) {
      if (this.props.playing) {
        this.play();
      } else {
        this.pause();
      }
    }

    if (prevProps.source !== this.props.source) {
      this.setSrc(this.props.source);
      this.load();
      this.play();
    }

    if (prevProps.volume !== this.props.volume) {
      this.setVolume(this.props.volume);
    }
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

  render() {
    return null;
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
      if (this.props.playing) {
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
}

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

    this._onChange = this._onChange.bind(this);
  }

  componentDidMount() {
    VolumeStore.addChangeListener(this._onChange);
    NowPlayingStore.addChangeListener(this._onChange);
  }

  componentWillUnmount() {
    VolumeStore.removeChangeListener(this._onChange);
    NowPlayingStore.removeChangeListener(this._onChange);
  }

  render() {
    return <AudioPlayer {...this.state} />;
  }

  _onChange() {
    this.setState(getPlayerState());
  }
}
