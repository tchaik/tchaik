"use strict";

import React from "react/addons";

import SearchActions from "../actions/SearchActions.js";

function dedupArray(arr) {
  var t = {};
  var result = [];
  arr.forEach(function(item) {
    if (t.hasOwnProperty(item)) {
      return;
    }
    t[item] = true;
    result.push(item);
  });
  return result;
}

export default class GroupAttributes extends React.Component {
  render() {
    var _this = this;
    var list = dedupArray(this.props.list);
    list = list.map(function(attr) {
      return [
        <a onClick={_this._onClickAttribute.bind(_this, attr)}>{attr}</a>,
        <span className="bull">&bull;</span>,
      ];
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
