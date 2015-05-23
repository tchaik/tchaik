"use strict";

import React from "react/addons";

import classNames from "classnames";

import Icon from "./Icon.js";
import TimeFormatter from "./TimeFormatter.js";
import ArtworkImage from "./ArtworkImage.js";
import NowPlayingStore from "../stores/NowPlayingStore.js";
import PlayingStatusStore from "../stores/PlayingStatusStore.js";

import LeftColumnActions from "../actions/LeftColumnActions.js";
import LeftColumnConstants from "../constants/LeftColumnConstants.js";

import GroupAttributes from "../components/GroupAttributes.js";

function getNowPlayingState() {
  return {
    track: NowPlayingStore.getTrack(),
    buffered: PlayingStatusStore.getBuffered(),
    duration: PlayingStatusStore.getDuration(),
    currentTime: PlayingStatusStore.getTime(),
    error: PlayingStatusStore.getError(),
  };
}

export default class NowPlaying extends React.Component {
  constructor(props) {
    super(props);

    this.state = getNowPlayingState();
    this._onChange = this._onChange.bind(this);
    this._onClickArtwork = this._onClickArtwork.bind(this);
  }

  componentDidMount() {
    NowPlayingStore.addChangeListener(this._onChange);
    PlayingStatusStore.addChangeListener(this._onChange);
  }

  componentWillUnmount() {
    NowPlayingStore.removeChangeListener(this._onChange);
    PlayingStatusStore.removeChangeListener(this._onChange);
  }

  render() {
    var track = this.state.track;
    if (track === null) {
      return <div className="now-playing-track" />;
    }

    var remainingTime = parseInt(this.state.duration) - parseInt(this.state.currentTime);

    var attributes = [];
    var attributesElement = null;

    for (var attribute of ["Artist", "GroupName"]) {
      if (track[attribute]) {
        attributes.push(<span>{track[attribute]}</span>);
      }
    }

    if (attributes) {
      attributesElement = <GroupAttributes list={attributes} />;
    }

    var className = classNames({
      "now-playing-track": true,
      "error": (this.state.error !== null),
    });

    return (
      <div className={className}>
        <ArtworkImage path={`/artwork/${track.TrackID}`} onClick={this._onClickArtwork}/>
        <div className="track-info">
          <div className="container">
            <div className="title">
              {track.Name}
              <span className="hover-show">
                <BitRate track={track} />
                <a className="goto" href={`#track_${track.TrackID}`}>
                  <Icon icon="share-alt" />
                </a>
              </span>
            </div>
            {attributesElement}

            <div className="times">
              <TimeFormatter className="current-time" time={this.state.currentTime} />
              <TimeFormatter className="track-length" time={this.state.duration} />
              <TimeFormatter className="remaining" time={remainingTime} />
            </div>
          </div>
        </div>
      </div>
    );
  }

  _onChange() {
    this.setState(getNowPlayingState());
  }

  _onClickArtwork() {
    LeftColumnActions.mode(LeftColumnConstants.RETRO);
  }
}


class BitRate extends React.Component {
  constructor(props) {
    super(props);

    this.state = {expanded: false};
    this._onClick = this._onClick.bind(this);
  }

  render() {
    var bitRate = null;
    if (this.state.expanded) {
      bitRate = (
        <span className="value">{this.props.track.BitRate} kbps</span>
      );
    }

    var className = classNames({
      "bitrate": true,
      "expanded": this.state.expanded,
    });
    return (
      <span className={className} onClick={this._onClick}>
        <Icon icon="equalizer" />
        {bitRate}
      </span>
    );
  }

  _onClick() {
    this.setState({expanded: !this.state.expanded});
  }
}
