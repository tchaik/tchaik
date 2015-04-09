/** @jsx React.DOM */
'use strict';

var React = require('react/addons');

var Icon = require('./Icon.js');

var WebsocketApi = require('../api/WebsocketApi.js');
var WebsocketApiActions = require('../actions/WebsocketApiActions.js');

var classNames = require('classnames');

function getStatus() {
  return WebsocketApi.getStatus();
}

var StatusView = React.createClass({
  getInitialState: function() {
    return getStatus();
  },

  componentDidMount: function() {
    WebsocketApi.addChangeListener(this._onChange);
  },

  componentWillUnmount: function() {
    WebsocketApi.removeChangeListener(this._onChange);
  },

  render: function() {
    var classes = {
      'item': true,
      'status': true,
      'closed': !this.state.open
    };
    var title = this.state.open ? "Online" : "Offline";

    return (
      <span className={classNames(classes)} onClick={this._onClick}>
        <Icon icon="flash" title={title} />
      </span>
    );
  },

  _onChange: function() {
    this.setState(getStatus());
  },

  _onClick: function() {
    WebsocketApiActions.reconnect();
  }
});

module.exports = StatusView;
