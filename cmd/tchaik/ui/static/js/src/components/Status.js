/** @jsx React.DOM */
'use strict';

var React = require('react/addons');

var Icon = require('./Icon.js');

var WebsocketApi = require('../api/WebsocketApi.js');

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
      'open': this.state.open
    };
    var title = "Connection " + (this.state.open ? "open" : "closed");

    return (
      <span className={classNames(classes)}><Icon icon="flash" title={title} /></span>
    );
  },

  _onChange: function() {
    this.setState(getStatus());
  }
});

module.exports = StatusView;
