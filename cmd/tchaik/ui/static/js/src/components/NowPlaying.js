/** @jsx React.DOM */
'use strict';

var React = require('react/addons');

var Icon = require('./Icon.js');
var TimeFormatter = require('./TimeFormatter.js');
var GroupAttributes = require('./GroupAttributes.js');
var ArtworkImage = require('./ArtworkImage.js');

var NowPlayingStore = require('../stores/NowPlayingStore.js');
var NowPlayingActions = require('../actions/NowPlayingActions.js');

var PlayingStatusStore = require('../stores/PlayingStatusStore.js');


var NowPlaying = React.createClass({
  render: function() {
    return (
      <div className="now-playing-track">
        <TrackInfo />
      </div>
    );
  },
});

function getTrackInfoState() {
  return {
    track: NowPlayingStore.getTrack(),
    buffered: PlayingStatusStore.getBuffered(),
    duration: PlayingStatusStore.getDuration(),
    currentTime: PlayingStatusStore.getTime(),
  };
}

var TrackInfo = React.createClass({
  getInitialState: function() {
    return getTrackInfoState();
  },

  componentDidMount: function() {
    NowPlayingStore.addChangeListener(this._onChange);
    PlayingStatusStore.addChangeListener(this._onChange);
  },

  componentWillUnmount: function() {
    NowPlayingStore.removeChangeListener(this._onChange);
    PlayingStatusStore.removeChangeListener(this._onChange);
  },

  render: function() {
    var track = this.state.track;
    var fields = ['Album', 'Artist', 'Year'];
    var attributeArr = [];
    fields.forEach(function(f) {
      if (track[f]) {
        attributeArr.push(track[f]);
      }
    });

    var attributes = <GroupAttributes list={attributeArr} />;
    var remainingTime = parseInt(this.state.duration) - parseInt(this.state.currentTime);

    return (
      <div className="info">
        <ArtworkImage path={"/artwork/" + track.TrackID} />
        <span className="title">{track.Name}<a className="goto" href={"#track_"+track.TrackID}><Icon icon="share-alt" /></a></span>
        {attributes}

        <PlayProgress markerWidth={2} current={this.state.currentTime} duration={this.state.duration} buffered={this.state.buffered} setCurrentTime={NowPlayingActions.setCurrentTime} />
        <span className="times">
          <TimeFormatter className="currentTime" time={this.state.currentTime} />
          <TimeFormatter className="remaining" time={remainingTime} />
        </span>
        <div style={{"clear": "both"}} />
      </div>
    );
  },

  _onChange: function() {
    this.setState(getTrackInfoState());
  }
});

function _getOffsetLeft(elem) {
    var offsetLeft = 0;
    do {
      if (!isNaN(elem.offsetLeft)) {
          offsetLeft += elem.offsetLeft;
      }
    } while ((elem = elem.offsetParent));
    return offsetLeft;
}

var PlayProgress = React.createClass({
  propTypes: {
    markerWidth: React.PropTypes.number.isRequired,
    current: React.PropTypes.number.isRequired,
    buffered: React.PropTypes.number.isRequired,
    duration: React.PropTypes.number.isRequired,
  },

  render: function() {
    var wpc = (this.props.current / this.props.duration) * 100;
    var w = "calc("+Math.min(wpc, 100.0)+"% - " + this.props.markerWidth + "px)";
    var bpc = (this.props.buffered / this.props.duration) * 100 - wpc;
    var b = "calc("+Math.min(bpc, 100.0)+"% - " + this.props.markerWidth + "px)";

    return (
      <span className="playProgress" onMouseDown={this._onMouseDown} onWheel={this._onWheel}>
        <span className="bar">
          <span className="current" style={{width: w}} />
          <span className="marker" style={{width: this.props.markerWidth}} />
          <span className="buffered" style={{width: b}} />
        </span>
      </span>
    );
  },

  _onMouseDown: function(evt) {
    var pos = evt.pageX - _getOffsetLeft(evt.currentTarget);
    var width = evt.currentTarget.offsetWidth;
    this.props.setCurrentTime((pos / width) * this.props.duration);
  },

  _onWheel: function(evt) {
    evt.stopPropagation();
    var t = this.props.current + (0.01 * this.props.duration * evt.deltaY);
    if (t > this.props.duration) {
      t = this.props.duration;
    } else if (t < 0.00) {
      t = 0.0;
    }
    this.props.setCurrentTime(t);
  },
});

module.exports = NowPlaying;
