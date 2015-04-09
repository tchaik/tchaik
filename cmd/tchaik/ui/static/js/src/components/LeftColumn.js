/** @jsx React.DOM */
'use strict';

var React = require('react/addons');

var classNames = require('classnames');

var Icon = require('./Icon.js');
var StatusView = require('./Status.js');
var ApiKeyView = require('./ApiKey.js');

var RootCollection = require('./Collection.js').RootCollection;
var Search = require('./Search.js');
var Covers = require('./Covers.js');

var LeftColumnStore = require('../stores/LeftColumnStore.js');
var LeftColumnActions = require('../actions/LeftColumnActions.js');

var NowPlayingStore = require('../stores/NowPlayingStore.js');
var NowPlayingActions = require('../actions/NowPlayingActions.js');


function getToolBarItemState(mode) {
  return {
    selected: mode === LeftColumnStore.getMode(),
  };
}

var ToolbarItem = React.createClass({
  getInitialState: function() {
    return getToolBarItemState(this.props.mode);
  },

  componentDidMount: function() {
    LeftColumnStore.addChangeListener(this._onChange);
  },

  componentWillUnmount: function() {
    LeftColumnStore.removeChangeListener(this._onChange);
  },

  render: function() {
    var {...other} = this.props;
    var classes = {
      item: true,
      toolbar: true,
      selected: this.state.selected,
    };
    return (
      <span className={classNames(classes)} onClick={this._onClick}>
        <Icon {...other} />
      </span>
    );
  },

  _onClick: function() {
    LeftColumnActions.mode(this.props.mode);
  },

  _onChange: function() {
    this.setState(getToolBarItemState(this.props.mode));
  }
});

function leftColumnState() {
  return {
    mode: LeftColumnStore.getMode(),
  };
}

var LeftColumn = React.createClass({
  getInitialState: function() {
    return leftColumnState();
  },

  componentDidMount: function() {
    LeftColumnStore.addChangeListener(this._onChange);
  },

  componentWillUnmount: function() {
    LeftColumnStore.removeChangeListener(this._onChange);
  },

  render: function() {
    var container = null;
    switch (this.state.mode) {
    case "All":
      container = <RootCollection />;
      break;
    case "Search":
      container = <Search />;
      break;
    case "Covers":
      container = <Covers />;
      break;
    }

    return (
      <div>
        <div id="header">
          <ToolbarItem mode="Search" icon="search" title="Search" />
          <ToolbarItem mode="All" icon="align-justify" title="All" />
          <ToolbarItem mode="Browse" icon="list" title="Albums" />
          <ToolbarItem mode="Covers" icon="th-large" title="Covers" />
          <div className="bottom">
            <Volume height={100} markerHeight={2} />
            <StatusView />
            <ApiKeyView />
          </div>
        </div>
        <div id="container">
          {container}
        </div>
      </div>
    );
  },

  _onChange: function() {
    this.setState(leftColumnState());
  },
});

function _getOffsetTop(elem) {
  var offsetTop = 0;
  do {
    if (!isNaN(elem.offsetTop)) {
        offsetTop += elem.offsetTop;
    }
  } while ((elem = elem.offsetParent));
  return offsetTop;
}

function getVolumeStatus() {
  return {
    volume: NowPlayingStore.getVolume(),
  };
}

var Volume = React.createClass({
  propTypes: {
    height: React.PropTypes.number.isRequired,
    markerHeight: React.PropTypes.number.isRequired,
  },

  getInitialState: function() {
    return getVolumeStatus();
  },

  componentDidMount: function() {
    NowPlayingStore.addChangeListener(this._onChange);
  },

  componentWillUnmount: function() {
    NowPlayingStore.removeChangeListener(this._onChange);
  },

  render: function() {
    var volume = this.state.volume;
    var classSuffix;
    if (volume === 0.00) {
      classSuffix = 'off';
    } else if (volume < 0.5) {
      classSuffix = 'down';
    } else {
      classSuffix = 'up';
    }

    var h = this.props.height - parseInt(volume * this.props.height);
    return (
      <div className="volume" onWheel={this._onWheel} >
        <span className="bar" onMouseOver={this._onMouseOver} onMouseDown={this._onMouseDown} style={{height: this.props.height}}>
          <span className="current" style={{height: h}} />
          <span className="marker" style={{height: this.props.markerHeight}} />
        </span>
        <Icon icon={'volume-' + classSuffix} onMouseDown={this._toggleMute} />
      </div>
    );
  },

  _toggleMute: function(evt) {
    evt.stopPropagation();
    NowPlayingActions.toggleVolumeMute();
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
    var pos = evt.pageY - _getOffsetTop(evt.currentTarget);
    var height = evt.currentTarget.offsetHeight;
    NowPlayingActions.volume(1 - pos/height);
  },

  _onChange: function() {
    this.setState(getVolumeStatus());
  },
});


module.exports = LeftColumn;
