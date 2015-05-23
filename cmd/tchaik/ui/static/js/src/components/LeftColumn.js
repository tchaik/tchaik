"use strict";

import React from "react/addons";

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
      item: true,
      toolbar: true,
      selected: this.state.selected,
    };
    return (
      <span className={classNames(classes)} onClick={this._onClick}>
        <Icon {...other} />
      </span>
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
    var toolbar = null;
    if (!this.state.hidden) {
      toolbar = (
        <div className="control-bar">
          <ToolbarItem mode={ContainerConstants.ALL} icon="align-justify" title="All" />
          <ToolbarItem mode={ContainerConstants.ARTISTS} icon="list" title="Artists" />
          <ToolbarItem mode={ContainerConstants.COVERS} icon="th-large" title="Covers" />
          <ToolbarItem mode={ContainerConstants.RECENT} icon="time" title="Recently Added" />
          <ToolbarItem mode={ContainerConstants.RETRO} icon="cd" title="Reto" />
          <ToolbarItem mode={ContainerConstants.SETTINGS} icon="cog" title="Settings" />
          <div className="bottom">
            <StatusView />
            <PlayerKeyView />
          </div>
        </div>
      );
    }

    return toolbar;
  }

  _onChange() {
    this.setState(leftColumnState());
  }
}
