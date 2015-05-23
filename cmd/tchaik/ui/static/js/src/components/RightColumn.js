"use strict";

import React from "react/addons";

import Playlist from "./Playlist.js";

import LeftColumnStore from "../stores/LeftColumnStore.js";

function getColumnState() {
  return {
    showPlaylist: LeftColumnStore.getMode() !== "Retro",
  };
}

export default class RightColumn extends React.Component {
  constructor(props) {
    super(props);

    this.state = getColumnState();
    this._onChange = this._onChange.bind(this);
  }

  componentDidMount() {
    LeftColumnStore.addChangeListener(this._onChange);
  }

  componentWillUnmount() {
    LeftColumnStore.removeChangeListener(this._onChange);
  }

  render() {
    var playlist = null;
    if (this.state.showPlaylist) {
      playlist = <Playlist />;
    }

    return (
      <div className="now-playing">
        {playlist}
      </div>
    );
  }

  _onChange() {
    this.setState(getColumnState());
  }
}
