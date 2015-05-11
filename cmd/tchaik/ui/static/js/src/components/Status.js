'use strict';

var React = require('react/addons');

var Icon = require('./Icon.js');

var WebsocketAPI = require('../utils/WebsocketAPI.js');
var WebsocketActions = require('../actions/WebsocketActions.js');

var classNames = require('classnames');

function getStatus() {
  return WebsocketAPI.getStatus();
}

var StatusView = React.createClass({
  getInitialState: function() {
    return getStatus();
  },

  componentDidMount: function() {
    WebsocketAPI.addChangeListener(this._onChange);
  },

  componentWillUnmount: function() {
    WebsocketAPI.removeChangeListener(this._onChange);
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
    WebsocketActions.reconnect();
  }
});

module.exports = StatusView;
