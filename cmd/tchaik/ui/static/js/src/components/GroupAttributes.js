'use strict';

import React from 'react/addons';

export default class GroupAttributes extends React.Component {
  render() {
    var list = this.props.list.map(function(attr) {
      return [
        <span>{attr}</span>,
        <span className="bull">&bull;</span>,
      ];
    });
    list[list.length-1].pop();

    return (
      <div className="attributes">
        {list}
      </div>
    );
  }
}
