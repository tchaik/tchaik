"use strict";

import React from "react";

import classNames from "classnames";

import Icon from "./Icon.js";
import TimeFormatter from "./TimeFormatter.js";
import GroupAttributes from "./GroupAttributes.js";
import ArtworkImage from "./ArtworkImage.js";

import CollectionStore from "../stores/CollectionStore.js";
import CollectionActions from "../actions/CollectionActions.js";

import PlaylistStore from "../stores/PlaylistStore.js";
import PlaylistActions from "../actions/PlaylistActions.js";

import CursorStore from "../stores/CursorStore.js";
import CursorActions from "../actions/CursorActions.js";


import NowPlayingStore from "../stores/NowPlayingStore.js";


export default class Playlist extends React.Component {
  constructor(props) {
    super(props);

    this.state = {list: PlaylistStore.getPlaylist()};
    this._onChange = this._onChange.bind(this);
  }

  componentDidMount() {
    PlaylistStore.addChangeListener(this._onChange);
    PlaylistActions.fetch();
  }

  componentWillUnmount() {
    PlaylistStore.removeChangeListener(this._onChange);
  }

  render() {
    var items = this.state.list;
    if (items.length === 0) {
      return (
        <div className="playlist">
          <div className="no-items"><Icon icon="queue_music" />Empty playlist</div>
        </div>
      );
    }

    let pathCount = {};
    items = items.map(function(item, i) {
      let path = item.path();
      if (!pathCount[path]) {
        pathCount[path] = 0;
      }
      pathCount[path]++;
      return <RootGroup path={path} key={path + pathCount[path]} itemIndex={i} />;
    });

    return (
      <div className="playlist">
        {items}
      </div>
    );
  }

  _onChange() {
    this.setState({list: PlaylistStore.getPlaylist()});
  }
}


function getItem(path) {
  var c = CollectionStore.getCollection(path);
  if (c === undefined) {
    return null;
  }
  return c;
}

function getRootGroupState(props) {
  return {item: getItem(props.path)};
}

class RootGroup extends React.Component {
  constructor(props) {
    super(props);

    this.state = getRootGroupState(this.props);
    this._onChange = this._onChange.bind(this);
  }

  componentDidMount() {
    CollectionStore.addChangeListener(this._onChange);
    CollectionActions.fetch(this.props.path);
  }

  componentWillUnmount() {
    CollectionStore.removeChangeListener(this._onChange);
  }

  render() {
    if (this.state.item === null) {
      return null;
    }

    return (
      <Group root={true} item={this.state.item} path={this.props.path} itemIndex={this.props.itemIndex} />
    );
  }

  _onChange(keyPath) {
    if (CollectionStore.pathToKey(this.props.path) === keyPath) {
      this.setState(getRootGroupState(this.props));
    }
  }
}

RootGroup.propTypes = {
  path: React.PropTypes.array.isRequired,
  itemIndex: React.PropTypes.number.isRequired,
};


class Group extends React.Component {
  constructor(props) {
    super(props);

    this.state = {expanded: true, common: {}};
    this.setCommon = this.setCommon.bind(this);
    this._onClick = this._onClick.bind(this);
    this._onClickRemove = this._onClickRemove.bind(this);
  }

  setCommon(c) {
    this.setState({common: c});
  }

  render() {
    var groupClasses = {
      "group": true,
      "expanded": this.state.expanded,
    };

    var image = null;
    if (this.props.root && this.state.common.ID) {
      image = <ArtworkImage path={`/artwork/${this.state.common.ID}`} />;
    }

    var duration = null;
    if (this.state.common.TotalTime) {
      duration = <TimeFormatter className="duration" time={parseInt(this.state.common.TotalTime / 1000)} />;
    }

    var common = this.state.common;
    var fields = ["Artist", "Composer", "Year"];
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
        <div className="group-info-container" onClick={this._onClick}>
          {image}
          <div className="group-info">
            <div className="group-details">
              <div className="name">{this.props.item.Name}</div>
              {attributes}
              <div className="attributes duration">{duration}</div>
            </div>
            <div className="controls">
              <Icon icon="clear"onClick={this._onClickRemove} />
            </div>
          </div>
        </div>
        {content}
      </div>
    );
  }

  _onClickRemove(e) {
    e.stopPropagation();
    PlaylistActions.remove(this.props.itemIndex, this.props.path);
  }

  _onClick() {
    this.setState({expanded: !this.state.expanded});
  }
}

Group.propTypes = {
  path: React.PropTypes.array.isRequired,
  itemIndex: React.PropTypes.number.isRequired,
  item: React.PropTypes.object.isRequired,
  root: React.PropTypes.bool,
};


class GroupContent extends React.Component {
  constructor(props) {
    super(props);

    this.state = {item: null};
    this._onChange = this._onChange.bind(this);
  }

  componentDidMount() {
    CollectionStore.addChangeListener(this._onChange);
    CollectionActions.fetch(this.props.path);
  }

  componentWillUnmount() {
    CollectionStore.removeChangeListener(this._onChange);
  }

  render() {
    var item = this.state.item;
    if (item === null) {
      return null;
    }

    var keys = PlaylistStore.getItemKeys(this.props.itemIndex, this.props.path);
    if (item.Groups) {
      return <GroupList path={this.props.path} list={item.Groups} itemIndex={this.props.itemIndex} keys={keys} />;
    }
    return <TrackList path={this.props.path} list={item.Tracks} listStyle={item.ListStyle} itemIndex={this.props.itemIndex} keys={keys} />;
  }

  _onChange(keyPath) {
    if (CollectionStore.pathToKey(this.props.path) === keyPath) {
      var item = CollectionStore.getCollection(this.props.path);

      var common = {};
      var fields = ["TotalTime", "Artist", "Composer", "ID", "Year"];
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
}

GroupContent.propTypes = {
  path: React.PropTypes.array.isRequired,
  itemIndex: React.PropTypes.number.isRequired,
};


class GroupList extends React.Component {
  render() {
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
      <div>
        {groups}
      </div>
    );
  }
}

GroupList.propTypes = {
  path: React.PropTypes.array.isRequired,
  itemIndex: React.PropTypes.number.isRequired,
  keys: React.PropTypes.array.isRequired,
};


class TrackList extends React.Component {
  render() {
    var list = this.props.list;
    var keys = this.props.keys;
    var tracks = keys.map(function(i) {
      return <Track key={list[i].ID} data={list[i]} path={this.props.path.concat([i])} itemIndex={this.props.itemIndex} index={i} />;
    }.bind(this));

    return (
      <div className="tracks">
        <ol className={this.props.listStyle}>
          {tracks}
        </ol>
      </div>
    );
  }
}

TrackList.propTypes = {
  itemIndex: React.PropTypes.number.isRequired,
  path: React.PropTypes.array.isRequired,
  list: React.PropTypes.array.isRequired,
  keys: React.PropTypes.array.isRequired,
  listStyle: React.PropTypes.string.isRequired,
};


function pathsEqual(p1, p2) {
  return CollectionStore.pathToKey(p1) === CollectionStore.pathToKey(p2);
}

function isCurrent(i, p) {
  var c = CursorStore.getCurrent();
  if (c === null) {
    return false;
  }
  return pathsEqual(c.path(), p) && (i === c.index());
}

function isPlaying(id) {
  var playing = NowPlayingStore.getPlaying();
  if (!playing) {
    return false;
  }

  var t = NowPlayingStore.getTrack();
  if (t) {
    return t.ID === id;
  }
  return false;
}

function getTrackState(props) {
  return {
    isCurrent: isCurrent(props.itemIndex, props.path),
    isPlaying: isPlaying(props.data.ID),
  };
}

class Track extends React.Component {
  constructor(props) {
    super(props);

    this.state = getTrackState(this.props);
    this._onClick = this._onClick.bind(this);
    this._onClickRemove = this._onClickRemove.bind(this);
    this._onChange = this._onChange.bind(this);
  }

  componentDidMount() {
    CursorStore.addChangeListener(this._onChange);
    NowPlayingStore.addChangeListener(this._onChange);
  }

  componentWillUnmount() {
    CursorStore.removeChangeListener(this._onChange);
    NowPlayingStore.removeChangeListener(this._onChange);
  }

  render() {
    var durationSecs = parseInt(this.props.data.TotalTime / 1000);
    var style = {
      current: this.state.isCurrent,
      "is-playing": this.state.isPlaying,
    };

    return (
      <li onClick={this._onClick} style={{"counterReset": "li " + (this.props.index + 1)}} className={classNames(style)}>
        <span id={"track_" + this.props.data.ID} className="name">{this.props.data.Name}</span>
        <span className="info">
          <Icon icon="clear"onClick={this._onClickRemove} />
          <TimeFormatter className="duration" time={durationSecs} />
        </span>
      </li>
    );
  }

  _onClick() {
    CursorActions.set(this.props.itemIndex, this.props.path);
  }

  _onClickRemove(e) {
    e.stopPropagation();
    PlaylistActions.remove(this.props.itemIndex, this.props.path);
  }

  _onChange() {
    this.setState(getTrackState(this.props));
  }
}

Track.propTypes = {
  itemIndex: React.PropTypes.number.isRequired,
  index: React.PropTypes.number.isRequired,
  path: React.PropTypes.array.isRequired,
};
