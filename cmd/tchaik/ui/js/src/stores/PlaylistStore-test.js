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
        AlbumArtist: "Camo & Krooked",
        Key: "19sea9",
        ListStyle: "",
        Name: "Above & Beyond",
        TotalTime: 0,
        ID: "5557b3238ec1cdbd69db5e049954488446cce05e",
        Tracks: [
          {
            AlbumArtist: "Camo & Krooked",
            Artist: "Camo & Krooked",
            DiscNumber: 1,
            GroupName: "Above & Beyond",
            Key: "0",
            Name: "You Cry",
            ID: "5557b3238ec1cdbd69db5e049954488446cce05e",
            Year: 2010,
          },
          {
            AlbumArtist: "Camo & Krooked",
            Artist: "Camo & Krooked",
            DiscNumber: 1,
            GroupName: "Above & Beyond",
            Key: "1",
            Name: "Walk on Air",
            ID: "ffa00899e388b81cd3f2a0a971743aeffa447fe6",
            Year: 2010,
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
