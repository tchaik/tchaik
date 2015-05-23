"use strict";

import React from "react/addons";

import Icon from "./Icon.js";
import MenuToggleButton from "./MenuToggleButton.js";

import SearchActions from "../actions/SearchActions.js";


export default class Top extends React.Component {
  constructor(props) {
    super(props);

    this._onChange = this._onChange.bind(this);
    this._onClick = this._onClick.bind(this);
  }

  render() {
    return (
      <div>
        <MenuToggleButton />
        <div className="search">
          <Icon icon="search" />
          <input type="text" name="search" placeholder="Search Music" onChange={this._onChange} onClick={this._onClick} />
        </div>
      </div>
    );
  }

  _onChange(e) {
    SearchActions.search(e.currentTarget.value);
  }

  _onClick(e) {
    e.currentTarget.select();
  }
}
