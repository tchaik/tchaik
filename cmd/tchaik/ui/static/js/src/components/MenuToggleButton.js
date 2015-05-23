import React from "react/addons";

import Icon from "./Icon.js";

import LeftColumnActions from "../actions/LeftColumnActions.js";

export default class MenuButton extends React.Component {
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

  _onClick(evt) {
    evt.stopPropagation();

    LeftColumnActions.toggleVisibility();
  }
}
