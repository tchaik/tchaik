'use strict';

import React from 'react/addons';

import classNames from 'classnames';

import Icon from './Icon.js';
import TimeFormatter from './TimeFormatter.js';
import GroupAttributes from './GroupAttributes.js';
import ArtworkImage from './ArtworkImage.js';

import NowPlayingStore from '../stores/NowPlayingStore.js';
import NowPlayingActions from '../actions/NowPlayingActions.js';

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

    var fields = ['Album', 'Artist', 'Year'];
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

          <PlayProgress markerWidth={10} current={this.state.currentTime} duration={this.state.duration} buffered={this.state.buffered} setCurrentTime={NowPlayingActions.setCurrentTime} />
          <div className="times">
            <TimeFormatter className="currentTime" time={this.state.currentTime} />
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

function _getOffsetLeft(elem) {
    var offsetLeft = 0;
    do {
      if (!isNaN(elem.offsetLeft)) {
          offsetLeft += elem.offsetLeft;
      }
    } while ((elem = elem.offsetParent));
    return offsetLeft;
}

class PlayProgress extends React.Component {
  constructor(props) {
    super(props);

    this._onClick = this._onClick.bind(this);
    this._onWheel = this._onWheel.bind(this);
  }

  render() {
    var wpc = (this.props.current / this.props.duration) * 100;
    var w = "calc("+Math.min(wpc, 100.0)+"% - " + this.props.markerWidth + "px)";
    var bpc = (this.props.buffered / this.props.duration) * 100 - wpc;
    var b = "calc("+Math.min(bpc, 100.0)+"% - " + this.props.markerWidth + "px)";

    return (
      <div className="playProgress" onClick={this._onClick} onWheel={this._onWheel}>
        <div className="bar">
          <div className="current" style={{width: w}} />
          <div className="marker" style={{width: this.props.markerWidth}} />
          <div className="buffered" style={{width: b}} />
        </div>
      </div>
    );
  }

  _onClick(evt) {
    var pos = evt.pageX - _getOffsetLeft(evt.currentTarget);
    var width = evt.currentTarget.offsetWidth;
    this.props.setCurrentTime((pos / width) * this.props.duration);
  }

  _onWheel(evt) {
    evt.stopPropagation();
    var t = this.props.current + (0.01 * this.props.duration * evt.deltaY);
    if (t > this.props.duration) {
      t = this.props.duration;
    } else if (t < 0.00) {
      t = 0.0;
    }
    this.props.setCurrentTime(t);
  }
}

PlayProgress.propTypes = {
  markerWidth: React.PropTypes.number.isRequired,
  current: React.PropTypes.number.isRequired,
  buffered: React.PropTypes.number.isRequired,
  duration: React.PropTypes.number.isRequired,
};
