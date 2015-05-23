var assign = require("object-assign");
var config = require("./webpack.common.config.js");

module.exports = assign({}, config, {
  entry: {
    app: ["./static/js/src/app.js"],
  },

  output: {
    path: "static/js/build/",
    pathInfo: true,
    publicPath: "/static/js/build/",
    filename: "tchaik.js",
  },
});
