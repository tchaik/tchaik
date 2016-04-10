var webpack = require("webpack");

var ExtractTextPlugin = require("extract-text-webpack-plugin");

var autoprefixer = require("autoprefixer-core");
var postcssNested = require("postcss-nested");
var postcssSassyMixins = require("postcss-sassy-mixins");
var postcssSimpleVars = require("postcss-simple-vars");
var postcssImport = require("postcss-import");
var postcssColorFunction = require("postcss-color-function");
var postcssColorRGBAFallback = require("postcss-color-rgba-fallback");

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
        loaders: ["react-hot", "babel"],
      },
      {
        test: /\.css$/,
        loader: ExtractTextPlugin.extract(
          "style-loader", "css-loader!postcss-loader"
        ),
      },
    ],
  },

  postcss: function() {
    return [
      postcssImport({
        onImport: function (files) {
          files.forEach(this.addDependency);
        }.bind(this),
      }),
      postcssSassyMixins,
      postcssNested,
      postcssSimpleVars,
      postcssColorFunction,
      postcssColorRGBAFallback,
      autoprefixer,
    ];
  },
};
