"use strict";

import React from "react";

import classNames from "classnames";

import Icon from "./Icon.js";
import TimeFormatter from "./TimeFormatter.js";
import GroupAttributes from "./GroupAttributes.js";
import ArtworkImage from "./ArtworkImage.js";

import CollectionStore from "../stores/CollectionStore.js";
import CollectionActions from "../actions/CollectionActions.js";

import PlaylistActions from "../actions/PlaylistActions.js";

import NowPlayingStore from "../stores/NowPlayingStore.js";


export class RootCollection extends React.Component {
  render() {
    let path = ["Root"];
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
      favourite: false,
      checklist: false,
      common: {},
    };

    this.setCommon = this.setCommon.bind(this);
    this.setFavourite = this.setFavourite.bind(this);
    this.setChecklist = this.setChecklist.bind(this);

    this._onClick = this._onClick.bind(this);
    this._onClickHeader = this._onClickHeader.bind(this);
    this._onClickImage = this._onClickImage.bind(this);
    this._onPlayNow = this._onPlayNow.bind(this);
    this._onQueue = this._onQueue.bind(this);
    this._toggleFavourite = this._toggleFavourite.bind(this);
    this._toggleChecklist = this._toggleChecklist.bind(this);
  }

  setCommon(c) {
    this.setState({common: c});
  }

  setFavourite(v) {
    this.setState({favourite: v});
  }

  setChecklist(v) {
    this.setState({checklist: v});
  }

  render() {
    let content = null;
    let play = null;
    let attributes = null;
    let headerDiv = null;
    let image = null;

    if (this.state.expanded || this.props.depth > 1) {
      content = [
        <GroupContent path={this.props.path} depth={this.props.depth} setCommon={this.setCommon} setFavourite={this.setFavourite} setChecklist={this.setChecklist} key="GroupContent0" />,
      ];

      if (this.props.depth === 1) {
        content.push(
          <div style={{clear: "both"}} key="GroupContent1" />
        );
      }

      const common = this.state.common;
      let duration = null;
      if (common.totalTime) {
        duration = (
          <span>
            <Icon icon="schedule" extraClasses={{duration: true}}/>
            <TimeFormatter className="time" time={parseInt(common.totalTime / 1000)} />
          </span>
        );
      }

      attributes = <GroupAttributes data={common} attributes={["albumArtist", "artist", "composer", "year"]} />;

      const favouriteIcon = this.state.favourite ? "favorite" : "favorite_border";
      const checklistIcon = this.state.checklist ? "check_circle" : "check";
      const checklistTitle = this.state.checklist ? "Remove from Checklist" : "Add to Checklist";

      play = (
        <span className="controls">
          <Kind kind={this.state.common.kind} />
          {duration}
          <Icon icon="play_arrow"title="Play Now" onClick={this._onPlayNow} />
          <Icon icon="playlist_add"title="Queue" onClick={this._onQueue} />
          <Icon icon={favouriteIcon} title="Favourite" extraClasses={{enabled: this.state.favourite}} onClick={this._toggleFavourite} />
          <Icon icon={checklistIcon} title={checklistTitle} extraClasses={{enabled: this.state.checklist}} onClick={this._toggleChecklist} />
        </span>
      );

      if (this.state.common.id && this.props.depth === 1) {
        image = <ArtworkImage path={`/artwork/${common.id}`} onClick={this._onClickImage}/>;
      }
    }

    // FIXME: this is to get around the "everything else" group which has no name
    // (and can"t be closed)
    let itemName = this.props.item.name;
    if (this.props.depth === 1 && itemName === "") {
      itemName = "Misc";
    }

    let albumArtist = null;
    if (this.props.depth === 1 && (this.props.item.albumArtist || this.props.item.artist)) {
      const artist = this.props.item.albumArtist || this.props.item.artist;
      albumArtist = <span className="album-artist">{artist.join(", ")}</span>;
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

    const groupClasses = {
      "group": true,
      "expanded": this.state.expanded || this.props.depth > 1,
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
    PlaylistActions.addItemPlayNow(this.props.path);
    e.stopPropagation();
  }

  _onQueue(e) {
    PlaylistActions.addItem(this.props.path);
    e.stopPropagation();
  }

  _toggleFavourite(e) {
    CollectionActions.setFavourite(this.props.path, !this.state.favourite);
    e.stopPropagation();
    this.setState({favourite: !this.state.favourite});
  }

  _toggleChecklist(e) {
    CollectionActions.setChecklist(this.props.path, !this.state.checklist);
    e.stopPropagation();
    this.setState({checklist: !this.state.checklist});
  }

}

class FavouriteIcon extends React.Component {

}

Group.propTypes = {
  path: React.PropTypes.array.isRequired,
  item: React.PropTypes.object.isRequired,
  depth: React.PropTypes.number.isRequired,
};

class Kind extends React.Component {
  constructor(props) {
    super(props);
  }

  render() {
    if (this.props.kind === "Apple Lossless audio file") {
      return <span className="alac">ALAC</span>;
    }
    return null;
  }
}

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

    if (item.groups) {
      return <GroupList path={this.props.path} depth={this.props.depth} list={item.groups} />;
    }
    return <TrackList path={this.props.path} list={item.tracks} listStyle={item.listStyle} />;
  }

  _onChange(keyPath) {
    if (CollectionStore.pathToKey(this.props.path) === keyPath) {
      const item = CollectionStore.getCollection(this.props.path);

      const common = {};
      for (const f of ["totalTime", "albumArtist", "artist", "id", "composer", "year", "kind"]) {
        if (item[f]) {
          common[f] = item[f];
        }
      }

      if (Object.keys(common).length > 0 && this.props.setCommon) {
        this.props.setCommon(common);
      }
      if (this.props.setFavourite) {
        this.props.setFavourite(item.favourite);
      }
      if (this.props.setChecklist) {
        this.props.setChecklist(item.checklist);
      }
      this.setState({item: item});
    }
  }
}

GroupContent.propTypes = {
  path: React.PropTypes.array.isRequired,
  depth: React.PropTypes.number.isRequired,
};


export class GroupList extends React.Component {
  render() {
    const list = this.props.list.map(function(item) {
      return <Group path={this.props.path.concat([item.key])} depth={this.props.depth + 1} item={item} key={item.key} />;
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
    const discs = {};
    const discIndices = [];
    let trackNumber = 0;
    this.props.list.forEach(function(track) {
      const discNumber = track.discNumber;
      if (!discs[discNumber]) {
        discs[discNumber] = [];
        discIndices.push(discNumber);
      }
      track.key = String(trackNumber++);
      discs[discNumber].push(track);
    });

    const ols = [];
    const buildTrack = function(track) {
      return <Track key={track.id} data={track} path={this.props.path.concat([track.key])} />;
    }.bind(this);

    for (const discIndex of discIndices) {
      const disc = discs[discIndex];
      ols.push(
        <ol key={`disc${discIndex}`} className={this.props.listStyle}>
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


function isCurrent(id) {
  const t = NowPlayingStore.getTrack();
  if (t) {
    return t.id === id;
  }
  return false;
}


class Track extends React.Component {
  constructor(props) {
    super(props);

    this.state = {
      current: isCurrent(this.props.data.id),
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
    const durationSecs = parseInt(this.props.data.totalTime / 1000);
    const liClasses = {
      "current": this.state.current,
      "playing": this.state.current && this.state.playing,
    };

    let expanded = null;
    if (this.state.expanded) {
      expanded = (
        <span className="expanded">
          <GroupAttributes data={this.props.data} attributes={["albumArtist", "artist", "composer", "year"]} />
        </span>
      );
    }

    return (
      <li className={classNames(liClasses)} onClick={this._onClick}>
        <span id={`track_${this.props.data.id}`} className="name">{this.props.data.name}</span>
        <TimeFormatter className="duration" time={durationSecs} />
        <span className="controls">
          <Icon icon="more_vert"title="More" onClick={this._onMore} />
          <Icon icon="play_arrow"title="Play Now" onClick={this._onPlayNow} />
          <Icon icon="playlist_add"title="Queue" onClick={this._onQueue} />
        </span>
        {expanded}
      </li>
    );
  }

  _onChange() {
    this.setState({
      current: isCurrent(this.props.data.id),
      playing: NowPlayingStore.getPlaying(),
    });
  }

  _onClick(e) {
    e.stopPropagation();
    CollectionActions.setCurrentTrack(this.props.data);
  }

  _onPlayNow(e) {
    e.stopPropagation();
    PlaylistActions.addItemPlayNow(this.props.path);
  }

  _onQueue(e) {
    e.stopPropagation();
    PlaylistActions.addItem(this.props.path);
  }

  _onMore(e) {
    e.stopPropagation();
    this.setState({
      expanded: !this.state.expanded,
    });
  }
}

Track.propTypes = {path: React.PropTypes.array.isRequired};
