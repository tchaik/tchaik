"use strict";

import React from "react";

import PathListStore from "../stores/PathListStore.js";
import PathListActions from "../actions/PathListActions.js";

import {GroupList as GroupList} from "./Collection.js";


function getPathListState(name) {
  return {
    items: PathListStore.getPaths(name),
  };
}

export default class PathList extends React.Component {
  constructor(props) {
    super(props);

    this.state = getPathListState(this.props.name);
    this._onChange = this._onChange.bind(this);
  }

  componentDidMount() {
    PathListStore.addChangeListener(this._onChange);
    PathListActions.fetch(this.props.name);
  }

  componentWillUnmount() {
    PathListStore.removeChangeListener(this._onChange);
  }

  componentWillReceiveProps(nextProps) {
    this.setState(getPathListState(nextProps.name));
    PathListActions.fetch(nextProps.name);
  }

  render() {
    return <GroupList path={["Root"]} list={this.state.items} depth={0} />;
  }

  _onChange() {
    this.setState(getPathListState(this.props.name));
  }
}

PathList.propTypes = {
  name: React.PropTypes.string.isRequired,
};
