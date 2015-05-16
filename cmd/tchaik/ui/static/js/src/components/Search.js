'use strict';

import React from 'react/addons';

import SearchStore from '../stores/SearchStore.js';
import SearchActions from '../actions/SearchActions.js';

import {Group as Group} from './Collection.js';
import CollectionStore from '../stores/CollectionStore.js';
import CollectionActions from '../actions/CollectionActions.js';


export class Search extends React.Component {
  componentDidMount() {
    var input = SearchStore.getInput();
    if (input && input !== "") {
      SearchActions.search(input);
      React.findDOMNode(this.refs.input).value = input;
    }
  }

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
    var list = this.state.results.map(function(path) {
      return <RootGroup path={path} key={CollectionStore.pathToKey(path)} />;
    });
    return (
      <div className="collection">
        {list}
      </div>
    );
  }

  _onChange() {
    this.setState(getResultsState());
  }
}


function getItem(path) {
  var c = CollectionStore.getCollection(path);
  if (c === undefined) {
    return null;
  }
  return c;
}

function getRootGroupState(props) {
  return {item: getItem(props.path)};
}

export class RootGroup extends React.Component {
  constructor(props) {
    super(props);

    this.state = getRootGroupState(this.props);
    this._onChange = this._onChange.bind(this);
  }

  componentDidMount() {
    CollectionStore.addChangeListener(this._onChange);
    CollectionActions.fetch(this.props.path);
  }

  componentWillUnmount() {
    CollectionStore.removeChangeListener(this._onChange);
  }

  render() {
    if (this.state.item === null) {
      return null;
    }

    return (
      <Group item={this.state.item} path={this.props.path} depth={1} />
    );
  }

  _onChange(keyPath) {
    if (CollectionStore.pathToKey(this.props.path) === keyPath) {
       this.setState(getRootGroupState(this.props));
    }
  }
}

RootGroup.propTypes = {
  path: React.PropTypes.array.isRequired,
};
