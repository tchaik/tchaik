/** @jsx React.DOM */
'use strict';

var React = require('react/addons');

var classNames = require('classnames');

var FilterStore = require('../stores/FilterStore.js');
var FilterActions = require('../actions/FilterActions.js');

var CollectionStore = require('../stores/CollectionStore.js');

var RootGroup = require('./Search.js').RootGroup;

var Filter = React.createClass({
  getInitialState: function() {
    return {
      items: [],
      current: null,
    };
  },

  componentDidMount: function() {
    FilterStore.addChangeListener(this._onChange);
    FilterActions.fetchList(this.props.name);
  },
  
  componentWillUnmount: function() {
    FilterStore.removeChangeListener(this._onChange);
  },

  setCurrent: function(itemName) {
    FilterActions.setItem(this.props.name, itemName);
    FilterActions.fetchPaths(this.props.name, itemName);
  },

  render: function() {
    var items = this.state.items.map(function(item) {
      return <Item item={item}
                   key={item}
               current={this.state.current === item}
            setCurrent={this.setCurrent} />;
    }.bind(this));

    return (
      <div className="filter">
        <div className="sidebar">
          <ul>{items}</ul>
        </div>
        <div className="collection">
          <Results filterName={this.props.name} itemName={this.state.current} />
        </div>
      </div>
    );
  },

  _onChange: function() {
    this.setState({
      items: FilterStore.getItems(this.props.name),
      current: FilterStore.getCurrentItem(this.props.name),
    });
  }
});

var Item = React.createClass({
  render: function() {
    return <li onClick={this._onClick} className={classNames({'selected':this.props.current})}>{this.props.item}</li>;
  },

  _onClick: function() {
    this.props.setCurrent(this.props.item);
  }
});

var Results = React.createClass({
  getInitialState: function() {
    return {items:[]};
  },

  componentDidMount: function() {
    FilterStore.addChangeListener(this._onChange);
    FilterActions.fetchPaths(this.props.filterName, this.props.itemName);
  },

  componentWillUnmount: function() {
    FilterStore.removeChangeListener(this._onChange);
  },

  render: function() {
    var list = this.state.items.map(function(path) {
      return <RootGroup path={path} key={CollectionStore.pathToKey(path)} />;
    });
    return (
      <div className="collection" key={this.props.itemName}>
        {list}
      </div>
    );
  },
  
  _onChange: function() {
    this.setState({
      items: FilterStore.getPaths(this.props.filterName, FilterStore.getCurrentItem(this.props.filterName)),
    });
  }
});

module.exports = Filter;