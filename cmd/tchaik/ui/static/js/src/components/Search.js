/** @jsx React.DOM */
'use strict';

var React = require('react/addons');

var SearchStore = require('../stores/SearchStore.js');
var SearchActions = require('../actions/SearchActions.js');

var Group = require('./Collection.js').Group;
var CollectionStore = require('../stores/CollectionStore.js');
var CollectionActions = require('../actions/CollectionActions.js');

var Search = React.createClass({
  componentDidMount: function() {
    var input = SearchStore.getInput();
    if (input && input !== "") {
      SearchActions.search(input);
      React.findDOMNode(this.refs.input).value = input;
    }
  },

  render: function() {
    return (
      <div className="search-collection">
        <div id="search">
          <input ref="input" type="text" onChange={this._onChange} />
        </div>
        <div className="collection">
          <Results />
        </div>
      </div>
    );
  },

  _onChange: function(e) {
    SearchActions.search(e.currentTarget.value);
  },
});


var Results = React.createClass({
  getInitialState: function() {
    return {
      results: SearchStore.getResults(),
    };
  },

  componentDidMount: function() {
    SearchStore.addChangeListener(this._onChange);
  },

  componentWillUnmount: function() {
    SearchStore.removeChangeListener(this._onChange);
  },

  render: function() {
    var list = this.state.results.map(function(path) {
      return <RootGroup path={path} key={CollectionStore.pathToKey(path)} />;
    });
    return (
      <div className="collection">
        {list}
      </div>
    );
  },

  _onChange: function() {
    this.setState({
      results: SearchStore.getResults(),
    });
  }
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
      <Group item={this.state.item} path={this.props.path} depth={1} />
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

module.exports.Search = Search;
module.exports.RootGroup = RootGroup;
