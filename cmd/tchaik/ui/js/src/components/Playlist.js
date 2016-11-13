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
    let items = this.state.list;
    if (items.length === 0) {
      return (
        <div className="playlist">
          <div className="no-items"><Icon icon="queue_music" />Empty playlist</div>
        </div>
      );
    }

    const pathCount = {};
    items = items.map(function(item, i) {
      const path = item.path();
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
  const c = CollectionStore.getCollection(path);
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
    const groupClasses = {
      "group": true,
      "expanded": this.state.expanded,
    };

    let image = null;
    if (this.props.root && this.state.common.id) {
      image = <ArtworkImage path={`/artwork/${this.state.common.id}`} />;
    }

    let duration = null;
    if (this.state.common.totalTime) {
      duration = <TimeFormatter className="duration" time={parseInt(this.state.common.totalTime / 1000)} />;
    }

    let content = null;
    if (this.state.expanded) {
      content = <GroupContent path={this.props.path} setCommon={this.setCommon} itemIndex={this.props.itemIndex} />;
    }

    return (
      <div className={classNames(groupClasses)}>
        <div className="info-container" onClick={this._onClick}>
          {image}
          <div className="info">
            <div className="details">
              <div className="name">{this.props.item.name}</div>
              <GroupAttributes data={this.state.common} attributes={["artist", "composer", "year"]} />
              <div className="attributes duration"><Icon icon="schedule" />{duration}</div>
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
    const item = this.state.item;
    if (item === null) {
      return null;
    }

    const keys = PlaylistStore.getItemKeys(this.props.itemIndex, this.props.path);
    if (item.groups) {
      return <GroupList path={this.props.path} list={item.groups} itemIndex={this.props.itemIndex} keys={keys} />;
    }
    return <TrackList path={this.props.path} list={item.tracks} listStyle={item.listStyle} itemIndex={this.props.itemIndex} keys={keys} />;
  }

  _onChange(keyPath) {
    if (CollectionStore.pathToKey(this.props.path) === keyPath) {
      const item = CollectionStore.getCollection(this.props.path);

      const common = {};
      const fields = ["totalTime", "artist", "composer", "id", "year"];
      for (let f of fields) {
        if (item[f]) {
          common[f] = item[f];
        }
      }

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
    const itemByKey = {};
    for (const item of this.props.list) {
      itemByKey[item.key] = item;
    }

    const keys = this.props.keys;
    let groups = keys.map(function(key) {
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


function pathsEqual(p1, p2) {
  return CollectionStore.pathToKey(p1) === CollectionStore.pathToKey(p2);
}

function pathsPrefix(p, prefix) {
  const k = CollectionStore.pathToKey(p);
  if (k == null) {
    return false;
  }
  return k.indexOf(CollectionStore.pathToKey(prefix)) === 0;
}

function hasCurrent(i, p) {
  const c = CursorStore.getCurrent();
  if (c === null) {
    return false;
  }
  return (i === c.index()) && pathsPrefix(c.path(), p);
}

function getTrackListState(props) {
  return {
    hasCurrent: hasCurrent(props.itemIndex, props.path),
  };
}

class TrackList extends React.Component {
  constructor(props) {
    super(props);

    this.state = getTrackListState(this.props);
    this._onChange = this._onChange.bind(this);
  }

  componentDidMount() {
    CursorStore.addChangeListener(this._onChange);
  }

  componentWillUnmount() {
    CursorStore.removeChangeListener(this._onChange);
  }

  shouldComponentUpdate(nextProps, nextState) {
    if (this.state.hasCurrent || this.state.hasCurrent !== nextState.hasCurrent) {
      return true;
    }
    if (this.props.list.length !== nextProps.list.length) {
      return true;
    }
    return false;
  }

  render() {
    let currentPath = null;
    if (this.state.hasCurrent) {
      currentPath = CursorStore.getCurrent().path();
    }

    const list = this.props.list;
    let tracks = this.props.keys.map(function(i) {
      const path = this.props.path.concat([i]);
      let isCurrent = this.state.hasCurrent && pathsEqual(path, currentPath);

      return <Track key={list[i].id} data={list[i]} path={path} isCurrent={isCurrent} itemIndex={this.props.itemIndex} index={i} />;
    }.bind(this));

    return (
      <div className="tracks">
        <ol className={this.props.listStyle}>
          {tracks}
        </ol>
      </div>
    );
  }

  _onChange() {
    this.setState(getTrackListState(this.props));
  }
}

TrackList.propTypes = {
  itemIndex: React.PropTypes.number.isRequired,
  path: React.PropTypes.array.isRequired,
  list: React.PropTypes.array.isRequired,
  keys: React.PropTypes.array.isRequired,
  listStyle: React.PropTypes.string.isRequired,
};

function isPlaying(id) {
  if (!NowPlayingStore.getPlaying()) {
    return false;
  }

  const t = NowPlayingStore.getTrack();
  if (t) {
    return t.id === id;
  }
  return false;
}

function getTrackState(props) {
  return {
    isPlaying: isPlaying(props.data.id),
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
    NowPlayingStore.addChangeListener(this._onChange);
  }

  componentWillUnmount() {
    NowPlayingStore.removeChangeListener(this._onChange);
  }

  shouldComponentUpdate(nextProps, nextState) {
    if (nextState.playing !== this.state.isPlaying) {
      return true;
    }
    if (nextProps.isCurrent !== this.props.isCurrent) {
      return true;
    }
    return false;
  }

  render() {
    const durationSecs = parseInt(this.props.data.totalTime / 1000);
    const style = {
      current: this.props.isCurrent,
      "is-playing": this.state.isPlaying,
    };

    return (
      <li onClick={this._onClick} style={{"counterReset": "li " + (this.props.index + 1)}} className={classNames(style)}>
        <span id={"track_" + this.props.data.id} className="name">{this.props.data.name}</span>
        <span className="info">
          <Icon icon="clear" onClick={this._onClickRemove} />
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
