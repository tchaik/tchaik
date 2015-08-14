"use strict";

import React from "react";

import FavouriteStore from "../stores/FavouriteStore.js";
import FavouriteActions from "../actions/FavouriteActions.js";

import {GroupList as GroupList} from "./Collection.js";


function getFavouriteState() {
  return {
    items: FavouriteStore.getPaths(),
  };
}

export default class Favourite extends React.Component {
  constructor(props) {
    super(props);

    this.state = getFavouriteState();
    this._onChange = this._onChange.bind(this);
  }

  componentDidMount() {
    FavouriteStore.addChangeListener(this._onChange);
    FavouriteActions.fetch();
  }

  componentWillUnmount() {
    FavouriteStore.removeChangeListener(this._onChange);
  }

  render() {
    return <GroupList path={["Root"]} list={this.state.items} depth={0} />;
  }

  _onChange() {
    this.setState(getFavouriteState());
  }
}
