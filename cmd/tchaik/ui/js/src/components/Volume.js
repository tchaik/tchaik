"use strict";

import React from "react";

import Icon from "./Icon.js";

import { connect, Provider } from "react-redux";
import { setVolume, setMute } from "../redux/Volume.js";

function _getOffsetLeft(elem) {
  var offsetLeft = 0;
  do {
    if (!isNaN(elem.offsetLeft)) {
      offsetLeft += elem.offsetLeft;
    }
  } while ((elem = elem.offsetParent));
  return offsetLeft;
}

class volumeBar extends React.Component {
  constructor(props) {
    super(props);

    this._onClick = this._onClick.bind(this);
    this._onWheel = this._onWheel.bind(this);
    this._toggleMute = this._toggleMute.bind(this);
  }

  render() {
    let { volume, mute } = this.props;
    if (mute) {
      volume = 0.00;
    }

    let classSuffix = "";
    if (volume === 0.00) {
      classSuffix = "mute";
    } else if (volume < 0.5) {
      classSuffix = "down";
    } else {
      classSuffix = "up";
    }

    const w = `${Math.min(volume * 100.0, 100.0)}%`;
    const rest = `${100 - Math.min(volume * 100.0, 100.0)}%`;
    return (
      <div className="volume" onWheel={this._onWheel}>
        <div className="bar" onClick={this._onClick}>
          <div className="current" style={{width: w}} />
          <div className="marker" />
          <div className="rest" style={{width: rest}} />
        </div>
        <Icon icon={"volume_" + classSuffix} onClick={this._toggleMute} />
      </div>
    );
  }

  _toggleMute(evt) {
    evt.stopPropagation();
    this.props.setMute(!this.props.mute)
  }

  _onWheel(evt) {
    evt.stopPropagation();
    let v = this.props.volume + 0.05 * evt.deltaY;
    if (v > 1.0) {
      v = 1.0;
    } else if (v < 0.00) {
      v = 0.0;
    }
    this.props.setVolume(v);
  }

  _onClick(evt) {
    const pos = evt.pageX - _getOffsetLeft(evt.currentTarget);
    const width = evt.currentTarget.offsetWidth;
    this.props.setVolume(pos / width);
  }
}

const Volume = connect(
  state => (state),
  dispatch => ({
    setVolume: volume => dispatch(setVolume(volume)),
    setMute: mute => dispatch(setMute(mute)),
  })
)(volumeBar)

export default Volume
