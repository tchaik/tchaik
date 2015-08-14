import WebsocketAPI from "../utils/WebsocketAPI.js";

import FavouriteConstants from "../constants/FavouriteConstants.js";


var FavouriteActions = {

  fetch: function() {
    WebsocketAPI.send(FavouriteConstants.FETCH_FAVOURITE);
  },

};

export default FavouriteActions;
