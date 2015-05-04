/** @jsx React.DOM */
'use strict';

var React = require('react/addons');

var Icon = require('./Icon.js');

var classNames = require('classnames');

var PlayerKeyStore = require('../stores/PlayerKeyStore.js');
var PlayerKeyActions = require('../actions/PlayerKeyActions.js');

function getStatus() {
  return {
    set: (PlayerKeyStore.isKeySet()),
  };
}

function randomString(len)
{
  var alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789";
  var res = "";
  for (var i = 0; i < len; i++) {
    res += alphabet.charAt(Math.floor(Math.random() * alphabet.length));
  }
  return res;
}

var PlayerKeyView = React.createClass({
  getInitialState: function() {
    return getStatus();
  },

  componentDidMount: function() {
    var k = PlayerKeyStore.getKey();
    if (k !== null) {
      PlayerKeyActions.setKey(k);
    }
    PlayerKeyStore.addChangeListener(this._onChange);
  },

  componentWillUnmount: function() {
    PlayerKeyStore.removeChangeListener(this._onChange);
  },

  render: function() {
    var classes = {
      'item': true,
      'key': true,
      'set': this.state.set,
    };
    var title = this.state.set ? "API: Enabled" : "API: Disabled";
    return (
      <span className={classNames(classes)}>
        <Icon icon="barcode" onClick={this._onClick} title={title} />
      </span>
    );
  },

  _onClick: function() {
    var key = PlayerKeyStore.getKey();
    if (key === null || key === "") {
      key = randomString(20);
    }
    key = prompt("Enter an API key", key);
    if (key !== null) {
      PlayerKeyActions.setKey(key);
    }
  },

  _onChange: function() {
    this.setState(getStatus());
  }
});

module.exports = PlayerKeyView;
