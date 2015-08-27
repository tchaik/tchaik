"use strict";

import React from "react";
import classNames from "classnames";

import Playlist from "./Playlist.js";

import ContainerConstants from "../constants/ContainerConstants.js";
import ContainerStore from "../stores/ContainerStore.js";

import RightColumnStore from "../stores/RightColumnStore.js";

function getState() {
  return {
    showPlaylist: ContainerStore.getMode() !== ContainerConstants.RETRO,
    hidden: RightColumnStore.getIsHidden(),
  };
}

export default class RightColumn extends React.Component {
  constructor(props) {
    super(props);

    this.state = getState();
    this._onChange = this._onChange.bind(this);
  }

  componentDidMount() {
    ContainerStore.addChangeListener(this._onChange);
    RightColumnStore.addChangeListener(this._onChange);
  }

  componentWillUnmount() {
    ContainerStore.removeChangeListener(this._onChange);
    RightColumnStore.removeChangeListener(this._onChange);
  }

  render() {
    var playlist = null;
    if (this.state.showPlaylist) {
      var classes = classNames("now-playing", { hidden: this.state.hidden });
      playlist = (
        <div className={classes}>
          <Playlist />
        </div>
      );
    }

    return playlist;
  }

  _onChange() {
    this.setState(getState());
  }
}
