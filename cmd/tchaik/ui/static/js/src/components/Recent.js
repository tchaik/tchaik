'use strict';

var React = require('react/addons');

var RecentStore = require('../stores/RecentStore.js');
var RecentActions = require('../actions/RecentActions.js');

var CollectionStore = require('../stores/CollectionStore.js');

var RootGroup = require('./Search.js').RootGroup;

function getRecentState() {
  return {
    items: RecentStore.getPaths()
  };
}

var Recent = React.createClass({
  getInitialState: function() {
    return getRecentState();
  },

  componentDidMount: function() {
    RecentStore.addChangeListener(this._onChange);
    RecentActions.fetch();
  },

  componentWillUnmount: function() {
    RecentStore.removeChangeListener(this._onChange);
  },

  render: function() {
    var list = this.state.items.map(function(path) {
      return <RootGroup path={path} key={"rootgroup-"+CollectionStore.pathToKey(path)} />;
    });

    return (
      <div className="collection">
        {list}
      </div>
    );
  },

  _onChange: function() {
    this.setState(getRecentState());
  }
});

module.exports = Recent;
