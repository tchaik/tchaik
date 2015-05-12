'use strict';

import React from 'react/addons';

import ArtworkImage from './ArtworkImage.js';
import Icon from './Icon.js';

import CollectionStore from '../stores/CollectionStore.js';
import CollectionActions from '../actions/CollectionActions.js';


export default class Covers extends React.Component {
  constructor(props) {
    super(props);

    this.state = {list: []};
    this._onChange = this._onChange.bind(this);
  }

  componentDidMount() {
    CollectionStore.addChangeListener(this._onChange);
    CollectionActions.fetch(["Root"]);
  }

  componentWillUnmount() {
    CollectionStore.removeChangeListener(this._onChange);
  }

  render() {
    var covers = this.state.list.map(function(item) {
      return <Cover path={["Root"].concat(item.Key)} key={item.Key} item={item} />;
    });

    return (
      <div className="covers">
        {covers}
      </div>
    );
  }

  _onChange(path) {
    if (path === CollectionStore.pathToKey(["Root"])) {
      this.setState({
        list: CollectionStore.getCollection(["Root"]).Groups.slice(0, 30),
      });
    }
  }
}


class Cover extends React.Component {
  constructor(props) {
    super(props);

    this.state = {item: this.props.item};
    this._onChange = this._onChange.bind(this);
    this._onPlayNow = this._onPlayNow.bind(this);
    this._onQueue = this._onQueue.bind(this);
  }

  componentDidMount() {
    CollectionStore.addChangeListener(this._onChange);
    CollectionActions.fetch(this.props.path);
  }

  componentWillUnmount() {
    CollectionStore.removeChangeListener(this._onChange);
  }

  render() {
    return (
      <div className="cover">
        <ArtworkImage path={"/artwork/"+this.state.item.TrackID} />
        <span className="controls">
          <Icon icon="play" title="Play Now" onClick={this._onPlayNow} />
          <Icon icon="list" title="Queue" onClick={this._onQueue} />
        </span>
      </div>
    );
  }

  _onChange(keyPath) {
    if (keyPath == CollectionStore.pathToKey(this.props.path)) {
      this.setState({item: CollectionStore.getCollection(this.props.path)});
    }
  }

  _onPlayNow(e) {
    e.stopPropagation();
    CollectionActions.playNow(this.props.path);
  }

  _onQueue(e) {
    e.stopPropagation();
    CollectionActions.appendToPlaylist(this.props.path);
  }
}
