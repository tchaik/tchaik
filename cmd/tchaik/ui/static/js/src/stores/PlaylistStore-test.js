let rewire = require("rewire");

import PlaylistConstants from "../constants/PlaylistConstants.js";
import CollectionStoreMock from "./CollectionStore-mock.js";

describe("PlaylistStore", () => {
  let PlaylistStore;
  let dispatcherCallback, registerSpy, collectionStoreMock, collection;

  beforeEach(() => {
    // we have to clear localstorage, otherwise we have leaks from other tests
    localStorage.clear();


    // We can"t quite get the rewriting working, due to the fact that
    // AppDispatcher is called before anything is exported from the stores. The
    // only way we can mock the register function is by applying the spy on the
    // raw module using require rather than rewire.
    var AppDispatcher = require("../dispatcher/AppDispatcher.js");
    registerSpy = sinon.spy(function(callback) {
      dispatcherCallback = callback;
    });
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
        TrackID: "5557b3238ec1cdbd69db5e049954488446cce05e",
        Tracks: [
          {
            AlbumArtist: "Camo & Krooked",
            Artist: "Camo & Krooked",
            DiscNumber: 1,
            GroupName: "Above & Beyond",
            Key: "0",
            Name: "You Cry",
            TrackID: "5557b3238ec1cdbd69db5e049954488446cce05e",
            Year: 2010,
          },
          {
            AlbumArtist: "Camo & Krooked",
            Artist: "Camo & Krooked",
            DiscNumber: 1,
            GroupName: "Above & Beyond",
            Key: "1",
            Name: "Walk on Air",
            TrackID: "ffa00899e388b81cd3f2a0a971743aeffa447fe6",
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

  describe("Dispatcher callback", () => {
    describe("View actions", () => {
      describe("ADD_ITEM", () => {
        beforeEach(() => {
          dispatcherCallback({
            source: "VIEW_ACTION",
            action: {
              actionType: PlaylistConstants.ADD_ITEM,
              path: ["Root", "19sea9"],
            },
          });
        });

        it("adds the items to the playlist", () => {
          expect(PlaylistStore.getPlaylist()).to.eql([{
            data: {
              "Root:19sea9": {
                type: "TYPE_TRACKS",
                keys: [0, 1],
              },
            },
            paths: [[]],
            root: ["Root", "19sea9"],
            tracks: [[0], [1]],
          }, ]);
        });
      });
    });
  });
});
