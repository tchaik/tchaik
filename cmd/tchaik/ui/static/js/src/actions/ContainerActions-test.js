var rewire = require("rewire");

import AppDispatcherMock from "../dispatcher/AppDispatcher-mock.js";
import ContainerConstants from "../constants/ContainerConstants.js";

describe("ContainerActions", () => {
  var ContainerActions;

  var appDispatcherMock;

  beforeEach(() => {
    ContainerActions = rewire("./ContainerActions.js");

    appDispatcherMock = new AppDispatcherMock();
    ContainerActions.__Rewire__("AppDispatcher", appDispatcherMock);
    sinon.stub(appDispatcherMock, "handleViewAction");
  });

  describe("function: mode", () => {
    var mode;

    beforeEach(() => {
      mode = "RETRO";
      ContainerActions.mode(mode);
    });

    it("should call AppDispatcher.handleViewAction", () => {
      expect(appDispatcherMock.handleViewAction).to.have.been.calledWith({
        actionType: ContainerConstants.MODE,
        mode: mode,
      });
    });
  });
});
