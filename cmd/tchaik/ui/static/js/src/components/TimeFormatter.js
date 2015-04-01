/** @jsx React.DOM */
'use strict';

var React = require('react/addons');

function zeroPad(n, width) {
  var s = "" + n;
  while (s.length < width) {
    s = "0" + s;
  }
  return s;
}

var TimeFormatter = React.createClass({

  render: function() {
    var {time, ...others} = this.props;
    var totalSeconds = parseInt(time);
    var timeText = "";
    var minsPad = 0;

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
});

module.exports = TimeFormatter;
