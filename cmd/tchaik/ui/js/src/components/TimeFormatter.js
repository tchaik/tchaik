"use strict";

import React from "react";

function zeroPad(n, width) {
  var s = "" + n;
  while (s.length < width) {
    s = "0" + s;
  }
  return s;
}

export default class TimeFormatter extends React.Component {
  shouldComponentUpdate(nextProps, nextState) {
    return this.props.time !== nextProps.time;
  }

  render() {
    var {time, ...others} = this.props;
    if (isNaN(time)) {
      return null;
    }

    var timeText = "";
    var minsPad = 0;

    var totalSeconds = parseInt(time);
    var hours = parseInt(totalSeconds / 3600);
    if (hours > 0) {
      timeText += hours + ":";
      totalSeconds %= 3600;
      minsPad = 2;
    }

    var mins = parseInt(totalSeconds / 60);
    var secs = parseInt(totalSeconds % 60);
    timeText += zeroPad(mins, minsPad) + ":" + zeroPad(secs, 2);

    return (
      <span {...others}>{timeText}</span>
    );
  }
}
