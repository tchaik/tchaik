"use strict";

import React from "react/addons";

import classNames from "classnames";

import Icon from "./Icon.js";
import StatusView from "./Status.js";
import PlayerKeyView from "./PlayerKey.js";
import MenuToggleButton from "./MenuToggleButton.js";

import ContainerActions from "../actions/ContainerActions.js";
import ContainerConstants from "../constants/ContainerConstants.js";
import ContainerStore from "../stores/ContainerStore.js";

import LeftColumnActions from "../actions/LeftColumnActions.js";
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
      <div className={classNames(classes)} onClick={this._onClick}>
        <Icon {...other} />
        <span className="title">{this.props.title}</span>
      </div>
    );
  }

  _onClick() {
    ContainerActions.mode(this.props.mode);
  }

  _onChange() {
    this.setState(getToolBarItemState(this.props.mode));
  }
}

class LinkItem extends React.Component {
  render() {
    return (
      <a className="menu-item" href={this.props.href} target="_blank">
        <span className="item">
          <Icon icon={this.props.icon} />
        </span>
        <span className="title">{this.props.title}</span>
      </a>
    );
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
    this._onClickMenuToggle = this._onClickMenuToggle.bind(this);
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
          <div className="menu-item menu-button-item" onClick={this._onClickMenuToggle}>
            <MenuToggleButton />
            <span className="title">Tchaik</span>
          </div>
          <div className="menu-items">
            <ToolbarItem mode={ContainerConstants.ALL} icon="align-justify" title="All" />
            <ToolbarItem mode={ContainerConstants.ARTISTS} icon="list" title="Artists" />
            <ToolbarItem mode={ContainerConstants.COVERS} icon="th-large" title="Covers" />
            <ToolbarItem mode={ContainerConstants.RECENT} icon="time" title="Recently Added" />
            <ToolbarItem mode={ContainerConstants.RETRO} icon="cd" title="Retro" />
            <ToolbarItem mode={ContainerConstants.SETTINGS} icon="cog" title="Settings" />
          </div>
          <div className="links">
            <LinkItem title="Github" href="https://github.com/tchaik/tchaik" icon="home" />
          </div>
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

  _onClickMenuToggle() {
    LeftColumnActions.toggleVisibility();
  }
}
