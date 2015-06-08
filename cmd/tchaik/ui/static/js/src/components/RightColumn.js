"use strict";

import React from "react";

import Playlist from "./Playlist.js";

import ContainerConstants from "../constants/ContainerConstants.js";
import ContainerStore from "../stores/ContainerStore.js";

function getColumnState() {
  return {
    showPlaylist: ContainerStore.getMode() !== ContainerConstants.RETRO,
  };
}

export default class RightColumn extends React.Component {
  constructor(props) {
    super(props);

    this.state = getColumnState();
    this._onChange = this._onChange.bind(this);
  }

  componentDidMount() {
    ContainerStore.addChangeListener(this._onChange);
  }

  componentWillUnmount() {
    ContainerStore.removeChangeListener(this._onChange);
  }

  render() {
    var playlist = null;
    if (this.state.showPlaylist) {
      playlist = (
        <div className="now-playing">
          <Playlist />
        </div>
      );
    }

    return playlist;
  }

  _onChange() {
    this.setState(getColumnState());
  }
}
