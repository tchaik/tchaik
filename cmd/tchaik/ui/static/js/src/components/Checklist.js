"use strict";

import React from "react";

import ChecklistStore from "../stores/ChecklistStore.js";
import ChecklistActions from "../actions/ChecklistActions.js";

import {GroupList as GroupList} from "./Collection.js";


function getChecklistState() {
  return {
    items: ChecklistStore.getPaths(),
  };
}

export default class Checklist extends React.Component {
  constructor(props) {
    super(props);

    this.state = getChecklistState();
    this._onChange = this._onChange.bind(this);
  }

  componentDidMount() {
    ChecklistStore.addChangeListener(this._onChange);
    ChecklistActions.fetch();
  }

  componentWillUnmount() {
    ChecklistStore.removeChangeListener(this._onChange);
  }

  render() {
    return <GroupList path={["Root"]} list={this.state.items} depth={0} />;
  }

  _onChange() {
    this.setState(getChecklistState());
  }
}
