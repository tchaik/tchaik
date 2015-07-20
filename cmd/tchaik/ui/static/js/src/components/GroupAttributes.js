"use strict";

import React from "react";

import SearchActions from "../actions/SearchActions.js";


function _dedupArray(arr, t, result) {
  arr.forEach(function(item) {
    if (Array.isArray(item)) {
      _dedupArray(item, t, result);
      return;
    }
    if (t.hasOwnProperty(item)) {
      return;
    }
    t[item] = true;
    result.push(item);
  });
  return result;
}

function dedupArray(arr) {
  let t = {};
  let result = [];
  _dedupArray(arr, t, result);
  return result;
}

function attributeLink(_this, x) {
  return [
    <a onClick={_this._onClickAttribute.bind(_this, x)}>{x}</a>,
    <span className="bull">&bull;</span>,
  ];
}

export default class GroupAttributes extends React.Component {
  render() {
    var _this = this;
    var list = dedupArray(this.props.list);
    list = list.map(function(attr) {
      return attributeLink(_this, attr);
    });

    if (list.length > 0) {
      list[list.length - 1].pop();
    }

    return (
      <div className="attributes">
        {list}
      </div>
    );
  }

  _onClickAttribute(attributeValue, evt) {
    evt.stopPropagation();
    SearchActions.search(attributeValue);
  }
}
