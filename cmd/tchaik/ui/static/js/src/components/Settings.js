/** @jsx React.DOM */
'use strict';

var React = require('react/addons');

var PlayerKeyStore = require('../stores/PlayerKeyStore.js');
var PlayerKeyActions = require('../actions/PlayerKeyActions.js');

function getPlayerKeySettingsState() {
  return {
    set: PlayerKeyStore.isKeySet(),
    key: PlayerKeyStore.getKey()
  };
}


var Settings = React.createClass({
  getInitialState: function() {
    return getPlayerKeySettingsState();
  },

  render: function() {
    return (
      <div className="settings">
        <PlayerKeyForm />
        <PushToPlayerKeyForm />
      </div>
    );
  }
});

function randomString(len)
{
  var alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789";
  var res = "";
  for (var i = 0; i < len; i++) {
    res += alphabet.charAt(Math.floor(Math.random() * alphabet.length));
  }
  return res;
}

function getPlayerKeyFormState() {
  return {
    set: PlayerKeyStore.isKeySet(),
    key: PlayerKeyStore.getKey(),
  };
}

var PlayerKeyForm = React.createClass({
  getInitialState: function() {
    return getPlayerKeyFormState();
  },

  componentDidMount: function() {
    PlayerKeyStore.addChangeListener(this._onChange);
  },

  componentWillUnmount: function() {
    PlayerKeyStore.removeChangeListener(this._onChange);
  },

  render: function() {
    var form = (
      <form onSubmit={this._onSetSubmit}>
        <input type="text" placeholder="Enter a player key" size="50" value={this.state.key} onChange={this._handleChange} />
        <button className={(!this.state.key) ? "" : "set"}>Set</button> or <button onClick={this._onGenerate}>Generate</button>
      </form>
    );

    if (this.state.set) {
      form = (
        <form onSubmit={this._onResetSubmit}>
          <input type="text" size="50" value={this.state.key} disabled />
          <button className="reset">Reset</button>
        </form>
      );
    }

    return (
      <div>
        <h3>Player Key</h3>
        <span>A <strong>Player Key</strong> uniquely identifies each music player and can be used to control it.</span>
        {form}
      </div>
    );
  },

  _onGenerate: function(e) {
    e.preventDefault();
    PlayerKeyActions.setKey(randomString(30));
  },

  _onChange: function() {
    this.setState(getPlayerKeyFormState());
  },

  _handleChange: function(event) {
     this.setState({key: event.target.value});
  },

  _onResetSubmit: function(e) {
    e.preventDefault();
    PlayerKeyActions.setKey("");
  },

  _onSetSubmit: function(e) {
    e.preventDefault();
    PlayerKeyActions.setKey(this.state.key);
  }
});

function getPushKeyFormState() {
  return {
    key: PlayerKeyStore.getPushKey(),
    set: PlayerKeyStore.isPushKeySet(),
  };
}

var PushToPlayerKeyForm = React.createClass({
  getInitialState: function() {
    return getPushKeyFormState();
  },

  componentDidMount: function() {
    PlayerKeyStore.addChangeListener(this._onChange);
  },

  componentWillUnmount: function() {
    PlayerKeyStore.removeChangeListener(this._onChange);
  },

  render: function() {
    var form = (
      <form onSubmit={this._onSetSubmit}>
        <input type="text" placeholder="Enter a player key" size="50" value={this.state.key} onChange={this._handleChange} />
        <button className={(!this.state.key) ? "" : "set"}>Set</button>
      </form>
    );

    if (this.state.set) {
      form = (
        <form onSubmit={this._onResetSubmit}>
          <input type="text" size="50" value={this.state.key} disabled />
          <button className="reset">Reset</button>
        </form>
      );
    }

    return (
      <div>
        <h3>Push Commands to a Player Key</h3>
        <span>Set a <strong>Player Key</strong> to push all UI commands to (instead of the UI player).</span>
        {form}
      </div>
    );
  },

  _onChange: function() {
    this.setState(getPushKeyFormState());
  },

  _handleChange: function(event) {
     this.setState({key: event.target.value});
  },

  _onResetSubmit: function(e) {
    e.preventDefault();
    PlayerKeyActions.setPushKey("");
  },

  _onSetSubmit: function(e) {
    e.preventDefault();
    PlayerKeyActions.setPushKey(this.state.key);
  }
});

module.exports = Settings;