var rewire = require("rewire");

import AppDispatcherMock from "../dispatcher/AppDispatcher-mock.js";
import CollectionConstants from "../constants/CollectionConstants.js";
import CollectionStoreMock from "../stores/CollectionStore-mock.js";
import NowPlayingConstants from "../constants/NowPlayingConstants.js";
import WebsocketAPIMock from "../utils/WebsocketAPI-mock.js";

describe("CollectionActions", () => {
  const PATH = [];

  var CollectionActions;

  var appDispatcherMock, collectionStoreMock;
  beforeEach(() => {
    CollectionActions = rewire("./CollectionActions.js");

    appDispatcherMock = new AppDispatcherMock();
    collectionStoreMock = new CollectionStoreMock();

    CollectionActions.__Rewire__("AppDispatcher", appDispatcherMock);
    CollectionActions.__Rewire__("CollectionStore", collectionStoreMock);

    sinon.stub(appDispatcherMock, "handleViewAction");
  });

  describe("function: fetch", () => {
    var websocketAPIMock;
    beforeEach(() => {
      websocketAPIMock = new WebsocketAPIMock();

      CollectionActions.__Rewire__("WebsocketAPI", websocketAPIMock);

      sinon.stub(collectionStoreMock, "getCollection");
      sinon.stub(collectionStoreMock, "emitChange");
      sinon.stub(websocketAPIMock, "send");
    });

    describe("if the path is already in the collection store", () => {
      beforeEach(() => {
        collectionStoreMock.getCollection.returns(true);

        CollectionActions.fetch(PATH);
      });

      it("should call emitChange on the CollectionStore", () => {
        expect(collectionStoreMock.emitChange).to.have.been.called;
      });

      it("should not send a message to the server through the websocket", () => {
        expect(websocketAPIMock.send).not.to.have.been.called;
      });
    });

    describe("if the path is not in the collection store", () => {
      beforeEach(() => {
        collectionStoreMock.getCollection.returns(false);
        CollectionActions.fetch(PATH);
      });

      it("should call emitChange on the CollectionStore", () => {
        expect(collectionStoreMock.emitChange).not.to.have.been.called;
      });

      it("should not send a message to the server through the websocket", () => {
        expect(websocketAPIMock.send).to.have.been.calledWith(
          CollectionConstants.FETCH, {path: PATH}
        );
      });
    });
  });

  describe("function: setCurrentTrack", () => {
    it("should call handleViewAction", () => {
      var track = "98789722";
      CollectionActions.setCurrentTrack(track);

      expect(appDispatcherMock.handleViewAction).to.have.been.calledWith({
        actionType: NowPlayingConstants.SET_CURRENT_TRACK,
        track: track,
        source: "collection",
      });
    });
  });

});
