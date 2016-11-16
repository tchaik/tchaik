"use strict";

import React from "react";

import classNames from "classnames";

import ArtworkImage from "./ArtworkImage.js";
import ContainerActions from "../actions/ContainerActions.js";
import ContainerConstants from "../constants/ContainerConstants.js";
import GroupAttributes from "../components/GroupAttributes.js";
import Icon from "./Icon.js";
import NowPlayingStore from "../stores/NowPlayingStore.js";
import PlayingStatusStore from "../stores/PlayingStatusStore.js";
import TimeFormatter from "./TimeFormatter.js";

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
    const track = this.state.track;
    if (track === null) {
      return <div className="track" />;
    }

    const remainingTime = parseInt(this.state.duration) - parseInt(this.state.currentTime);

    const className = classNames({
      "track": true,
      "error": (this.state.error !== null),
    });

    return (
      <div className={className}>
        <ArtworkImage path={`/artwork/${track.id}`} onClick={this._onClickArtwork}/>
        <div className="info">
          <div className="wrapper">
            <div className="title">
              {track.name}
              <span className="hover-show">
                <a className="goto" href={`#track_${track.id}`}>
                  <Icon icon="reply" />
                </a>
              </span>
            </div>
            <GroupAttributes data={track} attributes={["artist", "groupName", "composer"]} />
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
    ContainerActions.mode(ContainerConstants.RETRO);
  }
}

class BitRate extends React.Component {
  constructor(props) {
    super(props);

    this.state = {expanded: false};
    this._onClick = this._onClick.bind(this);
  }

  shouldComponentUpdate(nextProps, nextState) {
    if (this.props.track.bitRate != nextProps.track.bitRate) {
      return true;
    }
    if (this.state.expanded != nextState.expanded) {
      return true;
    }
    return false;
  }

  render() {
    let bitRate = null;
    if (this.state.expanded) {
      bitRate = (
        <span className="value">{this.props.track.bitRate} kbps</span>
      );
    }

    const className = classNames({
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
