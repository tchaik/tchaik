/** @jsx React.DOM */
'use strict';

var React = require('react/addons');

var Icon = require('./Icon.js');

var classNames = require('classnames');

var PlayerKeyStore = require('../stores/PlayerKeyStore.js');
var PlayerKeyActions = require('../actions/PlayerKeyActions.js');

function getStatus() {
  return {
    set: PlayerKeyStore.isKeySet(),
  };
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
    var title = this.state.set ? "Player Key: Set" : "";
    return (
      <span className={classNames(classes)}>
        <Icon icon="transfer" title={title} />
      </span>
    );
  },

  _onChange: function() {
    this.setState(getStatus());
  }
});

module.exports = PlayerKeyView;
