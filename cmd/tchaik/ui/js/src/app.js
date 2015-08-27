"use strict";

require("../../css/screen.css");
require("../../css/material-icons.css");

var React = require("react");

require("babel-core/polyfill");

var WebsocketAPI = require("./utils/WebsocketAPI.js");
var AudioAPI = require("./utils/AudioAPI.js");

var LeftColumn = require("./components/LeftColumn.js");
var RightColumn = require("./components/RightColumn.js");
var Bottom = require("./components/Bottom.js");
var Top = require("./components/Top.js");
var Container = require("./components/Container.js");

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

React.render(
  React.createFactory(LeftColumn)(),
  document.getElementById("left-column")
);

React.render(
  React.createFactory(RightColumn)(),
  document.getElementById("right-column")
);

React.render(
  React.createFactory(Bottom)(),
  document.getElementById("bottom")
);

React.render(
  React.createFactory(Top)(),
  document.getElementById("top")
);

React.render(
  React.createFactory(Container)(),
  document.getElementById("container")
);
