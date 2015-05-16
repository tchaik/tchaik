'use strict';

import React from 'react/addons';

import classNames from 'classnames';

import Icon from './Icon.js';
import TimeFormatter from './TimeFormatter.js';
import GroupAttributes from './GroupAttributes.js';
import ArtworkImage from './ArtworkImage.js';

import NowPlayingStore from '../stores/NowPlayingStore.js';

import PlayingStatusStore from '../stores/PlayingStatusStore.js';


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
      return null;
    }

    var fields = ['Album', 'Artist', 'AlbumArtist', 'Composer', 'Year'];
    if (track.Artist && track.AlbumArtist) {
      fields = ['Album', 'Artist', 'Composer', 'Year'];
    }
    var attributeArr = [];
    fields.forEach(function(f) {
      if (track[f]) {
        attributeArr.push(track[f]);
      }
    });

    var attributes = <GroupAttributes list={attributeArr} />;
    var remainingTime = parseInt(this.state.duration) - parseInt(this.state.currentTime);

    var className = classNames({
      'now-playing-track': true,
      'error': (this.state.error !== null),
    });

    return (
      <div className={className}>
        <ArtworkImage path={"/artwork/" + track.TrackID} />
        <div className="track-info">
          <div className="title">{track.Name}<BitRate track={track} /><a className="goto" href={"#track_"+track.TrackID}><Icon icon="share-alt" /></a></div>
          {attributes}

          <div className="times">
            <TimeFormatter className="current-time" time={this.state.currentTime} />
            <TimeFormatter className="track-length" time={this.state.duration} />
            <TimeFormatter className="remaining" time={remainingTime} />
          </div>
        </div>
      </div>
    );
  }

  _onChange() {
    this.setState(getNowPlayingState());
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
      'bitrate': true,
      'expanded': this.state.expanded,
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
