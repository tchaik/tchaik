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

const GroupAttributes = ({data, attributes}) => {
  if (attributes.length === 0) {
    return null;
  }

  const attr = [];
  for (const a of attributes) {
    if (data[a]) {
      attr.push(data[a]);
    }
  }

  if (attr.length === 0) {
    return null;
  }

  const onClick = (value, evt) => {
    evt.stopPropagation();
    SearchActions.search(value);
  }

  let list = dedupArray(attr);
  list = list.map(function(attr) {
    return [
      <a onClick={onClick.bind(null, attr)}>{attr}</a>,
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

GroupAttributes.propTypes = {
  attributes: React.PropTypes.array.isRequired,
  data: React.PropTypes.object.isRequired,
};

export default GroupAttributes;
