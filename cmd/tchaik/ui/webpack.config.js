var assign = require("object-assign");
var config = require("./webpack.common.config.js");

module.exports = assign({}, config, {
  entry: {
    app: ["./js/src/app.js"],
  },

  output: {
    path: "js/build/",
    pathInfo: true,
    publicPath: "/js/build/",
    filename: "tchaik.js",
  },
});
