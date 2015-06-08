"use strict";

import React from "react";

import Icon from "./Icon.js";
import MenuToggleButton from "./MenuToggleButton.js";

import SearchStore from "../stores/SearchStore.js";
import SearchActions from "../actions/SearchActions.js";

import ContainerStore from "../stores/ContainerStore.js";
import ContainerConstants from "../constants/ContainerConstants.js";

function getTopState() {
  return {
    searchValue: SearchStore.getInput(),
    title: ContainerStore.getTitle(),
  };
}

export default class Top extends React.Component {
  constructor(props) {
    super(props);

    this._onChange = this._onChange.bind(this);
    this._onInputChange = this._onInputChange.bind(this);
    this._onClick = this._onClick.bind(this);
    this.state = getTopState();
  }

  componentDidMount() {
    SearchStore.addChangeListener(this._onChange);
    ContainerStore.addChangeListener(this._onChange);
  }

  componentWillUnmount() {
    SearchStore.removeChangeListener(this._onChange);
    ContainerStore.removeChangeListener(this._onChange);
  }

  render() {
    return (
      <div className="top-container">
        <MenuToggleButton />
        <span className="title">{this.state.title}</span>
        <div className="search">
          <Icon icon="search" />
          <input type="text" name="search" placeholder="Search"
            value={this.state.searchValue}
            onChange={this._onInputChange} onClick={this._onClick} />
        </div>
      </div>
    );
  }

  _onChange() {
    var newState = getTopState();
    if (ContainerStore.getMode() !== ContainerConstants.SEARCH) {
      newState.searchValue = "";
    }
    this.setState(newState);
  }

  _onInputChange(e) {
    SearchActions.search(e.currentTarget.value);
  }

  _onClick(e) {
    e.currentTarget.select();
  }
}
