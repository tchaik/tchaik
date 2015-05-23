"use strict";

import React from "react/addons";

import classNames from "classnames";

import Icon from "./Icon.js";
import StatusView from "./Status.js";
import PlayerKeyView from "./PlayerKey.js";

import {RootCollection as RootCollection} from "./Collection.js";
import {Search as Search} from "./Search.js";
import Covers from "./Covers.js";
import Filter from "./Filter.js";
import Recent from "./Recent.js";
import Settings from "./Settings.js";
import Retro from "./Retro.js";

import LeftColumnConstants from "../constants/LeftColumnConstants.js";
import LeftColumnStore from "../stores/LeftColumnStore.js";
import LeftColumnActions from "../actions/LeftColumnActions.js";

import SearchStore from "../stores/SearchStore.js";


function getToolBarItemState(mode) {
  return {selected: mode === LeftColumnStore.getMode()};
}

class ToolbarItem extends React.Component {
  constructor(props) {
    super(props);

    this.state = getToolBarItemState(this.props.mode);
    this._onClick = this._onClick.bind(this);
    this._onChange = this._onChange.bind(this);
  }

  componentDidMount() {
    LeftColumnStore.addChangeListener(this._onChange);
  }

  componentWillUnmount() {
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
    LeftColumnActions.mode(this.props.mode);
  }

  _onChange() {
    this.setState(getToolBarItemState(this.props.mode));
  }
}


function leftColumnState() {
  return {
    mode: LeftColumnStore.getMode(),
    hidden: LeftColumnStore.getIsHidden(),
  };
}

export default class LeftColumn extends React.Component {
  constructor(props) {
    super(props);

    this.state = leftColumnState();
    this._onChange = this._onChange.bind(this);
    this._onSearch = this._onSearch.bind(this);
  }

  componentDidMount() {
    LeftColumnStore.addChangeListener(this._onChange);
    SearchStore.addChangeListener(this._onSearch);
  }

  componentWillUnmount() {
    LeftColumnStore.removeChangeListener(this._onChange);
    SearchStore.removeChangeListener(this._onSearch);
  }

  render() {
    var container = null;
    switch (this.state.mode) {
    case LeftColumnConstants.ARTISTS:
      container = <Filter name="Artist" />;
      break;
    case LeftColumnConstants.SEARCH:
      container = <Search />;
      break;
    case LeftColumnConstants.COVERS:
      container = <Covers />;
      break;
    case LeftColumnConstants.RECENT:
      container = <Recent />;
      break;
    case LeftColumnConstants.SETTINGS:
      container = <Settings />;
      break;
    case LeftColumnConstants.RETRO:
      container = <Retro />;
      break;
    case LeftColumnConstants.ALL:
    default:
      container = <RootCollection />;
      break;
    }

    var containerClasses = classNames({
      container: true,
      retro: this.state.mode === LeftColumnConstants.RETRO,
    });

    var toolbar = null;
    if (!this.state.hidden) {
      toolbar = (
        <div className="control-bar">
          <ToolbarItem mode={LeftColumnConstants.ALL} icon="align-justify" title="All" />
          <ToolbarItem mode={LeftColumnConstants.ARTISTS} icon="list" title="Artists" />
          <ToolbarItem mode={LeftColumnConstants.COVERS} icon="th-large" title="Covers" />
          <ToolbarItem mode={LeftColumnConstants.RECENT} icon="time" title="Recently Added" />
          <ToolbarItem mode={LeftColumnConstants.RETRO} icon="cd" title="Reto" />
          <ToolbarItem mode={LeftColumnConstants.SETTINGS} icon="cog" title="Settings" />
          <div className="bottom">
            <StatusView />
            <PlayerKeyView />
          </div>
        </div>
      );
    }

    return (
      <div>
        {toolbar}
        <div id="container" className={containerClasses}>
          {container}
        </div>
      </div>
    );
  }

  _onChange() {
    this.setState(leftColumnState());
  }

  _onSearch() {
    this.setState({mode: LeftColumnConstants.SEARCH});
  }
}
