var _ = require('lodash');
var RewirePlugin = require('rewire-webpack');

var config = _.clone(require('./webpack.common.config.js'));

config.plugins = config.plugins || [];
config.plugins.push(new RewirePlugin());

module.exports = config;
