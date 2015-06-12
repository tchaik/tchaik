"use strict";

import React from "react";

import RecentStore from "../stores/RecentStore.js";
import RecentActions from "../actions/RecentActions.js";

import {GroupList as GroupList} from "./Collection.js";


function getRecentState() {
  return {
    items: RecentStore.getPaths(),
  };
}

export default class Recent extends React.Component {
  constructor(props) {
    super(props);

    this.state = getRecentState();
    this._onChange = this._onChange.bind(this);
  }

  componentDidMount() {
    RecentStore.addChangeListener(this._onChange);
    RecentActions.fetch();
  }

  componentWillUnmount() {
    RecentStore.removeChangeListener(this._onChange);
  }

  render() {
    return <GroupList path={["Root"]} list={this.state.items} depth={0} />;
  }

  _onChange() {
    this.setState(getRecentState());
  }
}
