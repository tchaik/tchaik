/** @jsx React.DOM */
'use strict';

var React = require('react/addons');

var Icon = require('./Icon.js');

var classNames = require('classnames');

var CtrlKeyStore = require('../stores/CtrlKeyStore.js');
var CtrlKeyActions = require('../actions/CtrlKeyActions.js');

function getStatus() {
  return {
    set: (CtrlKeyStore.isKeySet()),
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

var CtrlKeyView = React.createClass({
  getInitialState: function() {
    return getStatus();
  },

  componentDidMount: function() {
    var k = CtrlKeyStore.getKey();
    if (k !== null) {
      CtrlKeyActions.setKey(k);
    }
    CtrlKeyStore.addChangeListener(this._onChange);
  },

  componentWillUnmount: function() {
    CtrlKeyStore.removeChangeListener(this._onChange);
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
    var key = CtrlKeyStore.getKey();
    if (key === null || key === "") {
      key = randomString(20);
    }
    key = prompt("Enter an API key", key);
    if (key !== null) {
      CtrlKeyActions.setKey(key);
    }
  },

  _onChange: function() {
    this.setState(getStatus());
  }
});

module.exports = CtrlKeyView;
