/** @jsx React.DOM */
'use strict';

var React = require('react/addons');

var classNames = require('classnames');

var Icon = require('./Icon.js');

var WebsocketApi = require('../api/WebsocketApi.js');

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
    var cx = classNames;
    var classes = cx({
      'status': true,
      'open': this.state.open
    });
    return (
      <div className={classes}>
        <Icon icon="flash"/>
      </div>
    );
  },

  _onChange: function() {
    this.setState(getStatus());
  }
});

module.exports = StatusView;
