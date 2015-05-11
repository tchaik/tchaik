'use strict';

var React = require('react/addons');

var Playlist = require('./Playlist.js');

var RightColumn = React.createClass({

  render: function() {
    return (
      <div className="now-playing">
        <Playlist />
      </div>
    );
  },

});

module.exports = RightColumn;
