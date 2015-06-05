var path = require("path");
var webpack = require("webpack");

var ExtractTextPlugin = require('extract-text-webpack-plugin');

module.exports = {
  plugins: [
    new webpack.DefinePlugin({
      "process.env": Object.keys(process.env).reduce(function(o, k) {
        o[k] = JSON.stringify(process.env[k]);
        return o;
      }, {}),
    }),
    new ExtractTextPlugin("styles.css"),
  ],

  devtool: "inline-source-map",

  module: {
    loaders: [
      {
        test: /\.js$/,
        exclude: /node_modules/,
        loaders: ["react-hot-loader", "babel-loader?stage=0", "eslint-loader"],
      },
      {
        test: /\.scss$/,
        loader: ExtractTextPlugin.extract(
          "css?sourceMap!sass?outputStyle=expanded&sourceMap&" +
          "includePaths[]=" +
            (path.resolve(__dirname, "./bower_components")) + "&" +
          "includePaths[]=" +
            (path.resolve(__dirname, "./node_modules"))
        ),
      },
    ],
  },
};
