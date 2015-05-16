'use strict';

require('../../sass/screen.scss');
require('../../sass/glyphicons.scss');

var React = require('react/addons');

require('babel-core/polyfill');

var WebsocketAPI = require('./utils/WebsocketAPI.js');
var AudioAPI = require('./utils/AudioAPI.js');

var LeftColumn = require('./components/LeftColumn.js');
var RightColumn = require('./components/RightColumn.js');
var Bottom = require('./components/Bottom.js');
var Top = require('./components/Top.js');

var socketAddr = document.location.host;

var protocol = "ws://";
if (window.location.protocol === "https:") {
  protocol = "wss://";
}

var websocketUrl = `${protocol}${socketAddr}/socket`;
if (process.env.WS_URL) {
  websocketUrl = process.env.WS_URL;
}
WebsocketAPI.init(websocketUrl);

AudioAPI.init();

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

var Bottom = React.createFactory(Bottom);
React.render(
  Bottom(),
  document.getElementById('bottom')
);

var Top = React.createFactory(Top);
React.render(
  Top(),
  document.getElementById('top')
);
