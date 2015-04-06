/** @jsx React.DOM */
'use strict';

var React = require('react/addons');

var classNames = require('classnames');

var Icon = require('./Icon.js');
var StatusView = require('./Status.js');
var ApiKeyView = require('./ApiKey.js');

var RootCollection = require('./Collection.js').RootCollection;
var Search = require('./Search.js');
var Covers = require('./Covers.js');

var LeftColumnStore = require('../stores/LeftColumnStore.js');
var LeftColumnActions = require('../actions/LeftColumnActions.js');


function getToolBarItemState(mode) {
  return {
    selected: mode === LeftColumnStore.getMode(),
  };
}

var ToolbarItem = React.createClass({
  getInitialState: function() {
    return getToolBarItemState(this.props.mode);
  },

  componentDidMount: function() {
    LeftColumnStore.addChangeListener(this._onChange);
  },

  componentWillUnmount: function() {
    LeftColumnStore.removeChangeListener(this._onChange);
  },

  render: function() {
    var {...other} = this.props;
    var classes = {
      item: true,
      toolbar: true,
      selected: this.state.selected,
    };
    return (
      <span className={classNames(classes)} onClick={this._onClick}>
        <Icon {...other} />
      </span>
    );
  },

  _onClick: function() {
    LeftColumnActions.mode(this.props.mode);
  },

  _onChange: function() {
    this.setState(getToolBarItemState(this.props.mode));
  }
});

function leftColumnState() {
  return {
    mode: LeftColumnStore.getMode(),
  };
}

var LeftColumn = React.createClass({
  getInitialState: function() {
    return leftColumnState();
  },

  componentDidMount: function() {
    LeftColumnStore.addChangeListener(this._onChange);
  },

  componentWillUnmount: function() {
    LeftColumnStore.removeChangeListener(this._onChange);
  },

  render: function() {
    var container = null;
    switch (this.state.mode) {
    case "All":
      container = <RootCollection />;
      break;
    case "Search":
      container = <Search />;
      break;
    case "Covers":
      container = <Covers />;
      break;
    }

    return (
      <div>
        <div id="header">
          <ToolbarItem mode="Search" icon="search" title="Search" />
          <ToolbarItem mode="All" icon="align-justify" title="All" />
          <ToolbarItem mode="Browse" icon="list" title="Albums" />
          <ToolbarItem mode="Covers" icon="th-large" title="Covers" />
          <div className="bottom">
            <StatusView />
            <ApiKeyView />
          </div>
        </div>
        <div id="container">
          {container}
        </div>
      </div>
    );
  },

  _onChange: function() {
    this.setState(leftColumnState());
  },
});


module.exports = LeftColumn;
