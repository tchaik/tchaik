'use strict';

import React from 'react/addons';

import RecentStore from '../stores/RecentStore.js';
import RecentActions from '../actions/RecentActions.js';

import CollectionStore from '../stores/CollectionStore.js';

import {RootGroup as RootGroup} from './Search.js';


function getRecentState() {
  return {
    items: RecentStore.getPaths()
  };
}

export default class Recent extends React.Component {
  constructor(props) {
    super(props);

    this.state = getRecentState();
    this._onChange = this._onChange.bind(this);
  }

  componentDidMount() {
    RecentStore.addChangeListener(this._onChange);
    RecentActions.fetch();
  }

  componentWillUnmount() {
    RecentStore.removeChangeListener(this._onChange);
  }

  render() {
    var list = this.state.items.map(function(path) {
      return <RootGroup path={path} key={"rootgroup-"+CollectionStore.pathToKey(path)} />;
    });

    return (
      <div className="collection">
        {list}
      </div>
    );
  }

  _onChange() {
    this.setState(getRecentState());
  }
}
