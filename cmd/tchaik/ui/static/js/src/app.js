'use strict';

var React = require('react/addons');

var WebsocketApi = require('./api/WebsocketApi.js');

var LeftColumn = require('./components/LeftColumn.js');
var RightColumn = require('./components/RightColumn.js');

var socketAddr = document.location.host;
var protocol = "ws://";
if (window.location.protocol === "https:") {
  protocol = "wss://";
}
WebsocketApi.init(protocol + socketAddr + "/socket");

var LeftColumn = React.createFactory(LeftColumn);
React.render(
  LeftColumn(),
  document.getElementById('left-column')
);

var RightColumn = React.createFactory(RightColumn);
React.render(
  RightColumn(),
  document.getElementById('right-column')
);
