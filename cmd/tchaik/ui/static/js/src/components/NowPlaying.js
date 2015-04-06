/** @jsx React.DOM */
'use strict';

var React = require('react/addons');

var classNames = require('classnames');

var Icon = require('./Icon.js');
var TimeFormatter = require('./TimeFormatter.js');
var GroupAttributes = require('./GroupAttributes.js');

var NowPlayingStore = require('../stores/NowPlayingStore.js');
var NowPlayingActions = require('../actions/NowPlayingActions.js');


var NowPlaying = React.createClass({
  getInitialState: function() {
    return {
      track: NowPlayingStore.getCurrent(),
      playing: NowPlayingStore.getPlaying(),
      currentTime: NowPlayingStore.getCurrentTime(),
      buffered: 0,
      duration: 0,
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
        this.setState({
          buffered: range.end(range.length-1),
        });
      }
      break;

    case "play":
      this.setState({playing: true});
      NowPlayingActions.playing(true);
      break;

    case "pause":
      this.setState({playing: false});
      NowPlayingActions.playing(false);
      break;

    case "ended":
      NowPlayingActions.ended(NowPlayingStore.getCurrentSource());
      break;

    case "timeupdate":
      var t = this.currentTime();
      this.setState({currentTime: t});
      NowPlayingActions.currentTime(t);
      break;

    case "loadedmetadata":
      this.setState({
        duration: this.duration()
      });

      this.setCurrentTime(NowPlayingStore.getCurrentTime());
      if (this.state.playing) {
        this.play();
      }
      break;

    case "loadstart":
      this.setState({
        duration: 0,
        currentTime: 0,
        buffered: 0,
      });
      break;

    default:
      console.warn("unhandled player event:");
      console.warn(evt);
      break;
    }
  },

  audioDOMNode: function() {
    return this.refs.ref_audio.getDOMNode();
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

    this.setVolume(NowPlayingStore.getVolume());
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

    if (this.volume() != NowPlayingStore.getVolume()) {
      this.setVolume(NowPlayingStore.getVolume());
    }
  },

  render: function() {
    var trackInfo = null;
    var trackID = "0";
    var source = null;

    if (this.state.track) {
      trackInfo = <TrackInfo track={this.state.track}
                       currentTime={this.state.currentTime}
                          duration={this.state.duration}
                          buffered={this.state.buffered}
                    setCurrentTime={this.setCurrentTime} />;
      trackID = this.state.track.TrackID;
      source = <source src={"/track/"+trackID} />;
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
      track: NowPlayingStore.getCurrent(),
      playing: NowPlayingStore.getPlaying(),
    });

    setTimeout(0, function() {
      this.refs.now_playing.load();
    }.bind(this));
  },
});

var TrackInfo = React.createClass({
  getInitialState: function() {
    return {
      showImage: false,
      largeImage: true,
    };
  },

  showImage: function() {
    this.setState({showImage: true});
  },

  onError: function() {
    this.setState({showImage: false});
  },

  toggleLargeImage: function() {
    this.setState({
      largeImage: !this.state.largeImage,
    });
  },

  render: function() {
    var remainingTime = parseInt(this.props.duration) - parseInt(this.props.currentTime);

    var imgClassSet = {
      'visible': this.state.showImage,
      'large': this.state.showImage && this.state.largeImage,
    };

    var attributeArr = [];
    if (this.props.track.Album) {
      attributeArr.push(this.props.track.Album);
    }

    if (this.props.track.Artist) {
      attributeArr.push(this.props.track.Artist);
    }

    if (this.props.track.Year) {
      attributeArr.push(this.props.track.Year);
    }

    var attributes = <GroupAttributes list={attributeArr} />;
//        <span className="album">{this.props.track.Album}</span>
    return (
      <div className="info">
        <span className="image">
          <img src={"/artwork/" + this.props.track.TrackID} key="img" className={classNames(imgClassSet)} onLoad={this.showImage} onError={this.onError} />
        </span>
        <span className="title">{this.props.track.Name}<a className="goto" href={"#track_"+this.props.track.TrackID}><Icon icon="share-alt" /></a></span>
        {attributes}

        <PlayProgress markerWidth={2} current={this.props.currentTime} duration={this.props.duration} buffered={this.props.buffered} setCurrentTime={this.props.setCurrentTime} />
        <span className="times">
          <TimeFormatter className="currentTime" time={this.props.currentTime} />
          <TimeFormatter className="remaining" time={remainingTime} />
        </span>
        <div style={{"clear": "both"}} />
      </div>
    );
  },
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
    var w = "calc("+wpc+"% - " + this.props.markerWidth + "px)";
    var bpc = (this.props.buffered / this.props.duration) * 100 - wpc;
    var b = "calc("+bpc+"% - " + this.props.markerWidth + "px)";

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
