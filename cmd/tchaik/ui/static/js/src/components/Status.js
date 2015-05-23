"use strict";

import React from "react/addons";

import Icon from "./Icon.js";

import WebsocketAPI from "../utils/WebsocketAPI.js";
import WebsocketActions from "../actions/WebsocketActions.js";

import classNames from "classnames";


function getStatus() {
  return WebsocketAPI.getStatus();
}

export default class StatusView extends React.Component {
  constructor(props) {
    super(props);
    this.state = getStatus();

    this._onChange = this._onChange.bind(this);
  }

  componentDidMount() {
    WebsocketAPI.addChangeListener(this._onChange);
  }

  componentWillUnmount() {
    WebsocketAPI.removeChangeListener(this._onChange);
  }

  render() {
    var classes = {
      "item": true,
      "status": true,
      "closed": !this.state.open,
    };
    var title = this.state.open ? "Online" : "Offline";

    return (
      <span className={classNames(classes)} onClick={this._onClick}>
        <Icon icon="flash" title={title} />
      </span>
    );
  }

  _onChange() {
    this.setState(getStatus());
  }

  _onClick() {
    WebsocketActions.reconnect();
  }
}
