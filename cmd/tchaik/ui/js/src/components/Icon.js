"use strict";

import React from "react";

import classNames from "classnames";


export default class Icon extends React.Component {
  shouldComponentUpdate(nextProps, nextState) {
    return this.props.icon != nextProps.icon;
  }

  render() {
    var {icon, extraClasses, ...others} = this.props;
    var classes = {
      icon: true,
    };

    classes["material-icons"] = true;

    if (extraClasses) {
      for (let k in extraClasses) {
        classes[k] = extraClasses[k];
      }
    }

    return <span {...others} className={classNames(classes)}>{icon}</span>;
  }
}

Icon.propTypes = {
  icon: React.PropTypes.string.isRequired,
  extraClasses: React.PropTypes.object,
};
