"use strict";

import React from "react";

import classNames from "classnames";

import FilterStore from "../stores/FilterStore.js";
import FilterActions from "../actions/FilterActions.js";

import {GroupList as GroupList} from "./Collection.js";

function getFilterState(name) {
  return {
    items: FilterStore.getItems(name),
    current: FilterStore.getCurrentItem(name),
  };
}

export default class Filter extends React.Component {
  constructor(props) {
    super(props);

    this.state = {items: [], current: null};
    this.setCurrent = this.setCurrent.bind(this);
    this._onChange = this._onChange.bind(this);
  }

  componentDidMount() {
    FilterStore.addChangeListener(this._onChange);
    FilterActions.fetchList(this.props.name);
  }

  componentWillReceiveProps(nextProps) {
    this.setState(getFilterState(nextProps.name));
  }

  componentWillUnmount() {
    FilterStore.removeChangeListener(this._onChange);
  }

  setCurrent(itemName) {
    FilterActions.setItem(this.props.name, itemName);
    FilterActions.fetchPaths(this.props.name, itemName);
  }

  render() {
    var items = this.state.items.map(function(item) {
      return <Item item={item}
                   key={item}
               current={this.state.current === item}
            setCurrent={this.setCurrent} />;
    }.bind(this));

    var results = null;
    if (this.state.current !== null) {
      results = <Results filterName={this.props.name} itemName={this.state.current} />;
    }

    return (
      <div className="filter">
        <div className="sidebar">
          <ul>{items}</ul>
        </div>
        <div className="collection">
          {results}
        </div>
      </div>
    );
  }

  _onChange() {
    this.setState(getFilterState(this.props.name));
  }
}


class Item extends React.Component {
  constructor(props) {
    super(props);

    this._onClick = this._onClick.bind(this);
  }

  render() {
    return <li onClick={this._onClick} className={classNames({"selected": this.props.current})}>{this.props.item}</li>;
  }

  _onClick() {
    this.props.setCurrent(this.props.item);
  }
}


class Results extends React.Component {
  constructor(props) {
    super(props);

    this.state = {items: []};
    this._onChange = this._onChange.bind(this);
  }

  componentDidMount() {
    FilterStore.addChangeListener(this._onChange);
    FilterActions.fetchPaths(this.props.filterName, this.props.itemName);
  }

  componentWillUnmount() {
    FilterStore.removeChangeListener(this._onChange);
  }

  render() {
    return <GroupList list={this.state.items} depth={0} path={["Root"]} />;
  }

  _onChange() {
    this.setState({
      items: FilterStore.getPaths(this.props.filterName, FilterStore.getCurrentItem(this.props.filterName)),
    });
  }
}
