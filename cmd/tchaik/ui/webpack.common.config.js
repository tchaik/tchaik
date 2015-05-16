var path = require('path');
var webpack = require('webpack');

module.exports = {
  plugins: [
    new webpack.DefinePlugin({
      "process.env": Object.keys(process.env).reduce(function(o, k) {
        o[k] = JSON.stringify(process.env[k]);
        return o;
      }, {})
    }),
  ],

  module: {
    loaders: [
      {
        test: /\.js$/,
        exclude: /node_modules/,
        loaders: ['react-hot-loader', 'babel-loader?stage=0'],
      },
      {
        test: /\.scss$/,
        loader: "style!css!sass?outputStyle=expanded&" +
          "includePaths[]=" +
            (path.resolve(__dirname, "./bower_components")) + "&" +
          "includePaths[]=" +
            (path.resolve(__dirname, "./node_modules"))
      },
    ]
  }
};
