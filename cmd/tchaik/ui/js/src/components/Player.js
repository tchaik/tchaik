"use strict";

import React from "react";

import NowPlayingConstants from "../constants/NowPlayingConstants.js";
import NowPlayingActions from "../actions/NowPlayingActions.js";
import NowPlayingStore from "../stores/NowPlayingStore.js";

import PlayingStatusStore from "../stores/PlayingStatusStore.js";

import { connect } from "react-redux";

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

    if (this.props.source) {
      this.setSrc(this.props.source);
      this.load();
    }

    if (this.props.playing) {
      this.play();
    }
  }

  componentWillUnmount() {
    this._audio.pause();

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

    if (prevProps.volume !== this.props.volume || prevProps.mute !== this.props.mute) {
      this.setVolume(this.props.volume, this.props.mute);
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

  setVolume(volume, mute) {
    if (mute) {
      volume = 0.00
    }
    this._audio.volume = volume;
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
        if (this.props.playing !== true) {
          NowPlayingActions.playing(true);
        }
        break;

      case "pause":
        if (this.props.playing !== false) {
          NowPlayingActions.playing(false);
        }
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
  let src = null;
  let track = NowPlayingStore.getTrack();
  if (track) {
    src = `/track/${track.id}`;
  }

  return {
    source: src,
    playing: NowPlayingStore.getPlaying(),
  };
}

class player extends React.Component {
  constructor(props) {
    super(props);

    this.state = getPlayerState();

    this._onChange = this._onChange.bind(this);
  }

  componentDidMount() {
    NowPlayingStore.addChangeListener(this._onChange);
  }

  componentWillUnmount() {
    NowPlayingStore.removeChangeListener(this._onChange);
  }

  render() {
    return <AudioPlayer {...this.props} {...this.state} />;
  }

  _onChange() {
    this.setState(getPlayerState());
  }
}

const Player = connect(
  state => (state)
)(player)

export default Player
