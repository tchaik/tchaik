import WebsocketAPI from "../utils/WebsocketAPI.js";

import PathListConstants from "../constants/PathListConstants.js";


var PathListActions = {

  fetch: function(name) {
    WebsocketAPI.send(PathListConstants.FETCH_PATHLIST, {name: name});
  },

};

export default PathListActions;
