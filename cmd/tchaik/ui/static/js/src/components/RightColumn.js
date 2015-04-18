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
        <span><Icon icon={this.state.playing ? "pause" : "play"} onClick={this._togglePlayPause} /></span>
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

    var favicon = document.querySelector("head link[rel=\"shortcut icon\"]");
    var currentTrack = NowPlayingStore.getTrack();
    if (currentTrack === null) {
      document.title = "tchaik";
      return;
    }
    document.title = currentTrack.Name;
    favicon.href = "/icon/" + currentTrack.TrackID;
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
    } else if (this.state.playing || NowPlayingStore.getTime() > 0) {
      NowPlayingActions.setCurrentTime(0);
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

module.exports = RightColumn;
