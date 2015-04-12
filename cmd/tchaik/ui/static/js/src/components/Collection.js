/** @jsx React.DOM */
'use strict';

var React = require('react/addons');

var classNames = require('classnames');

var Icon = require('./Icon.js');
var TimeFormatter = require('./TimeFormatter.js');
var GroupAttributes = require('./GroupAttributes.js');

var CollectionStore = require('../stores/CollectionStore.js');
var CollectionActions = require('../actions/CollectionActions.js');

var NowPlayingStore = require('../stores/NowPlayingStore.js');

var RootCollection = React.createClass({
  render: function() {
    return <GroupContent path={["Root"]} depth={0} />;
  }
});

var Group = React.createClass({
  propTypes: {
    path: React.PropTypes.array.isRequired,
    item: React.PropTypes.object.isRequired,
    depth: React.PropTypes.number.isRequired
  },

  getInitialState: function() {
    return {
      expanded: (this.props.depth !== 1) || CollectionStore.isExpanded(this.props.path),
      common: {},
      showImage: false,
    };
  },

  setCommon: function(c) {
    this.setState({
      common: c,
    });
  },

  showImage: function() {
    this.setState({showImage: true});
  },

  render: function() {
    var content = null;
    var play = null;
    var headerDiv = null;
    var image = null;

    if (this.state.expanded) {
      content = [
        <GroupContent path={this.props.path} depth={this.props.depth} setCommon={this.setCommon} key="GroupContent0" />
      ];

      if (this.props.depth === 1) {
        content.push(
          <div style={{clear: 'both'}} key="GroupContent1" />
        );
      }

      var duration = null;
      if (this.state.common.totalTime) {
        duration = <TimeFormatter className="duration" time={parseInt(this.state.common.totalTime/1000)} />;
      }

      var attributeArr = [];

      if (this.state.common.artist) {
        attributeArr.push(this.state.common.artist);
      }

      if (this.state.common.composer) {
         attributeArr.push(this.state.common.composer);
      }

      if (this.state.common.year) {
        attributeArr.push(this.state.common.year);
      }

      var attributes = null;
      if (attributeArr.length > 0) {
        attributes = <GroupAttributes list={attributeArr} />;
      }

      play = (
        <span className="controls">
          <Icon icon="play" title="Play Now" onClick={this._onPlayNow} />
          <Icon icon="list" title="Queue" onClick={this._onQueue} />
          {duration}
          {attributes}
        </span>
      );

      if (this.state.common.trackId && this.props.depth == 1) {
        image = (
          <img src={"/artwork/" + this.state.common.trackId}
               key="img"
               className={this.state.showImage === true ? "visible" : ""}
               onLoad={this.showImage} />
        );
      }
    }

    // FIXME: this is to get around the "everthing else" group which has no name
    // (and can't be closed)
    var itemName = this.props.item.Name;
    if (this.props.depth === 1 && itemName === "") {
      itemName = "Misc";
    }

    if (itemName !== "") {
      headerDiv = (
        <div className="header">
          {image}
          <span className="name" onClick={this._onClickName}>{itemName}</span>
          {play}
        </div>
      );
    }

    var groupClasses = {
      'group': true,
      'untitled': this.props.item.Name === "",
      'expanded': this.state.expanded,
    };

    return (
      <div className={classNames(groupClasses)} onClick={this._onClick}>
        {headerDiv}
        {content}
      </div>
    );
  },

  _onClick: function(e) {
    if (!this.state.expanded) {
      this._onClickName(e);
    }
  },

  _onClickName: function(e) {
    e.stopPropagation();
    this.setState({
      expanded: !this.state.expanded
    });
    CollectionActions.expandPath(this.props.path, !this.state.expanded);
  },

  _onPlayNow: function(e) {
    CollectionActions.playNow(this.props.path);
    e.stopPropagation();
  },

  _onQueue: function(e) {
    CollectionActions.appendToPlaylist(this.props.path);
    e.stopPropagation();
  },

});

var GroupContent = React.createClass({
  propTypes: {
    path: React.PropTypes.array.isRequired,
    depth: React.PropTypes.number.isRequired,
  },

  getInitialState: function() {
    return {item: null};
  },

  componentDidMount: function() {
    CollectionStore.addChangeListener(this._onChange);
    CollectionActions.fetch(this.props.path);
  },

  componentWillUnmount: function() {
    CollectionStore.removeChangeListener(this._onChange);
  },

  render: function() {
    var item = this.state.item;
    if (item === null) {
      return null;
    }

    if (item.Groups) {
      return <GroupList path={this.props.path} depth={this.props.depth} list={item.Groups} />;
    }
    return <TrackList path={this.props.path} list={item.Tracks} listStyle={item.ListStyle} />;
  },

  _onChange: function(keyPath) {
    if (CollectionStore.pathToKey(this.props.path) === keyPath) {
      var item = CollectionStore.getCollection(this.props.path);

      var common = {};
      var empty = true;
      if (item.TotalTime) {
        common.totalTime = item.TotalTime;
        empty = false;
      }

      if (item.Artist) {
        common.artist = item.Artist;
        empty = false;
      }

      if (item.TrackID) {
        common.trackId = item.TrackID;
        empty = false;
      }

      if (item.Year) {
        common.year = item.Year;
        empty = false;
      }

      if (item.Composer) {
        common.composer = item.Composer;
        empty = false;
      }

      if (!empty) {
        this.props.setCommon(common);
      }

      this.setState({item: item});
    }
  }
});

var GroupList = React.createClass({
  propTypes: {
    path: React.PropTypes.array.isRequired,
    depth: React.PropTypes.number.isRequired,
    list: React.PropTypes.array.isRequired,
  },

  render: function() {
    var list = this.props.list.map(function(item) {
      return <Group path={this.props.path.concat([item.Key])} depth={this.props.depth + 1} item={item} key={item.Key} />;
    }.bind(this));

    return (
      <div className="collection">
        {list}
      </div>
    );
  },
});

var TrackList = React.createClass({
  propTypes: {
    path: React.PropTypes.array.isRequired,
    list: React.PropTypes.array.isRequired,
    listStyle: React.PropTypes.string.isRequired,
  },

  getInitialState: function() {
    return {};
  },

  render: function() {
    var discs = {};
    var discIndices = [];
    var trackNumber = 0;
    this.props.list.forEach(function(track) {
      var discNumber = track.DiscNumber;
      if (!discs[discNumber]) {
        discs[discNumber] = [];
        discIndices.push(discNumber);
      }
      track.Key = ""+(trackNumber++);
      discs[discNumber].push(track);
    });

    var ols = [];
    var buildTrack = function(track) {
      return <Track key={track.TrackID} data={track} path={this.props.path.concat([track.Key])} />;
    }.bind(this);

    for (var i = 0; i < discIndices.length; i++) {
      var disc = discs[discIndices[i]];
      ols.push(
        <ol key={"disc"+discIndices[i]} className={this.props.listStyle}>
          {disc.map(buildTrack)}
        </ol>
      );
    }

    return (
      <div className="tracks">
        {ols}
      </div>
    );
  },
});

function isCurrent(trackId) {
  var t = NowPlayingStore.getCurrent();
  if (t) {
    return t.TrackID === trackId;
  }
  return false;
}

function getTrackState(trackID) {
  return {
    current: isCurrent(trackID),
    playing: NowPlayingStore.getPlaying(),
  };
}

var Track = React.createClass({
  propTypes: {
    path: React.PropTypes.array.isRequired,
  },

  getInitialState: function() {
    return getTrackState(this.props.data.TrackID);
  },

  componentDidMount: function() {
    NowPlayingStore.addChangeListener(this._onChange);
  },

  componentWillUnmount: function() {
    NowPlayingStore.removeChangeListener(this._onChange);
  },

  render: function() {
    var durationSecs = parseInt(this.props.data.TotalTime/1000);
    var liClasses = {
      'current': this.state.current,
      'playing': this.state.current && this.state.playing,
    };
    return (
      <li className={classNames(liClasses)} onClick={this._onClick}>
        <span id={"track_"+this.props.data.TrackID} className="name">{this.props.data.Name}</span>
        <TimeFormatter className="duration" time={durationSecs} />
        <span className="controls">
          <Icon icon="play" title="Play Now" onClick={this._onPlayNow} />
          <Icon icon="list" title="Queue" onClick={this._onQueue} />
        </span>
      </li>
    );
  },

  _onChange: function() {
    this.setState(getTrackState(this.props.data.TrackID));
  },

  _onClick: function(e) {
    e.stopPropagation();
    CollectionActions.setCurrentTrack(this.props.data);
  },

  _onPlayNow: function(e) {
    e.stopPropagation();
    CollectionActions.playNow(this.props.path.concat([this.props.Key]));
  },

  _onQueue: function(e) {
    e.stopPropagation();
    CollectionActions.appendToPlaylist(this.props.path.concat([this.props.Key]));
  },
});

module.exports.RootCollection = RootCollection;
module.exports.Group = Group;
