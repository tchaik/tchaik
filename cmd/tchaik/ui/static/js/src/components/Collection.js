"use strict";

import React from "react";

import classNames from "classnames";

import Icon from "./Icon.js";
import TimeFormatter from "./TimeFormatter.js";
import GroupAttributes from "./GroupAttributes.js";
import ArtworkImage from "./ArtworkImage.js";

import CollectionStore from "../stores/CollectionStore.js";
import CollectionActions from "../actions/CollectionActions.js";

import NowPlayingStore from "../stores/NowPlayingStore.js";


export class RootCollection extends React.Component {
  render() {
    var path = ["Root"];
    if (this.props.path) {
      path = this.props.path;
    }
    return <GroupContent path={path} depth={0} />;
  }
}


export class Group extends React.Component {
  constructor(props) {
    super(props);

    this.state = {
      expanded: (this.props.depth !== 1),
      common: {},
    };

    this.setCommon = this.setCommon.bind(this);

    this._onClick = this._onClick.bind(this);
    this._onClickHeader = this._onClickHeader.bind(this);
    this._onClickImage = this._onClickImage.bind(this);
    this._onPlayNow = this._onPlayNow.bind(this);
    this._onQueue = this._onQueue.bind(this);
  }

  setCommon(c) {
    this.setState({common: c});
  }

  render() {
    var content = null;
    var play = null;
    var attributes = null;
    var headerDiv = null;
    var image = null;

    if (this.state.expanded) {
      content = [
        <GroupContent path={this.props.path} depth={this.props.depth} setCommon={this.setCommon} key="GroupContent0" />,
      ];

      if (this.props.depth === 1) {
        content.push(
          <div style={{clear: "both"}} key="GroupContent1" />
        );
      }

      var common = this.state.common;
      var duration = null;
      if (this.state.common.TotalTime) {
        duration = (
          <span>
            <Icon icon="time" /><TimeFormatter className="duration" time={parseInt(common.TotalTime / 1000)} />
          </span>
        );
      }

      var attributeArr = [];
      var fields = ["AlbumArtist", "Artist", "Composer", "Year"];
      fields.forEach(function(f) {
        if (common[f]) {
          attributeArr.push(common[f]);
        }
      });

      if (attributeArr.length > 0) {
        attributes = <GroupAttributes list={attributeArr} />;
      }

      play = (
        <span className="controls">
          {duration}
          <Icon icon="play" title="Play Now" onClick={this._onPlayNow} />
          <Icon icon="list" title="Queue" onClick={this._onQueue} />
        </span>
      );

      if (this.state.common.TrackID && this.props.depth === 1) {
        image = <ArtworkImage path={`/artwork/${common.TrackID}`} onClick={this._onClickImage}/>;
      }
    }

    // FIXME: this is to get around the "everything else" group which has no name
    // (and can"t be closed)
    var itemName = this.props.item.Name;
    if (this.props.depth === 1 && itemName === "") {
      itemName = "Misc";
    }

    var albumArtist = null;
    if (this.props.depth === 1 && (this.props.item.AlbumArtist || this.props.item.Artist)) {
      var artist = this.props.item.AlbumArtist || this.props.item.Artist;
      albumArtist = <span className="group-album-artist">{artist}</span>;
    }

    if (itemName !== "") {
      headerDiv = (
        <div className="header" onClick={this._onClickHeader}>
          {image}
          <span className="name">{itemName}</span>{albumArtist}
          {play}
          {attributes}
          <div style={{"clear": "both"}} />
        </div>
      );
    }

    var groupClasses = {
      "group": true,
      "expanded": this.state.expanded,
    };

    return (
      <div className={classNames(groupClasses)} onClick={this._onClick}>
        {headerDiv}
        {content}
      </div>
    );
  }

  _onClick(e) {
    if (!this.state.expanded) {
      this._onClickHeader(e);
    }
  }

  _onClickHeader(e) {
    e.stopPropagation();
    this.setState({expanded: !this.state.expanded});
  }

  _onClickImage(e) {
    e.stopPropagation();
  }

  _onPlayNow(e) {
    CollectionActions.playNow(this.props.path);
    e.stopPropagation();
  }

  _onQueue(e) {
    CollectionActions.appendToPlaylist(this.props.path);
    e.stopPropagation();
  }

}

Group.propTypes = {
  path: React.PropTypes.array.isRequired,
  item: React.PropTypes.object.isRequired,
  depth: React.PropTypes.number.isRequired,
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

    if (item.Groups) {
      return <GroupList path={this.props.path} depth={this.props.depth} list={item.Groups} />;
    }
    return <TrackList path={this.props.path} list={item.Tracks} listStyle={item.ListStyle} />;
  }

  _onChange(keyPath) {
    if (CollectionStore.pathToKey(this.props.path) === keyPath) {
      var item = CollectionStore.getCollection(this.props.path);

      var common = {};
      var fields = ["TotalTime", "AlbumArtist", "Artist", "TrackID", "Composer", "Year"];
      fields.forEach(function(f) {
        if (item[f]) {
          common[f] = item[f];
        }
      });

      if (Object.keys(common).length > 0 && this.props.setCommon) {
        this.props.setCommon(common);
      }
      this.setState({item: item});
    }
  }
}

GroupContent.propTypes = {
  path: React.PropTypes.array.isRequired,
  depth: React.PropTypes.number.isRequired,
};


class GroupList extends React.Component {
  render() {
    var list = this.props.list.map(function(item) {
      return <Group path={this.props.path.concat([item.Key])} depth={this.props.depth + 1} item={item} key={item.Key} />;
    }.bind(this));

    return (
      <div className="collection">
        {list}
      </div>
    );
  }
}

GroupList.propTypes = {
  path: React.PropTypes.array.isRequired,
  depth: React.PropTypes.number.isRequired,
  list: React.PropTypes.array.isRequired,
};


class TrackList extends React.Component {
  render() {
    var discs = {};
    var discIndices = [];
    var trackNumber = 0;
    this.props.list.forEach(function(track) {
      var discNumber = track.DiscNumber;
      if (!discs[discNumber]) {
        discs[discNumber] = [];
        discIndices.push(discNumber);
      }
      track.Key = String(trackNumber++);
      discs[discNumber].push(track);
    });

    var ols = [];
    var buildTrack = function(track) {
      return <Track key={track.TrackID} data={track} path={this.props.path.concat([track.Key])} />;
    }.bind(this);

    for (var i = 0; i < discIndices.length; i++) {
      var disc = discs[discIndices[i]];
      ols.push(
        <ol key={`disc${discIndices[i]}`} className={this.props.listStyle}>
          {disc.map(buildTrack)}
        </ol>
      );
    }

    return (
      <div className="tracks">
        {ols}
      </div>
    );
  }
}

TrackList.propTypes = {
  path: React.PropTypes.array.isRequired,
  list: React.PropTypes.array.isRequired,
  listStyle: React.PropTypes.string.isRequired,
};


function isCurrent(trackId) {
  var t = NowPlayingStore.getTrack();
  if (t) {
    return t.TrackID === trackId;
  }
  return false;
}


class Track extends React.Component {
  constructor(props) {
    super(props);

    this.state = {
      current: isCurrent(this.props.data.TrackID),
      playing: NowPlayingStore.getPlaying(),
      expanded: false,
    };
    this._onChange = this._onChange.bind(this);
    this._onClick = this._onClick.bind(this);
    this._onPlayNow = this._onPlayNow.bind(this);
    this._onMore = this._onMore.bind(this);
    this._onQueue = this._onQueue.bind(this);
  }

  componentDidMount() {
    NowPlayingStore.addChangeListener(this._onChange);
  }

  componentWillUnmount() {
    NowPlayingStore.removeChangeListener(this._onChange);
  }

  render() {
    var durationSecs = parseInt(this.props.data.TotalTime / 1000);
    var liClasses = {
      "current": this.state.current,
      "playing": this.state.current && this.state.playing,
    };

    var expanded = null;
    if (this.state.expanded) {
      var data = this.props.data;
      var attributeArr = [];
      var fields = ["AlbumArtist", "Artist", "Composer", "Year"];
      fields.forEach(function(f) {
        if (data[f]) {
          attributeArr.push(data[f]);
        }
      });

      var attributes = null;
      if (attributeArr.length > 0) {
        attributes = <GroupAttributes list={attributeArr} />;
      }
      expanded = <span className="expanded">{attributes}</span>;
    }

    return (
      <li className={classNames(liClasses)} onClick={this._onClick}>
        <span id={`track_${this.props.data.TrackID}`} className="name">{this.props.data.Name}</span>
        <TimeFormatter className="duration" time={durationSecs} />
        <span className="controls">
          <Icon icon="option-vertical" title="More" onClick={this._onMore} />
          <Icon icon="play" title="Play Now" onClick={this._onPlayNow} />
          <Icon icon="list" title="Queue" onClick={this._onQueue} />
        </span>
        {expanded}
      </li>
    );
  }

  _onChange() {
    this.setState({
      current: isCurrent(this.props.data.TrackID),
      playing: NowPlayingStore.getPlaying(),
    });
  }

  _onClick(e) {
    e.stopPropagation();
    CollectionActions.setCurrentTrack(this.props.data);
  }

  _onPlayNow(e) {
    e.stopPropagation();
    CollectionActions.playNow(this.props.path.concat([this.props.Key]));
  }

  _onQueue(e) {
    e.stopPropagation();
    CollectionActions.appendToPlaylist(this.props.path.concat([this.props.Key]));
  }

  _onMore(e) {
    e.stopPropagation();
    this.setState({
      expanded: !this.state.expanded,
    });
  }
}

Track.propTypes = {path: React.PropTypes.array.isRequired};
