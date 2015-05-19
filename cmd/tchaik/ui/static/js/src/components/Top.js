'use strict';

import React from 'react/addons';

import Icon from './Icon.js';

import LeftColumnActions from '../actions/LeftColumnActions.js';
import SearchActions from '../actions/SearchActions.js';


class MenuButton extends React.Component {
  constructor(props) {
    super(props);

    this._onClick = this._onClick.bind(this);
  }

  render() {
    return (
      <div className="menu-button" onClick={this._onClick}>
        <Icon icon="menu-hamburger"/>
      </div>
    );
  }

  _onClick() {
    LeftColumnActions.toggleVisibility();
  }
}

export default class Top extends React.Component {
  constructor(props) {
    super(props);

    this._onChange = this._onChange.bind(this);
    this._onClick = this._onClick.bind(this);
  }

  render() {
    return (
      <div>
        <MenuButton />
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
