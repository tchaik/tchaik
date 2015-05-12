'use strict';

import React from 'react/addons';

import classNames from 'classnames';


class ArtworkImage extends React.Component {
  constructor(props) {
    super(props);
    this.state = {visible: false};

    this._onLoad = this._onLoad.bind(this);
    this._onError = this._onError.bind(this);
  }

  render() {
    var classes = {
      'visible': this.state.visible,
      'artwork': true,
    };
    return (
      <img src={this.props.path} className={classNames(classes)} onLoad={this._onLoad} onError={this._onError} />
    );
  }

  _onLoad() {
    this.setState({visible: true});
  }

  _onError() {
    this.setState({visible: false});
  }
}

export default ArtworkImage;
