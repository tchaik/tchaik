'use strict';

import React from 'react/addons';

import classNames from 'classnames';

import Icon from './Icon.js';
import StatusView from './Status.js';
import PlayerKeyView from './PlayerKey.js';

import {RootCollection as RootCollection} from './Collection.js';
import {Search as Search} from './Search.js';
import Covers from './Covers.js';
import Filter from './Filter.js';
import Recent from './Recent.js';
import Settings from './Settings.js';
import Volume from './Volume.js';

import LeftColumnStore from '../stores/LeftColumnStore.js';
import LeftColumnActions from '../actions/LeftColumnActions.js';


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
      selected: this.state.selected
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
  return {mode: LeftColumnStore.getMode()};
}

export default class LeftColumn extends React.Component {
  constructor(props) {
    super(props);

    this.state = leftColumnState();
    this._onChange = this._onChange.bind(this);
  }

  componentDidMount() {
    LeftColumnStore.addChangeListener(this._onChange);
  }

  componentWillUnmount() {
    LeftColumnStore.removeChangeListener(this._onChange);
  }

  render() {
    var container = null;
    switch (this.state.mode) {
    case "All":
      container = <RootCollection />;
      break;
    case "Artists":
      container = <Filter name="Artist" />;
      break;
    case "Search":
      container = <Search />;
      break;
    case "Covers":
      container = <Covers />;
      break;
    case "Recent":
      container = <Recent />;
      break;
    case "Settings":
      container = <Settings />;
    }

    return (
      <div>
        <div id="header">
          <ToolbarItem mode="Search" icon="search" title="Search" />
          <ToolbarItem mode="All" icon="align-justify" title="All" />
          <ToolbarItem mode="Artists" icon="list" title="Artists" />
          <ToolbarItem mode="Covers" icon="th-large" title="Covers" />
          <ToolbarItem mode="Recent" icon="time" title="Recently Added" />
          <ToolbarItem mode="Settings" icon="cog" title="Settings" />
          <div className="bottom">
            <Volume height={100} markerHeight={10} />
            <StatusView />
            <PlayerKeyView />
          </div>
        </div>
        <div id="container">
          {container}
        </div>
      </div>
    );
  }

  _onChange() {
    this.setState(leftColumnState());
  }
}
