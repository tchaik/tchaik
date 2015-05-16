'use strict';

import React from 'react/addons';

import Icon from './Icon.js';

import SearchActions from '../actions/SearchActions.js';


export default class Top extends React.Component {
  constructor(props) {
    super(props);

    this._onChange = this._onChange.bind(this);
  }

  render() {
    return (
      <div>
        <Icon icon="search" />
        <input type="text" name="search" onChange={this._onChange} />
      </div>
    );
  }

  _onChange(e) {
    SearchActions.search(e.currentTarget.value);
  }
}
