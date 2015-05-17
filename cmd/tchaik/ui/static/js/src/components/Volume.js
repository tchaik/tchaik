'use strict';

import React from 'react/addons';

import Icon from './Icon.js';

import VolumeStore from '../stores/VolumeStore.js';
import VolumeActions from '../actions/VolumeActions.js';


function _getOffsetLeft(elem) {
    var offsetLeft = 0;
    do {
      if (!isNaN(elem.offsetLeft)) {
          offsetLeft += elem.offsetLeft;
      }
    } while ((elem = elem.offsetParent));
    return offsetLeft;
}

function getVolumeState() {
  return {volume: VolumeStore.getVolume()};
}


export default class Volume extends React.Component {
  constructor(props) {
    super(props);

    this.state = getVolumeState();
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

    var w = `${Math.min(volume * 100.0, 100.0)}%`;
    var rest = `${100 - Math.min(volume * 100.0, 100.0)}%`;
    return (
      <div className="volume" onWheel={this._onWheel}>
        <div className="bar" onClick={this._onClick}>
          <div className="current" style={{width: w}} />
          <div className="marker" />
          <div className="rest" style={{width: rest}} />
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
    var pos = evt.pageX - _getOffsetLeft(evt.currentTarget);
    var width = evt.currentTarget.offsetWidth;
    VolumeActions.volume(pos/width);
  }

  _onChange() {
    this.setState(getVolumeState());
  }
}
