'use strict';

var React = require('react/addons');

var classNames = require('classnames');

var Icon = React.createClass({
  propTypes: {
    icon: React.PropTypes.string.isRequired,
    extraClasses: React.PropTypes.object,
  },

  render: function() {
    var {icon, extraClasses, ...others} = this.props;
    var classes = {
      glyphicon: true,
    };
    classes['glyphicon-' + icon] = true;

    if (extraClasses) {
      for (var k in extraClasses) {
        classes[k] = extraClasses[k];
      }
    }

    return (
      <span {...others} className={classNames(classes)} />
    );
  }
});

module.exports = Icon;
