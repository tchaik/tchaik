'use strict';

import React from 'react/addons';

import ArtworkImage from './ArtworkImage.js';
import Playlist from './Playlist.js';

import NowPlayingStore from '../stores/NowPlayingStore.js';


export default class Retro extends React.Component {
  constructor(props) {
    super(props);

    this.state = {track: NowPlayingStore.getTrack()};
    this._onChange = this._onChange.bind(this);
  }

  componentDidMount() {
    NowPlayingStore.addChangeListener(this._onChange);
  }

  componentWillUnmount() {
    NowPlayingStore.removeChangeListener(this._onChange);
  }

  render() {
    var artworkImage = null;
    var imagePath = null;
    if (this.state.track) {
      imagePath = "/artwork/" + this.state.track.TrackID;
      artworkImage = <ArtworkImage path={imagePath} />;
    }

    return (
      <div className="retro">
        <div className="blur" style={{'backgroundImage': "url("+imagePath+")"}} />
          <div className="current-artwork">
            {artworkImage}
          </div>
        <div className="content">
          <Playlist />
        </div>
      </div>
    );
  }

  _onChange() {
    this.setState({track: NowPlayingStore.getTrack()});
  }
}
