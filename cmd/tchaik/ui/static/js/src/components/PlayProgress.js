import React from 'react/addons';

import NowPlayingStore from '../stores/NowPlayingStore.js';
import NowPlayingActions from '../actions/NowPlayingActions.js';

import PlayingStatusStore from '../stores/PlayingStatusStore.js';

function getNowPlayingState() {
  return {
    buffered: PlayingStatusStore.getBuffered(),
    duration: PlayingStatusStore.getDuration(),
    currentTime: PlayingStatusStore.getTime(),
  };
}

function _getOffsetLeft(elem) {
    var offsetLeft = 0;
    do {
      if (!isNaN(elem.offsetLeft)) {
          offsetLeft += elem.offsetLeft;
      }
    } while ((elem = elem.offsetParent));
    return offsetLeft;
}

export default class PlayProgress extends React.Component {
  constructor(props) {
    super(props);

    this.state = getNowPlayingState();

    this._onChange = this._onChange.bind(this);
    this._onClick = this._onClick.bind(this);
    this._onWheel = this._onWheel.bind(this);
  }

  componentDidMount() {
    NowPlayingStore.addChangeListener(this._onChange);
    PlayingStatusStore.addChangeListener(this._onChange);
  }

  componentWillUnmount() {
    NowPlayingStore.removeChangeListener(this._onChange);
    PlayingStatusStore.removeChangeListener(this._onChange);
  }

  _onChange() {
    this.setState(getNowPlayingState());
  }

  render() {
    var wpc = (this.state.currentTime / this.state.duration) * 100;
    var w = "calc("+Math.min(wpc, 100.0)+"% - " + this.props.markerWidth + "px)";
    var bpc = (this.state.buffered / this.state.duration) * 100 - wpc;
    var b = "calc("+Math.min(bpc, 100.0)+"% - " + this.props.markerWidth + "px)";

    return (
      <div className="playProgress" onClick={this._onClick} onWheel={this._onWheel}>
        <div className="bar">
          <div className="current" style={{width: w}} />
          <div className="marker" style={{width: this.props.markerWidth}} />
          <div className="buffered" style={{width: b}} />
        </div>
      </div>
    );
  }

  _onClick(evt) {
    var pos = evt.pageX - _getOffsetLeft(evt.currentTarget);
    var width = evt.currentTarget.offsetWidth;
    var time = (pos / width) * this.state.duration;
    NowPlayingActions.setCurrentTime(time);
  }

  _onWheel(evt) {
    evt.stopPropagation();
    var t = this.state.current + (0.01 * this.state.duration * evt.deltaY);
    if (t > this.state.duration) {
      t = this.props.duration;
    } else if (t < 0.00) {
      t = 0.0;
    }
    NowPlayingActions.setCurrentTime(t);
  }
}
