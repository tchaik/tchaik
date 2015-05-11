'use strict';

var React = require('react/addons');

var ArtworkImage = require('./ArtworkImage.js');
var Icon = require('./Icon.js');

var CollectionStore = require('../stores/CollectionStore.js');
var CollectionActions = require('../actions/CollectionActions.js');

var Covers = React.createClass({
  getInitialState: function() {
    return {
      list: [],
    };
  },

  componentDidMount: function() {
    CollectionStore.addChangeListener(this._onChange);
    CollectionActions.fetch(["Root"]);
  },

  componentWillUnmount: function() {
    CollectionStore.removeChangeListener(this._onChange);
  },

  render: function() {
    var covers = this.state.list.map(function(item) {
      return <Cover path={["Root"].concat(item.Key)} key={item.Key} item={item} />;
    });

    return (
      <div className="covers">
        {covers}
      </div>
    );
  },

  _onChange: function(path) {
    if (path === CollectionStore.pathToKey(["Root"])) {
      this.setState({
        list: CollectionStore.getCollection(["Root"]).Groups.slice(0, 30),
      });
    }
  },
});

var Cover = React.createClass({
  componentDidMount: function() {
    CollectionStore.addChangeListener(this._onChange);
    CollectionActions.fetch(this.props.path);
  },

  componentWillUnmount: function() {
    CollectionStore.removeChangeListener(this._onChange);
  },

  getInitialState: function() {
    return {
      item: this.props.item,
    };
  },

  render: function() {
    return (
      <div className="cover">
        <ArtworkImage path={"/artwork/"+this.state.item.TrackID} />
          <span className="controls">
            <Icon icon="play" title="Play Now" onClick={this._onPlayNow} />
            <Icon icon="list" title="Queue" onClick={this._onQueue} />
          </span>
      </div>
    );
  },

  _onChange: function(keyPath) {
    if (keyPath == CollectionStore.pathToKey(this.props.path)) {
      var item = CollectionStore.getCollection(this.props.path);
      this.setState({
        item: item,
      });
    }
  },

  _onPlayNow: function(e) {
    e.stopPropagation();
    CollectionActions.playNow(this.props.path);
  },

  _onQueue: function(e) {
    e.stopPropagation();
    CollectionActions.appendToPlaylist(this.props.path);
  },
});

module.exports = Covers;
