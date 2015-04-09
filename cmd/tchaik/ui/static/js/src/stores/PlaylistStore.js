'use strict';

var AppDispatcher = require('../dispatcher/AppDispatcher');
var EventEmitter = require('eventemitter3').EventEmitter;
var assign = require('object-assign');

var CollectionStore = require('./CollectionStore.js');

var CollectionConstants = require('../constants/CollectionConstants.js');
var PlaylistConstants = require('../constants/PlaylistConstants.js');
var NowPlayingConstants = require('../constants/NowPlayingConstants.js');
var ControlApiConstants = require('../constants/ControlApiConstants.js');

var CHANGE_EVENT = 'change';

var defaultPlaylist = {
  items: [],
  current: null,
};

var currentPlaylist = null;

function buildPlaylistItem(root) {
  var queue = [[]];
  var tracks = [];
  function getGroupKey(path) {
    return function(g) {
      var k = g.Key;
      queue.push(path.concat(g.Key));
      return k;
    };
  }

  function getTrackKey(path) {
    var i = 0;
    return function() {
      tracks.push(path.concat(i));
      return i++;
    };
  }

  var data = {};
  var paths = [];
  while (queue.length > 0) {
    var p = queue.shift();
    var rp = root.concat(p); // path from root
    var c = CollectionStore.getCollection(rp);

    var x = {};
    if (c.Groups) {
      x = {
        type: PlaylistConstants.TYPE_GROUP,
        keys: c.Groups.map(getGroupKey(p)),
      };
    } else if (c.Tracks) {
      x = {
        type: PlaylistConstants.TYPE_TRACKS,
        keys: c.Tracks.map(getTrackKey(p)),
      };
    } else {
      console.error("expected c.Groups or c.Tracks to be non-null");
    }
    paths.push(p);
    data[CollectionStore.pathToKey(rp)] = x;
  }

  return {
    root: root,
    paths: paths,
    data: data,
    tracks: tracks,
  };
}

function reset() {
  var p = playlist();
  if (p.items.length > 0) {
    p.current = {
      item: 0,
      track: 0,
      path: p.items[0].tracks[0],
    };
  } else {
    p.current = null;
  }
  setPlaylist(p);
}

function next() {
  var p = playlist();
  var c = p.current;
  if (c === null) {
    reset();
    return;
  }

  var item = p.items[c.item];
  var tracks = item.tracks;
  if (c.track + 1 < tracks.length) {
    c.path = item.root.concat(tracks[++c.track]);
  } else if (c.item + 1 < p.items.length) {
    item = p.items[++c.item];
    c.track = 0;
    c.path = item.root.concat(item.tracks[0]);
  } else {
    // Overflow - at the end - must be last item.
    p.current = null;
  }
  setPlaylist(p);
}

function prev() {
  var p = playlist();
  var c = p.current;
  if (c === null) {
    reset();
    return;
  }

  var item = p.items[c.item];
  var tracks = item.tracks;
  if (c.track - 1 >= 0) {
    c.path = item.root.concat(tracks[--c.track]);
    setPlaylist(p);
    return;
  }

  if (c.item - 1 >= 0) {
    item = p.items[--c.item];
    c.track = item.tracks.length-1;
    c.path = item.root.concat(item.tracks[c.track]);
    setPlaylist(p);
    return;
  }
}

function canPrev() {
  var p = playlist();
  var c = p.current;
  if (c === null) {
    return false;
  }

  if (c.track > 0 || (c.track === 0 && c.item > 0)) {
    return true;
  }
  return false;
}

function canNext() {
  var p = playlist();
  var c = p.current;
  if (c === null) {
    return false;
  }

  var tracks = p.items[c.item].tracks;
  if (c.track < (tracks.length - 1) || (c.item < (p.items.length-1))) {
    return true;
  }
  return false;
}

function isPathPrefix(path, prefix) {
  if (prefix.length > path.length) {
    return false;
  }
  for (var i = 0; i < prefix.length; i++) {
    if (path[i] != prefix[i]) {
      return false;
    }
  }
  return true;
}

function pathsEqual(p1, p2) {
  if (p1.length != p2.length) {
    return false;
  }
  return isPathPrefix(p1, p2);
}

function _removeTracks(tracks, path) {
  var i = 0;
  while (i < tracks.length) {
    if (isPathPrefix(tracks[i], path)) {
      tracks.splice(i, 1);
      continue;
    }
    i++;
  }
}

function _removePaths(paths, data, path) {
  var i = 0;
  while (i < paths.length) {
    if (isPathPrefix(paths[i], path)) {
      delete data[CollectionStore.pathToKey(paths[i])];
      paths.splice(i, 1);
      continue;
    }
    i++;
  }
}

function _removeItem(playlist, itemIndex) {
  if (playlist.current !== null) {
    var current = playlist.current;
    if (current.item > itemIndex) {
      current.item--;
    }
  }
  playlist.items.splice(itemIndex, 1);
}

function remove(itemIndex, path) {
  var p = playlist();
  var item = p.items[itemIndex];

  var data;
  do {
    if (pathsEqual(path, item.root)) {
      _removeItem(p, itemIndex);
      break;
    }

    _removeTracks(item.tracks, path);
    _removePaths(item.paths, item.data, path);

    var last = path.pop();
    data = item.data[CollectionStore.pathToKey(path)];
    if (!data || !data.keys) {
      break;
    }
    data.keys.splice(data.keys.indexOf(last), 1);
  } while(data.keys.length === 0 && path.length > 1);

  p.current = null;
  setPlaylist(p);
}

function setCurrent(itemIndex, path) {
  var p = playlist();
  var item = p.items[itemIndex];
  var key = CollectionStore.pathToKey(path);

  var track = -1;
  for (var i = 0; i < item.tracks.length; i++) {
    if (CollectionStore.pathToKey(item.root.concat(item.tracks[i])) === key) {
      track = i;
      break;
    }
  }
  if (track === -1) {
    console.error("Could not find track for path:"+path);
  }

  p.current = {
    item: itemIndex,
    path: path,
    track: track,
  };
  setPlaylist(p);
}

function currentTrack() {
  var p = playlist();
  var c = p.current;

  if (c === null) {
    return null;
  }

  var path = c.path.slice(0);
  var i = path.pop();
  var t = CollectionStore.getCollection(path);
  if (t === null) {
    console.log("Could not find collection with path:"+path);
  }

  if (t.Tracks) {
    var track = t.Tracks[i];
    if (track) {
      return track;
    }
    console.log("No track found for index: "+i);
    console.log(t.Tracks);
    return null;
  }

  console.log("Collection item did not have Tracks property");
  return null;
}

function append(path) {
  var p = playlist();
  p.items.push(buildPlaylistItem(path));
  setPlaylist(p);
}

function now(path) {
  var p = playlist();
  p.items.unshift(buildPlaylistItem(path));
  setPlaylist(p);
  reset();
}

function setPlaylist(p) {
  currentPlaylist = p;
  localStorage.setItem("playlist", JSON.stringify(p));
}

function _playlist() {
  var p = localStorage.getItem("playlist");
  if (p === null) {
    return defaultPlaylist;
  }
  return JSON.parse(p);
}

function playlist() {
  if (currentPlaylist === null) {
    currentPlaylist = _playlist();
  }
  return currentPlaylist;
}

var PlaylistStore = assign({}, EventEmitter.prototype, {

  getPlaylist: function() {
    return playlist().items;
  },

  getCurrent: function() {
    return playlist().current;
  },

  getCurrentTrack: function() {
    return currentTrack();
  },

  canPrev: function() {
    return canPrev();
  },

  canNext: function() {
    return canNext();
  },

  getItemKeys: function(index, path) {
    var p = playlist();
    var items = p.items;
    var key = CollectionStore.pathToKey(path);
    var item = items[index];

    return item.data[key];
  },

  removeItem: function(itemIndex, path) {
    remove(itemIndex, path);
  },

  emitChange: function() {
    this.emit(CHANGE_EVENT);
  },

  /**
   * @param {function} callback
   */
  addChangeListener: function(callback) {
    this.on(CHANGE_EVENT, callback);
  },

  /**
   * @param {function} callback
   */
  removeChangeListener: function(callback) {
    this.removeListener(CHANGE_EVENT, callback);
  }

});

PlaylistStore.dispatchToken = AppDispatcher.register(function(payload) {
  var action = payload.action;
  var source = payload.source;

  if (source === 'SERVER_ACTION') {
    if (action.actionType === ControlApiConstants.CTRL) {
      switch (action.data) {
        case ControlApiConstants.NEXT:
          next();
          PlaylistStore.emitChange();
          break;

        case ControlApiConstants.PREV:
          prev();
          PlaylistStore.emitChange();
          break;
      }
    }
  }

  if (source === 'VIEW_ACTION') {
    switch (action.actionType) {

      case NowPlayingConstants.ENDED:
        if (action.source !== "playlist") {
          break;
        }
        /* falls through */
      case PlaylistConstants.NEXT:
        next();
        PlaylistStore.emitChange();
        break;

      case PlaylistConstants.PREV:
        prev();
        PlaylistStore.emitChange();
        break;

      case PlaylistConstants.REMOVE:
        remove(action.itemIndex, action.path);
        PlaylistStore.emitChange();
        break;

      case CollectionConstants.APPEND_TO_PLAYLIST:
        append(action.path);
        PlaylistStore.emitChange();
        break;

      case CollectionConstants.PLAY_NOW:
        now(action.path);
        PlaylistStore.emitChange();
        break;

      case PlaylistConstants.PLAY_ITEM:
        setCurrent(action.itemIndex, action.path);
        PlaylistStore.emitChange();
        break;

      default:
        break;
    }
  }

  return true;
});

module.exports = PlaylistStore;
