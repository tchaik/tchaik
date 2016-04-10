"use strict";

import {ChangeEmitter} from "../utils/ChangeEmitter.js";

import AppDispatcher from "../dispatcher/AppDispatcher";

import FilterConstants from "../constants/FilterConstants.js";


let _currentItems = null;
const _filters = {};
const _filterPaths = {};

function _getCurrentItem(name) {
  if (_currentItems === null) {
    const fi = localStorage.getItem("filterCurrentItems");
    if (fi) {
      _currentItems = JSON.parse(fi);
    }
  }
  if (!_currentItems || !_currentItems[name]) {
    return null;
  }
  return _currentItems[name];
}

function _setCurrentItem(name, itemName) {
  if (_currentItems === null) {
    _currentItems = {};
  }
  _currentItems[name] = itemName;
  localStorage.setItem("filterCurrentItems", JSON.stringify(_currentItems));
}

function _addFilterPaths(name, itemName, paths) {
  let filterPaths = _filterPaths[name];
  if (!filterPaths) {
    filterPaths = {};
  }
  filterPaths[itemName] = paths;
  _filterPaths[name] = filterPaths;
}


class FilterStore extends ChangeEmitter {
  getCurrentItem(name) {
    return _getCurrentItem(name);
  }

  getItems(name) {
    let x = _filters[name];
    if (!x) {
      x = [];
    }
    return x;
  }

  getPaths(name, itemName) {
    const x = _filterPaths[name];
    if (!x) {
      return [];
    }
    let paths = x[itemName];
    if (!paths) {
      paths = [];
    }
    return paths;
  }
}

const _filterStore = new FilterStore();

_filterStore.dispatchToken = AppDispatcher.register(function(payload) {
  const action = payload.action;
  const source = payload.source;

  if (source === "SERVER_ACTION") {
    switch (action.actionType) {
      case FilterConstants.FILTER_LIST:
        _filters[action.data.name] = action.data.items;
        _filterStore.emitChange();
        break;
      case FilterConstants.FILTER_PATHS:
        const path = action.data.path; // [name, item]
        _addFilterPaths(path[0], path[1], action.data.paths.groups);
        _filterStore.emitChange();
        break;
    }
  } else if (source === "VIEW_ACTION") {
    switch (action.actionType) {
      case FilterConstants.SET_FILTER_ITEM:
        _setCurrentItem(action.name, action.itemName);
        _filterStore.emitChange();
        break;
    }
  }

  return true;
});

export default _filterStore;
