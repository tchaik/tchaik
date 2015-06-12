"use strict";

import React from "react";

import SearchStore from "../stores/SearchStore.js";

import {GroupList as GroupList} from "./Collection.js";

import Icon from "./Icon.js";

export class Search extends React.Component {
  render() {
    return (
      <div className="collection">
        <Results />
      </div>
    );
  }
}


function getResultsState() {
  return {results: SearchStore.getResults()};
}

class Results extends React.Component {
  constructor(props) {
    super(props);

    this.state = getResultsState();
    this._onChange = this._onChange.bind(this);
  }

  componentDidMount() {
    SearchStore.addChangeListener(this._onChange);
  }

  componentWillUnmount() {
    SearchStore.removeChangeListener(this._onChange);
  }

  render() {
    var list = this.state.results;
    if (list.length === 0) {
      return (
        <div className="collection">
          <div className="no-results"><Icon icon="headphones" />No results found</div>
        </div>
      );
    }
    return <GroupList path={["Root"]} list={list} depth={0} />;
  }

  _onChange() {
    this.setState(getResultsState());
  }
}
