/** @jsx React.DOM */
'use strict';

var React = require('react/addons');

var classNames = require('classnames');

var ArtworkImage = React.createClass({
  getInitialState: function() {
    return {
      visible: false,
    };
  },

  render: function() {
    var classes = {
      'visible': this.state.visible,
      'artwork': true,
    };
    return (
      <img src={this.props.path} className={classNames(classes)} onLoad={this._onLoad} onError={this._onError} />
    );
  },

  _onLoad: function() {
    this.setState({visible: true});
  },

  _onError: function() {
    this.setState({visible: false});
  },
});

module.exports = ArtworkImage;