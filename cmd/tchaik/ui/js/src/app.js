"use strict";

import "../../css/screen.css";
import "../../css/material-icons.css";

import React from "react";
import ReactDOM from "react-dom";

import WebsocketAPI from "./utils/WebsocketAPI.js";

import LeftColumn from "./components/LeftColumn.js";
import RightColumn from "./components/RightColumn.js";
import Bottom from "./components/Bottom.js";
import Top from "./components/Top.js";
import Container from "./components/Container.js";

const socketAddr = document.location.host;

let protocol = "ws://";
if (window.location.protocol === "https:") {
  protocol = "wss://";
}

let websocketUrl = `${protocol}${socketAddr}/socket`;
if (process.env.WS_URL) {
  websocketUrl = process.env.WS_URL;
}

WebsocketAPI.init(websocketUrl);

ReactDOM.render(
  React.createFactory(LeftColumn)(),
  document.getElementById("left-column")
);

ReactDOM.render(
  React.createFactory(RightColumn)(),
  document.getElementById("right-column")
);

ReactDOM.render(
  React.createFactory(Bottom)(),
  document.getElementById("bottom")
);

ReactDOM.render(
  React.createFactory(Top)(),
  document.getElementById("top")
);

ReactDOM.render(
  React.createFactory(Container)(),
  document.getElementById("container")
);
