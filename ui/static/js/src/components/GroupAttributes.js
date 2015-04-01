/** @jsx React.DOM */
'use strict';

var React = require('react/addons');

var GroupAttributes = React.createClass({
  render: function() {
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
});

module.exports = GroupAttributes;
