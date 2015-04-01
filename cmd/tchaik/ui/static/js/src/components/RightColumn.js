/** @jsx React.DOM */
'use strict';

var React = require('react/addons');

var Icon = require('./Icon.js');

var NowPlaying = require('./NowPlaying.js');
var NowPlayingActions = require('../actions/NowPlayingActions.js');
var NowPlayingStore = require('../stores/NowPlayingStore.js');

var Playlist = require('./Playlist.js');
var PlaylistStore = require('../stores/PlaylistStore.js');
var PlaylistActions = require('../actions/PlaylistActions.js');

var BACKWARD_TIMEOUT = 2000;

var RightColumn = React.createClass({

  render: function() {
    return (
      <div className="now-playing">
        <NowPlaying />
        <Controls />
        <Volume width={100} markerWidth={2} />
        <Playlist />
      </div>
    );
  },

});

var Controls = React.createClass({
  getInitialState: function() {
    return {
      playing: NowPlayingStore.getPlaying(),
      canNext: PlaylistStore.canNext(),
      canPrev: PlaylistStore.canPrev(),
    };
  },

  componentDidMount: function() {
    NowPlayingStore.addChangeListener(this._onChange);
    PlaylistStore.addChangeListener(this._onChangePlaylist);
  },

  componentWillUnmount: function() {
    NowPlayingStore.removeChangeListener(this._onChange);
    PlaylistStore.removeChangeListener(this._onChangePlaylist);
  },

  render: function() {
    var prevClasses = {'enabled': this.state.canPrev};
    var nextClasses = {'enabled': this.state.canNext};
    return (
      <div className="controls">
        <Icon icon="step-backward" extraClasses={prevClasses} onClick={this._onBackward} />
        <Icon icon={this.state.playing ? "pause" : "play"} onClick={this._togglePlayPause} />
        <Icon icon="step-forward" extraClasses={nextClasses} onClick={this._onForward} />
      </div>
    );
  },

  _onChangePlaylist: function() {
    this.setState({
      canNext: PlaylistStore.canNext(),
      canPrev: PlaylistStore.canPrev(),
    });
  },

  _onChange: function() {
    this.setState({playing: NowPlayingStore.getPlaying()});

    if (this.state.playing) {
      document.title = "tchaik: " + NowPlayingStore.getCurrent().Name;
    }
  },

  _togglePlayPause: function() {
    NowPlayingActions.playing(!this.state.playing);
    this.setState({
      playing: !this.state.playing,
    });
  },

  _backButtonTimerRunning: function() {
    if (this._backButtonTimer) {
      return true;
    }
    return false;
  },

  _backButtonTimerStart: function() {
    if (this._backButtonTimer) {
      clearTimeout(this._backButtonTimer);
    }
    this._backButtonTimer = setTimeout(this._backButtonTimerTimeout, BACKWARD_TIMEOUT);
  },

  _backButtonTimerTimeout: function() {
    this._backButtonTimer = null;
  },

  _prev: function() {
    PlaylistActions.prev();
  },

  _onBackward: function() {
    if (this._backButtonTimerRunning()) {
      PlaylistActions.prev();
    } else if (this.state.playing || NowPlayingStore.getCurrentTime() > 0) {
      NowPlayingActions.restartTrack();
    } else {
      PlaylistActions.prev();
    }
    this._backButtonTimerStart();
    return;
  },

  _onForward: function() {
    PlaylistActions.next();
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

var Volume = React.createClass({
  propTypes: {
    width: React.PropTypes.number.isRequired,
    markerWidth: React.PropTypes.number.isRequired,
  },

  getInitialState: function() {
    return {volume: NowPlayingStore.getVolume()};
  },

  componentDidMount: function() {
    NowPlayingStore.addChangeListener(this._onChange);
  },

  componentWillUnmount: function() {
    NowPlayingStore.removeChangeListener(this._onChange);
  },

  render: function() {
    var classSuffix;
    if (this.state.volume === 0.00) {
      classSuffix = 'off';
    } else if (this.state.volume < 0.5) {
      classSuffix = 'down';
    } else {
      classSuffix = 'up';
    }

    var w = parseInt(this.state.volume * this.props.width);
    var overflow = (w + this.props.markerWidth) - this.props.width;
    if (overflow > 0) {
      w -= overflow;
    }

    return (
      <div className="volume" onMouseDown={this._onMouseDown} onWheel={this._onWheel} style={{width: this.props.width}}>
        <Icon icon={'volume-' + classSuffix} onMouseDown={this._toggleMute} />
        <span className="bar">
          <span className="current" style={{width: w}} />
          <span className="marker" style={{width: this.props.markerWidth}} />
        </span>
      </div>
    );
  },

  _toggleMute: function(evt) {
    evt.stopPropagation();
    var v = (this.state.volume === 0.00) ? 0.75 : 0.00;
    NowPlayingActions.volume(v);
  },

  _onWheel: function(evt) {
    evt.stopPropagation();
    var v = this.state.volume + 0.05 * evt.deltaY;
    if (v > 1.0) {
      v = 1.0;
    } else if (v < 0.00) {
      v = 0.0;
    }
    NowPlayingActions.volume(v);
  },

  _onMouseDown: function(evt) {
    var pos = evt.pageX - _getOffsetLeft(evt.currentTarget);
    var width = evt.currentTarget.offsetWidth;

    NowPlayingActions.volume(pos/width);
  },

  _onChange: function() {
    this.setState({volume: NowPlayingStore.getVolume()});
  },
});

module.exports = RightColumn;
