var gulp = require("gulp");

var _ = require("lodash");
var eslint = require("gulp-eslint");
var gutil = require("gulp-util");
var karma = require("gulp-karma");
var uglify = require("gulp-uglify");
var webpack = require("webpack");
var WebpackDevServer = require("webpack-dev-server");

var paths = {
  js: {
    tests: "static/js/src/**/*-test.js",
    entry: "./static/js/src/app.js",
    bundleName: "tchaik.js",
    dest: "static/js/build",
  },
  eslint: {
    src: [
      "*.js",
      "static/js/src/**/*.js",
      "static/js/src/*.js",
    ],
  },
};

gulp.task("webpack", function(done) {
  webpack(
    require("./webpack.config.js"),
    function(err, stats) {
      if(err) {
        throw new gutil.PluginError("webpack", err);
      }
      gutil.log("[webpack]", stats.toString({
        // output options
      }));
      done();
    }
  );
});

gulp.task("eslint", function() {
  return gulp.src(paths.eslint.src)
    .pipe(eslint())
    .pipe(eslint.formatEach())
    .pipe(eslint.failOnError());
});

gulp.task("serve", function() {
  var webpackConfig = require("./webpack.config.js");
  // Webpack requires an absolute path for the dev server
  webpackConfig.output.path = "/";

  // Load the hot module replacement library
  webpackConfig.entry.app.unshift("webpack/hot/dev-server");
  // Load and configure the dev-server client.
  webpackConfig.entry.app.unshift("webpack-dev-server/client?http://localhost:3000");

  // Load the hot module replacement server plugin
  webpackConfig.plugins.push(new webpack.HotModuleReplacementPlugin());

  var compiler = webpack(webpackConfig);
  var server = new WebpackDevServer(compiler, {
    publicPath: "/static/js/build/",
    hot: true,
    stats: { colors: true },
    proxy: {
      "*": "http://localhost:8080",
    },
  });

  server.listen(3000);
});


function setupKarma(options) {
  return gulp.src([
    // Polyfill so we can use react in phantomjs
    "./node_modules/phantomjs-polyfill/bind-polyfill.js",
    "./node_modules/babel-core/browser-polyfill.js",
    // Test files
    paths.js.tests,
  ])
  .pipe(karma(_.assign({
    configFile: "karma.conf.js",
    browsers: ["PhantomJS"],
  }, options)));
}

gulp.task("karma:ci", function() {
  return setupKarma({
    action: "run",
  });
});

gulp.task("karma:dev", function() {
  return setupKarma({
    action: "watch",
  });
});

gulp.task("compress", function() {
  return gulp.src(paths.js.dest + "/" + paths.js.bundleName)
    .pipe(uglify())
    .pipe(gulp.dest(paths.js.dest));
});

gulp.task("default", ["webpack", "eslint"]);
gulp.task("lint", ["eslint"]);
gulp.task("test", ["karma:ci"]);
