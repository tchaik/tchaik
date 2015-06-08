"use strict";

import React from "react";

import Icon from "./Icon.js";

import classNames from "classnames";

import PlayerKeyStore from "../stores/PlayerKeyStore.js";
import PlayerKeyActions from "../actions/PlayerKeyActions.js";


function getStatus() {
  return {
    set: PlayerKeyStore.isKeySet(),
  };
}

export default class PlayerKeyView extends React.Component {
  constructor(props) {
    super(props);

    this.state = getStatus();
    this._onChange = this._onChange.bind(this);
  }

  componentDidMount() {
    var k = PlayerKeyStore.getKey();
    if (k !== null) {
      PlayerKeyActions.setKey(k);
    }
    PlayerKeyStore.addChangeListener(this._onChange);
  }

  componentWillUnmount() {
    PlayerKeyStore.removeChangeListener(this._onChange);
  }

  render() {
    var classes = {
      "item": true,
      "key": true,
      "set": this.state.set,
    };
    var title = this.state.set ? "Player Key: Set" : "";
    return (
      <span className={classNames(classes)}>
        <Icon icon="transfer" title={title} />
      </span>
    );
  }

  _onChange() {
    this.setState(getStatus());
  }
}
