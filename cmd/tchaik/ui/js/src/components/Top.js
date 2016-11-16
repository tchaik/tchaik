"use strict";

import React from "react";

import Icon from "./Icon.js";

import SearchStore from "../stores/SearchStore.js";
import SearchActions from "../actions/SearchActions.js";

import LeftColumnActions from "../actions/LeftColumnActions.js";

import ContainerStore from "../stores/ContainerStore.js";

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

    if (this.state.searchValue !== "") {
      SearchActions.search(this.state.searchValue);
    }
  }

  componentWillUnmount() {
    SearchStore.removeChangeListener(this._onChange);
    ContainerStore.removeChangeListener(this._onChange);
  }

  render() {
    return (
      <div className="container">
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
    this.setState(getTopState());
  }

  _onInputChange(e) {
    SearchActions.search(e.currentTarget.value);
  }

  _onClick(e) {
    e.currentTarget.select();
  }
}

const MenuToggleButton = () => {
  const onClick = (evt) => {
    evt.stopPropagation();
    LeftColumnActions.toggle();
  }

  return (
    <div className="menu-button" onClick={onClick}>
      <Icon icon="menu" />
    </div>
  );
}
