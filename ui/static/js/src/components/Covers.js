/** @jsx React.DOM */
'use strict';

var React = require('react/addons');

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
        list: CollectionStore.getCollection(["Root"]).Groups,
      });
    }
  },
});

var Cover = React.createClass({
  getInitialState: function() {
    return {
      expanded: false,
    };
  },

  render: function() {
    return (
      <div className="cover">
        <img src={"/artwork/"+this.props.item.TrackID} />
      </div>
    );
  },
});

module.exports = Covers;
