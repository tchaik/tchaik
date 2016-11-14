"use strict";

import React from "react";

import classNames from "classnames";

const Icon = ({icon, extraClasses, ...others}) => {
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

Icon.propTypes = {
  icon: React.PropTypes.string.isRequired,
  extraClasses: React.PropTypes.object,
};

export default Icon;
