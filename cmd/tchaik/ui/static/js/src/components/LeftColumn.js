"use strict";

import React from "react";

import classNames from "classnames";

import Icon from "./Icon.js";
import StatusView from "./Status.js";
import PlayerKeyView from "./PlayerKey.js";

import ContainerActions from "../actions/ContainerActions.js";
import ContainerConstants from "../constants/ContainerConstants.js";
import ContainerStore from "../stores/ContainerStore.js";

import LeftColumnStore from "../stores/LeftColumnStore.js";

function getToolBarItemState(mode) {
  return {selected: mode === ContainerStore.getMode()};
}

class ToolbarItem extends React.Component {
  constructor(props) {
    super(props);

    this.state = getToolBarItemState(this.props.mode);
    this._onClick = this._onClick.bind(this);
    this._onChange = this._onChange.bind(this);
  }

  componentDidMount() {
    ContainerStore.addChangeListener(this._onChange);
    LeftColumnStore.addChangeListener(this._onChange);
  }

  componentWillUnmount() {
    ContainerStore.removeChangeListener(this._onChange);
    LeftColumnStore.removeChangeListener(this._onChange);
  }

  render() {
    var {...other} = this.props;
    var classes = {
      "menu-item": true,
      selected: this.state.selected,
    };
    return (
      <li className={classNames(classes)} onClick={this._onClick}>
        <Icon {...other} />
        <span className="title">{this.props.title}</span>
      </li>
    );
  }

  _onClick() {
    ContainerActions.mode(this.props.mode);
  }

  _onChange() {
    this.setState(getToolBarItemState(this.props.mode));
  }
}

function leftColumnState() {
  return {
    mode: ContainerStore.getMode(),
    hidden: LeftColumnStore.getIsHidden(),
  };
}

export default class LeftColumn extends React.Component {
  constructor(props) {
    super(props);

    this.state = leftColumnState();
    this._onChange = this._onChange.bind(this);
  }

  componentDidMount() {
    ContainerStore.addChangeListener(this._onChange);
    LeftColumnStore.addChangeListener(this._onChange);
  }

  componentWillUnmount() {
    ContainerStore.removeChangeListener(this._onChange);
    LeftColumnStore.removeChangeListener(this._onChange);
  }

  render() {
    var classes = classNames("control-bar", {hidden: this.state.hidden});
    return (
      <div className={classes}>
        <div className="top">
          <ul className="menu">
            <ToolbarItem mode={ContainerConstants.ALL} icon="library_music" title="Library" />
            <ToolbarItem mode={ContainerConstants.ARTISTS} icon="group" title="Artists" />
            <ToolbarItem mode={ContainerConstants.COVERS} icon="view_module" title="Covers" />
            <ToolbarItem mode={ContainerConstants.RECENT} icon="schedule" title="Recently Added" />
            <ToolbarItem mode={ContainerConstants.FAVOURITE} icon="favorite" title="Favourite" />
            <ToolbarItem mode={ContainerConstants.CHECKLIST} icon="done_all" title="Checklist" />
            <ToolbarItem mode={ContainerConstants.RETRO} icon="album" title="Retro" />
            <ToolbarItem mode={ContainerConstants.SETTINGS} icon="settings" title="Settings" />
          </ul>
        </div>
        <div className="middle"></div>
        <div className="bottom">
          <div className="bottom-item"><StatusView /></div>
          <div className="bottom-item"><PlayerKeyView /></div>
        </div>
      </div>
    );
  }

  _onChange() {
    this.setState(leftColumnState());
  }
}
