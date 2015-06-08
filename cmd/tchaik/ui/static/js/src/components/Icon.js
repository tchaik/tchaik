"use strict";

import React from "react";

import classNames from "classnames";


export default class Icon extends React.Component {
  render() {
    var {icon, extraClasses, fa, ...others} = this.props;
    var classes = {
      icon: true,
    };

    if (fa) {
      classes.fa = true;
      classes["fa-" + icon] = true;
    } else {
      classes.glyphicon = true;
      classes["glyphicon-" + icon] = true;
    }

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
