"use strict";

import React from "react/addons";

import classNames from "classnames";


export default class Icon extends React.Component {
  render() {
    var {icon, extraClasses, ...others} = this.props;
    var classes = {
      glyphicon: true,
    };
    classes["glyphicon-" + icon] = true;

    if (extraClasses) {
      for (var k in extraClasses) {
        classes[k] = extraClasses[k];
      }
    }

    return (
      <span {...others} className={classNames(classes)} />
    );
  }
}

Icon.propTypes = {
  icon: React.PropTypes.string.isRequired,
  extraClasses: React.PropTypes.object,
};
