/** @jsx React.DOM */
'use strict';

var React = require('react/addons');

var Icon = require('./Icon.js');
var TimeFormatter = require('./TimeFormatter.js');
var GroupAttributes = require('./GroupAttributes.js');
var ArtworkImage = require('./ArtworkImage.js');

var NowPlayingStore = require('../stores/NowPlayingStore.js');
var NowPlayingActions = require('../actions/NowPlayingActions.js');
var NowPlayingConstants = require('../constants/NowPlayingConstants.js');

var PlayingStatusStore = require('../stores/PlayingStatusStore.js');


var NowPlaying = React.createClass({
  getInitialState: function() {
    return {
      track: NowPlayingStore.getTrack(),
      playing: NowPlayingStore.getPlaying(),
      currentTime: PlayingStatusStore.getTime(),
    };
  },

  onPlayerEvent: function(evt) {
    switch (evt.type) {

    case "error":
      console.log("Error received from <audio> tag:");
      console.error(evt);
      break;

    case "progress":
      var range = this.buffered();
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
      NowPlayingActions.currentTime(this.currentTime());
      break;

    case "loadedmetadata":
      NowPlayingActions.setDuration(this.duration());

      this.setCurrentTime(PlayingStatusStore.getTime());
      if (this.state.playing) {
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
  },

  audioDOMNode: function() {
    return React.findDOMNode(this.refs.ref_audio);
  },

  buffered: function() {
    return this.audioDOMNode().buffered;
  },

  play: function() {
    this.audioDOMNode().play();
  },

  pause: function() {
    this.audioDOMNode().pause();
  },

  load: function() {
    this.audioDOMNode().load();
  },

  setCurrentTime: function(t) {
    this.audioDOMNode().currentTime = t;
  },

  currentTime: function() {
    return this.audioDOMNode().currentTime;
  },

  setVolume: function(t) {
    this.audioDOMNode().volume = t;
  },

  volume: function() {
    return this.audioDOMNode().volume;
  },

  duration: function() {
    return this.audioDOMNode().duration;
  },

  componentDidMount: function() {
    var audioNode = this.audioDOMNode();
    audioNode.addEventListener('error', this.onPlayerEvent);
    audioNode.addEventListener('progress', this.onPlayerEvent);
    audioNode.addEventListener('play', this.onPlayerEvent);
    audioNode.addEventListener('pause', this.onPlayerEvent);
    audioNode.addEventListener('ended', this.onPlayerEvent);
    audioNode.addEventListener('timeupdate', this.onPlayerEvent);
    audioNode.addEventListener('loadedmetadata', this.onPlayerEvent);
    audioNode.addEventListener('loadstart', this.onPlayerEvent);

    NowPlayingStore.addChangeListener(this._onChange);
    NowPlayingStore.addControlListener(this._onControl);

    var volume = NowPlayingStore.getVolume();
    if (NowPlayingStore.getVolumeMute()) {
      volume = 0.0;
    }
    this.setVolume(volume);
  },

  componentWillUnmount: function() {
    var audioNode = this.audioDOMNode();
    audioNode.removeEventListener('error', this.onPlayerEvent);
    audioNode.removeEventListener('progress', this.onPlayerEvent);
    audioNode.removeEventListener('play', this.onPlayerEvent);
    audioNode.removeEventListener('pause', this.onPlayerEvent);
    audioNode.removeEventListener('ended', this.onPlayerEvent);
    audioNode.removeEventListener('timeupdate', this.onPlayerEvent);
    audioNode.removeEventListener('loadedmetadata', this.onPlayerEvent);
    audioNode.removeEventListener('loadstart', this.onPlayerEvent);

    NowPlayingStore.removeChangeListener(this._onChange);
    NowPlayingStore.removeControlListener(this._onControl);
  },

  componentDidUpdate: function(prevProps, prevState) {
    if (this.state.track) {
      if (prevState.track === null || prevState.track.TrackID !== this.state.track.TrackID) {
        this.load();
        this.play();
      }
    }

    if (prevState.playing != this.state.playing) {
      if (this.state.playing) {
        this.play();
      } else {
        this.pause();
      }
    }

    var volume = NowPlayingStore.getVolume();
    if (this.volume() != volume) {
      this.setVolume(volume);
    }
  },

  render: function() {
    var trackID = "0";
    var source = null;
    var trackInfo = null;

    if (this.state.track) {
      trackID = this.state.track.TrackID;
      source = <source src={"/track/"+trackID} />;
      trackInfo = <TrackInfo />;
    }

    return (
      <div className="now-playing-track">
        {trackInfo}
        <audio id={"player_"+trackID} ref="ref_audio">
         {source}
        </audio>
      </div>
    );
  },

  _onChange: function() {
    this.setState({
      track: NowPlayingStore.getTrack(),
      playing: NowPlayingStore.getPlaying(),
    });

    setTimeout(0, function() {
      this.refs.now_playing.load();
    }.bind(this));
  },

  _onControl: function(type, value) {
    switch (type) {
      case NowPlayingConstants.SET_CURRENT_TIME:
        this.setCurrentTime(value);
        break;
    }
  }
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
