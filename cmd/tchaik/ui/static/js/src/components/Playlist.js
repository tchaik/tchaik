/** @jsx React.DOM */
'use strict';

var React = require('react/addons');

var classNames = require('classnames');

var Icon = require('./Icon.js');
var TimeFormatter = require('./TimeFormatter.js');
var GroupAttributes = require('./GroupAttributes.js');
var ArtworkImage = require('./ArtworkImage.js');

var CollectionStore = require('../stores/CollectionStore.js');
var CollectionActions = require('../actions/CollectionActions.js');

var PlaylistStore = require('../stores/PlaylistStore.js');
var PlaylistActions = require('../actions/PlaylistActions.js');

var NowPlayingStore = require('../stores/NowPlayingStore.js');


var Playlist = React.createClass({

  getInitialState: function() {
    return {
      list: PlaylistStore.getPlaylist(),
    };
  },

  componentDidMount: function() {
    PlaylistStore.addChangeListener(this._onChange);
  },

  componentWillUnmount: function() {
    PlaylistStore.removeChangeListener(this._onChange);
  },

  render: function() {
    var rootCount = {};
    var items = this.state.list.map(function(item, i) {
      if (!rootCount[item.root]) {
        rootCount[item.root] = 0;
      }
      rootCount[item.root]++;
      return <RootGroup path={item.root} key={item.root+rootCount[item.root]} itemIndex={i} />;
    });

    return (
      <div className="playlist">
        {items}
      </div>
    );
  },

  _onChange: function() {
    this.setState({
      list: PlaylistStore.getPlaylist(),
    });
  },

});

function getItem(path) {
  var c = CollectionStore.getCollection(path);
  if (c === undefined) {
    return null;
  }
  return c;
}

var RootGroup = React.createClass({
  propTypes: {
    path: React.PropTypes.array.isRequired,
    itemIndex: React.PropTypes.number.isRequired,
  },

  getInitialState: function() {
    return {
      item: getItem(this.props.path),
    };
  },

  componentDidMount: function() {
    CollectionStore.addChangeListener(this._onChange);
    CollectionActions.fetch(this.props.path);
  },

  componentWillUnmount: function() {
    CollectionStore.removeChangeListener(this._onChange);
  },

  render: function() {
    if (this.state.item === null) {
      return null;
    }

    return (
      <Group item={this.state.item} path={this.props.path} itemIndex={this.props.itemIndex} />
    );
  },

  _onChange: function(keyPath) {
    if (CollectionStore.pathToKey(this.props.path) === keyPath) {
       this.setState({
         item: getItem(this.props.path)
       });
    }
  },
});

var Group = React.createClass({
  propTypes: {
    path: React.PropTypes.array.isRequired,
    itemIndex: React.PropTypes.number.isRequired,
    item: React.PropTypes.object.isRequired,
  },

  getInitialState: function() {
    return {
      expanded: true,
      common: {},
    };
  },

  setCommon: function(c) {
    this.setState({
      common: c,
    });
  },

  render: function() {
    var groupClasses = {
      'group': true,
      'expanded': this.state.expanded
    };

    var image = null;
    if (this.state.common.TrackID) {
      image = <ArtworkImage path={"/artwork/" + this.state.common.TrackID} />;
    }

    var duration = null;
    if (this.state.common.TotalTime) {
      duration = <TimeFormatter className="duration" time={parseInt(this.state.common.TotalTime/1000)} />;
    }

    var common = this.state.common;
    var fields = ['Artist', 'Composer', 'Year'];
    var attributeArr = [];
    fields.forEach(function(f) {
      if (common[f]) {
        attributeArr.push(common[f]);
      }
    });

    var attributes = null;
    if (attributeArr.length > 0) {
      attributes = <GroupAttributes list={attributeArr} />;
    }

    var content = null;
    if (this.state.expanded) {
      content = <GroupContent path={this.props.path} setCommon={this.setCommon} itemIndex={this.props.itemIndex} />;
    }

    return (
      <div className={classNames(groupClasses)}>
        <div className="name" onClick={this._onClick}>
          {image}
          <span className="name">{this.props.item.Name === "" ? "" : this.props.item.Name}</span>
          <span className="info">
            <Icon icon="remove" onClick={this._onClickRemove} />
            <span className="controls">{duration}</span>
          </span>
          {attributes}
          <div key="clear" style={{clear: 'both'}}/>
        </div>
        {content}
      </div>
    );
  },

  _onClickRemove: function(e) {
    e.stopPropagation();
    PlaylistActions.remove(this.props.itemIndex, this.props.path);
  },

  _onClick: function() {
    this.setState({
      expanded: !this.state.expanded
    });
  },
});

var GroupContent = React.createClass({
  propTypes: {
    path: React.PropTypes.array.isRequired,
    itemIndex: React.PropTypes.number.isRequired,
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

    var pathKeys = PlaylistStore.getItemKeys(this.props.itemIndex, this.props.path);
    if (item.Groups) {
      return <GroupList path={this.props.path} list={item.Groups} itemIndex={this.props.itemIndex} keys={pathKeys.keys} />;
    }
    return <TrackList path={this.props.path} list={item.Tracks} listStyle={item.ListStyle} itemIndex={this.props.itemIndex} keys={pathKeys.keys} />;
  },

  _onChange: function(keyPath) {
    if (CollectionStore.pathToKey(this.props.path) === keyPath) {
      var item = CollectionStore.getCollection(this.props.path);

      var common = {};
      var fields = ['TotalTime', 'Artist', 'Composer', 'TrackID', 'Year'];
      fields.forEach(function(f) {
        if (item[f]) {
          common[f] = item[f];
        }
      });

      if (Object.keys(common).length > 0) {
        this.props.setCommon(common);
      }
      this.setState({item: item});
    }
  }
});

var GroupList = React.createClass({
  propTypes: {
    path: React.PropTypes.array.isRequired,
    itemIndex: React.PropTypes.number.isRequired,
    keys: React.PropTypes.array.isRequired,
  },

  render: function() {
    var keys = this.props.keys;
    var list = this.props.list;

    var itemByKey = {};
    list.forEach(function(item) {
      itemByKey[item.Key] = item;
    });

    var groups = keys.map(function(key) {
      return <Group path={this.props.path.concat([key])} key={key} item={itemByKey[key]} itemIndex={this.props.itemIndex} />;
    }.bind(this));

    return (
      <div className="collection">
        {groups}
      </div>
    );
  },
});

var TrackList = React.createClass({
  propTypes: {
    itemIndex: React.PropTypes.number.isRequired,
    path: React.PropTypes.array.isRequired,
    list: React.PropTypes.array.isRequired,
    keys: React.PropTypes.array.isRequired,
    listStyle: React.PropTypes.string.isRequired,
  },

  render: function() {
    var list = this.props.list;
    var keys = this.props.keys;
    var tracks = keys.map(function(i) {
      return <Track key={list[i].TrackID} data={list[i]} path={this.props.path.concat([i])} itemIndex={this.props.itemIndex} index={i} />;
    }.bind(this));

    return (
      <div className="tracks">
        <ol className={this.props.listStyle}>
          {tracks}
        </ol>
      </div>
    );
  },
});

function pathsEqual(p1, p2) {
  return CollectionStore.pathToKey(p1) === CollectionStore.pathToKey(p2);
}

function isCurrent(i, p) {
  var c = PlaylistStore.getCurrent();
  if (c === null) {
    return false;
  }
  return pathsEqual(c.path, p) && (i === c.item);
}

function isPlaying(trackId) {
  var playing = NowPlayingStore.getPlaying();
  if (!playing) {
    return false;
  }

  var t = NowPlayingStore.getTrack();
  if (t) {
    return t.TrackID === trackId;
  }
  return false;
}

var Track = React.createClass({
  propTypes: {
    itemIndex: React.PropTypes.number.isRequired,
    index: React.PropTypes.number.isRequired,
    path: React.PropTypes.array.isRequired,
  },

  getInitialState: function() {
    return {
      isCurrent: isCurrent(this.props.itemIndex, this.props.path),
      isPlaying: isPlaying(this.props.data.TrackID),
    };
  },

  componentDidMount: function() {
    PlaylistStore.addChangeListener(this._onChange);
    NowPlayingStore.addChangeListener(this._onChange);
  },

  componentWillUnmount: function() {
    PlaylistStore.removeChangeListener(this._onChange);
    NowPlayingStore.removeChangeListener(this._onChange);
  },

  render: function() {
    var durationSecs = parseInt(this.props.data.TotalTime/1000);
    var style = {
      current: this.state.isCurrent,
      'is-playing': this.state.isPlaying,
    };

    return (
      <li onClick={this._onClick} style={{'counterReset': "li "+ (this.props.index+1)}} className={classNames(style)}>
        <span id={"track_"+this.props.data.TrackID} className="name">{this.props.data.Name}</span>
        <span className="info">
          <Icon icon="remove" onClick={this._onClickRemove} />
          <TimeFormatter className="duration" time={durationSecs} />
        </span>
      </li>
    );
  },

  _onClick: function() {
    PlaylistActions.play(this.props.itemIndex, this.props.path, this.props.data);
  },

  _onClickRemove: function(e) {
    e.stopPropagation();
    PlaylistActions.remove(this.props.itemIndex, this.props.path);
  },

  _onChange: function() {
    this.setState({
      isCurrent: isCurrent(this.props.itemIndex, this.props.path),
      isPlaying: isPlaying(this.props.data.TrackID),
    });
  },
});

module.exports = Playlist;
