import React from "react";

import ContainerConstants from "../constants/ContainerConstants.js";
import ContainerStore from "../stores/ContainerStore.js";
import SearchStore from "../stores/SearchStore.js";

import {RootCollection as RootCollection} from "./Collection.js";
import {Search as Search} from "./Search.js";
import Covers from "./Covers.js";
import Filter from "./Filter.js";
import PathList from "./PathList.js";
import Settings from "./Settings.js";
import Retro from "./Retro.js";

function getContainerState() {
  return {
    mode: ContainerStore.getMode(),
  };
}

export default class Container extends React.Component {
  constructor(props) {
    super(props);

    this.state = getContainerState();
    this._onChange = this._onChange.bind(this);
    this._onSearch = this._onSearch.bind(this);
  }

  componentDidMount() {
    ContainerStore.addChangeListener(this._onChange);
    SearchStore.addChangeListener(this._onSearch);
  }

  componentWillUnmount() {
    ContainerStore.removeChangeListener(this._onChange);
    SearchStore.removeChangeListener(this._onSearch);
  }

  render() {
    var content = null;
    switch (this.state.mode) {
      case ContainerConstants.ARTISTS:
        content = <Filter name="Artist" />;
        break;

      case ContainerConstants.COMPOSERS:
        content = <Filter name="Composer" />;
        break;

      case ContainerConstants.SEARCH:
        content = <Search />;
        break;

      case ContainerConstants.COVERS:
        content = <Covers />;
        break;

      case ContainerConstants.RECENT:
        content = <PathList name="recent" />;
        break;

      case ContainerConstants.FAVOURITE:
        content = <PathList name="favourite" />;
        break;

      case ContainerConstants.CHECKLIST:
        content = <PathList name="checklist" />;
        break;

      case ContainerConstants.SETTINGS:
        content = <Settings />;
        break;

      case ContainerConstants.RETRO:
        content = <Retro />;
        break;

      case ContainerConstants.ALL:
      default:
        content = <RootCollection />;
        break;
    }

    return <div>{content}</div>;
  }

  _onChange() {
    this.setState(getContainerState());
  }

  _onSearch() {
    if (SearchStore.getInput() !== "") {
      this.setState({mode: ContainerConstants.SEARCH});
    }
  }
}
