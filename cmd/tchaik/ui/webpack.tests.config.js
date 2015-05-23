var _ = require("lodash");
var RewirePlugin = require("rewire-webpack");

var config = _.clone(require("./webpack.common.config.js"));

config.plugins = config.plugins || [];
config.plugins.push(new RewirePlugin());

// Append the babel-plugin-rewrite parameter to the babel-loader plugin url
config.module.loaders.forEach(function(loader) {
  for (var i in loader.loaders) {
    if (loader.loaders[i].indexOf("babel-loader") === 0) {
      var originalValue = loader.loaders[i];
      loader.loaders[i] = originalValue + "&plugins=babel-plugin-rewire";
    }
  }
});

module.exports = config;
