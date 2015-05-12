'use strict';

import React from 'react/addons';

import PlayerKeyStore from '../stores/PlayerKeyStore.js';
import PlayerKeyActions from '../actions/PlayerKeyActions.js';


function getPlayerKeySettingsState() {
  return {
    set: PlayerKeyStore.isKeySet(),
    key: PlayerKeyStore.getKey()
  };
}

export default class Settings extends React.Component {
  constructor(props) {
    super(props);

    this.state = getPlayerKeySettingsState();
  }

  render() {
    return (
      <div className="settings">
        <PlayerKeyForm />
        <PushToPlayerKeyForm />
      </div>
    );
  }
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

function getPlayerKeyFormState() {
  return {
    set: PlayerKeyStore.isKeySet(),
    key: PlayerKeyStore.getKey(),
  };
}

class PlayerKeyForm extends React.Component {
  constructor(props) {
    super(props);

    this.state = getPlayerKeyFormState();
    this._onChange = this._onChange.bind(this);
    this._handleChange = this._handleChange.bind(this);
    this._onSetSubmit = this._onSetSubmit.bind(this);
  }

  componentDidMount() {
    PlayerKeyStore.addChangeListener(this._onChange);
  }

  componentWillUnmount() {
    PlayerKeyStore.removeChangeListener(this._onChange);
  }

  render() {
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
  }

  _onGenerate(e) {
    e.preventDefault();
    PlayerKeyActions.setKey(randomString(30));
  }

  _onChange() {
    this.setState(getPlayerKeyFormState());
  }

  _handleChange(event) {
     this.setState({key: event.target.value});
  }

  _onResetSubmit(e) {
    e.preventDefault();
    PlayerKeyActions.setKey("");
  }

  _onSetSubmit(e) {
    e.preventDefault();
    PlayerKeyActions.setKey(this.state.key);
  }
}

function getPushKeyFormState() {
  return {
    key: PlayerKeyStore.getPushKey(),
    set: PlayerKeyStore.isPushKeySet(),
  };
}

class PushToPlayerKeyForm extends React.Component {
  constructor(props) {
    super(props);

    this.state = getPushKeyFormState();
    this._onChange = this._onChange.bind(this);
    this._handleChange = this._handleChange.bind(this);
    this._onSetSubmit = this._onSetSubmit.bind(this);
  }

  componentDidMount() {
    PlayerKeyStore.addChangeListener(this._onChange);
  }

  componentWillUnmount() {
    PlayerKeyStore.removeChangeListener(this._onChange);
  }

  render() {
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
  }

  _onChange() {
    this.setState(getPushKeyFormState());
  }

  _handleChange(event) {
     this.setState({key: event.target.value});
  }

  _onResetSubmit(e) {
    e.preventDefault();
    PlayerKeyActions.setPushKey("");
  }

  _onSetSubmit(e) {
    e.preventDefault();
    PlayerKeyActions.setPushKey(this.state.key);
  }
}
