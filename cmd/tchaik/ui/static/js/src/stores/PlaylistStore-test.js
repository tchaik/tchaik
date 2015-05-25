var rewire = require("rewire");

describe("PlaylistStore", () => {
  let PlaylistStore;
  //let callback;
  let registerSpy;

  beforeEach(() => {
    // we have to clear localstorage, otherwise we have leaks from other tests
    localStorage.clear();


    // We can"t quite get the rewriting working, due to the fact that
    // AppDispatcher is called before anything is exported from the stores. The
    // only way we can mock the register function is by applying the spy on the
    // raw module using require rather than rewire.
    var AppDispatcher = require("../dispatcher/AppDispatcher.js");
    registerSpy = sinon.spy(function() {
      //callback = cb;
    });
    AppDispatcher.register = registerSpy;

    PlaylistStore = rewire("./PlaylistStore.js");
  });

  describe("on initialisation", () => {
    it("registers a callback with the dispatcher", () => {
      expect(registerSpy).to.have.been.called;
    });
  });

  describe("function: getPlaylist", () => {
    describe("on initialisation", () => {
      it("should return an empty array", () => {
        expect(PlaylistStore.getPlaylist()).to.be.empty;
      });
    });
  });

  describe("function: getCurrent", () => {
    describe("on initialisation", () => {
      it("should return null", () => {
        expect(PlaylistStore.getCurrent()).to.be.null;
      });
    });
  });

  describe("function: getCurrentTrack", () => {
    describe("on initialisation", () => {
      it("should return null", () => {
        expect(PlaylistStore.getCurrentTrack()).to.be.null;
      });
    });
  });

  describe("function: canPrev", () => {
    describe("on initialisation", () => {
      it("should return falsy", () => {
        expect(PlaylistStore.canPrev()).to.be.falsy;
      });
    });
  });

  describe("function: canNext", () => {
    describe("on initialisation", () => {
      it("should return falsy", () => {
        expect(PlaylistStore.canNext()).to.be.falsy;
      });
    });
  });

  describe("function: getNext", () => {
    describe("on initialisation", () => {
      it("should return null", () => {
        expect(PlaylistStore.getNext()).to.be.null;
      });
    });
  });

  describe("function: getItemKeys", () => {
  });
});
