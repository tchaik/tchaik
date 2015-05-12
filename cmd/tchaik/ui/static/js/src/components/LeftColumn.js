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

import LeftColumnStore from '../stores/LeftColumnStore.js';
import LeftColumnActions from '../actions/LeftColumnActions.js';

import VolumeStore from '../stores/VolumeStore.js';
import VolumeActions from '../actions/VolumeActions.js';


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


function _getOffsetTop(elem) {
  var offsetTop = 0;
  do {
    if (!isNaN(elem.offsetTop)) {
        offsetTop += elem.offsetTop;
    }
  } while ((elem = elem.offsetParent));
  return offsetTop;
}

function getVolumeStatus() {
  return {
    volume: VolumeStore.getVolume(),
  };
}

class Volume extends React.Component {
  constructor(props) {
    super(props);

    this.state = getVolumeStatus();
    this._onChange = this._onChange.bind(this);
    this._onWheel = this._onWheel.bind(this);
  }

  componentDidMount() {
    VolumeStore.addChangeListener(this._onChange);
  }

  componentWillUnmount() {
    VolumeStore.removeChangeListener(this._onChange);
  }

  render() {
    var volume = this.state.volume;
    var classSuffix;
    if (volume === 0.00) {
      classSuffix = 'off';
    } else if (volume < 0.5) {
      classSuffix = 'down';
    } else {
      classSuffix = 'up';
    }

    var h = this.props.height - parseInt(volume * this.props.height);
    return (
      <div className="volume" onWheel={this._onWheel} >
        <div className="bar" onClick={this._onClick} style={{height: this.props.height}}>
          <div className="current" style={{height: h}} />
          <div className="marker" style={{height: this.props.markerHeight}} />
        </div>
        <Icon icon={'volume-' + classSuffix} onClick={this._toggleMute} />
      </div>
    );
  }

  _toggleMute(evt) {
    evt.stopPropagation();
    VolumeActions.toggleVolumeMute();
  }

  _onWheel(evt) {
    evt.stopPropagation();
    var v = this.state.volume + 0.05 * evt.deltaY;
    if (v > 1.0) {
      v = 1.0;
    } else if (v < 0.00) {
      v = 0.0;
    }
    VolumeActions.volume(v);
  }

  _onClick(evt) {
    var pos = evt.pageY - _getOffsetTop(evt.currentTarget);
    var height = evt.currentTarget.offsetHeight;
    VolumeActions.volume(1 - pos/height);
  }

  _onChange() {
    this.setState(getVolumeStatus());
  }
}

Volume.propTypes = {
  height: React.PropTypes.number.isRequired,
  markerHeight: React.PropTypes.number.isRequired,
};
