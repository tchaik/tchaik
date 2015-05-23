"use strict";

import React from "react/addons";

import Icon from "./Icon.js";
import MenuToggleButton from "./MenuToggleButton.js";
import SearchStore from "../stores/SearchStore.js";
import SearchActions from "../actions/SearchActions.js";

function getTopState() {
  return {
    searchValue: SearchStore.getInput(),
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
  }

  componentWillUnmount() {
    SearchStore.removeChangeListener(this._onChange);
  }

  render() {
    return (
      <div>
        <MenuToggleButton />
        <div className="search">
          <Icon icon="search" />
          <input type="text" name="search" placeholder="Search Music"
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
