let rewire = require("rewire");

import CollectionStoreMock from "./CollectionStore-mock.js";

describe("PlaylistStore", () => {
  let PlaylistStore;
  let registerSpy, collectionStoreMock, collection;

  beforeEach(() => {
    // we have to clear localstorage, otherwise we have leaks from other tests
    localStorage.clear();


    // We can"t quite get the rewriting working, due to the fact that
    // AppDispatcher is called before anything is exported from the stores. The
    // only way we can mock the register function is by applying the spy on the
    // raw module using require rather than rewire.
    var AppDispatcher = require("../dispatcher/AppDispatcher.js");
    registerSpy = sinon.spy();
    AppDispatcher.register = registerSpy;

    PlaylistStore = rewire("./PlaylistStore.js");

    collectionStoreMock = new CollectionStoreMock();
    PlaylistStore.__Rewire__("CollectionStore", collectionStoreMock);


    collection = {
      "Root:19sea9": {
        albumArtist: "Camo & Krooked",
        key: "19sea9",
        listStyle: "",
        name: "Above & Beyond",
        totalTime: 0,
        id: "5557b3238ec1cdbd69db5e049954488446cce05e",
        tracks: [
          {
            albumArtist: "Camo & Krooked",
            artist: "Camo & Krooked",
            discNumber: 1,
            groupName: "Above & Beyond",
            key: "0",
            name: "You Cry",
            id: "5557b3238ec1cdbd69db5e049954488446cce05e",
            year: 2010,
          },
          {
            albumArtist: "Camo & Krooked",
            artist: "Camo & Krooked",
            discNumber: 1,
            groupName: "Above & Beyond",
            key: "1",
            name: "Walk on Air",
            id: "ffa00899e388b81cd3f2a0a971743aeffa447fe6",
            year: 2010,
          },
        ],
      },
    };
    sinon.stub(collectionStoreMock, "getCollection", (path) => {
      return collection[path.join(":")];
    });
    sinon.stub(collectionStoreMock, "pathToKey", (path) => {
      if (path) {
        return path.join(":");
      }
      return null;
    });
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

  describe("function: getItemKeys", () => {
  });
});
