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

import { Provider } from "react-redux"
import { volumeStore } from "./redux/Volume.js";


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
  <Provider store={volumeStore}>
    <Bottom />
  </Provider>,
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

// ChromeCast
// function sessionListener(e) {
//   let session = e;
//   console.log("sessionListener", session);
//   if (session.media.length !== 0) {
//     console.log("onRequestSessionSuccess", session.media[0]);
//   }
// }
//
// function receiverListener(e) {
//   let session = e;
//   if (e === chrome.cast.ReceiverAvailability.AVAILABLE) {
//     chrome.cast.requestSession(onRequestSessionSuccess, onLaunchError);
//   }
//   console.log("receiverListener", session);
//   if (session.media.length !== 0) {
//     console.log("onRequestReceiverSuccess", session.media[0]);
//   }
// }
//
// function onRequestSessionSuccess(e) {
//   session = e;
// }
//
// function onInitSuccess() {
//   console.log("init success!");
// }
//
// function onError() {
//   console.log("onError!");
// }
//
// var initializeCastApi = function() {
//   var sessionRequest = new chrome.cast.SessionRequest(chrome.cast.media.DEFAULT_MEDIA_RECEIVER_APP_ID);
//   var apiConfig = new chrome.cast.ApiConfig(sessionRequest,
//     sessionListener,
//     receiverListener);
//   chrome.cast.initialize(apiConfig, onInitSuccess, onError);
// };
//
// window.__onGCastApiAvailable = function(loaded, errorInfo) {
//   if (loaded) {
//     initializeCastApi();
//   } else {
//     console.log("error from onGCastApiAvailable");
//     console.log(errorInfo);
//   }
// };
