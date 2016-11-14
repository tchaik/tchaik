"use strict";

import React from "react";

function zeroPad(n, width) {
  let s = "" + n;
  while (s.length < width) {
    s = "0" + s;
  }
  return s;
}

function timeText(time) {
  let text = "";
  let minsPad = 0;
  let totalSeconds = parseInt(time);

  const hours = parseInt(totalSeconds / 3600);
  if (hours > 0) {
    text += hours + ":";
    totalSeconds %= 3600;
    minsPad = 2;
  }

  const mins = parseInt(totalSeconds / 60);
  const secs = parseInt(totalSeconds % 60);
  text += zeroPad(mins, minsPad) + ":" + zeroPad(secs, 2);
  return text
}

const TimeFormatter = ({time, ...others}) => {
  if (isNaN(time)) {
    return null;
  }
  return <span {...others}>{timeText(time)}</span>;
}

export default TimeFormatter;
