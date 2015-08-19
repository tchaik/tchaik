import WebsocketAPI from "../utils/WebsocketAPI.js";

import ChecklistConstants from "../constants/ChecklistConstants.js";


var ChecklistActions = {

  fetch: function() {
    WebsocketAPI.send(ChecklistConstants.FETCH_CHECKLIST);
  },

};

export default ChecklistActions;
