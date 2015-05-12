'use strict';

import React from 'react/addons';

import Playlist from './Playlist.js';

export default class RightColumn extends React.Component {
  render() {
    return (
      <div className="now-playing">
        <Playlist />
      </div>
    );
  }
}
