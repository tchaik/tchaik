var gulp = require('gulp');

var _ = require('lodash');
var gutil = require('gulp-util');
var jshint = require('gulp-jshint');
var webpack = require('webpack');
var webpackDevServer = require('webpack-dev-server');

var paths = {
  sass: {
    src: ['static/sass/**/*.scss'],
    dest: 'static/css'
  },
  js: {
    entry: './static/js/src/app.js',
    bundleName: 'tchaik.js',
    dest: 'static/js/build'
  },
  jshint: {
    src: [
      'package.json',
      'gulpfile.js',
      'static/js/src/app.js',
      'static/js/src/actions/*.js',
      'static/js/src/stores/*.js',
      'static/js/src/constants/*.js',
      'static/js/src/utils/*.js'
    ]
  }
};

gulp.task('webpack', function(done) {
  webpack(
    require('./webpack.config.js'),
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

var jshintConfig = {
  esnext: true,
  browser: true,
  devel: true,
  curly: true,
  undef: true,
  unused: true,
  node: true,
  newcap: false,
};

gulp.task('jshint', function() {
  gulp.src(paths.jshint.src)
    .on('error',gutil.log.bind(gutil, 'JSHint Error'))
    .pipe(jshint(jshintConfig))
    .pipe(jshint.reporter('jshint-stylish'));
});

gulp.task('jshint:jsx', function() {
  return gulp.src(['static/js/src/components/*.js'])
    .pipe(jshint({
      linter: require('jshint-jsx').JSXHINT,
      esnext: true,
      browser: true,
      devel: true,
      jquery: true,
      curly: true,
      undef: true,
      unused: true,
      node: true,
      newcap: false
    }))
    .pipe(jshint.reporter('jshint-stylish'))
    .pipe(jshint.reporter('fail'));
});

gulp.task('serve', function() {
  var webpackConfig = require('./webpack.config.js');
  // Webpack requires an absolute path for the dev server
  webpackConfig.output.path = '/';

  // Load the hot module replacement library
  webpackConfig.entry.app.unshift('webpack/hot/dev-server');
  // Load and configure the dev-server client.
  webpackConfig.entry.app.unshift('webpack-dev-server/client?http://localhost:3000');

  // Load the hot module replacement server plugin
  webpackConfig.plugins.push(new webpack.HotModuleReplacementPlugin());

  var compiler = webpack(webpackConfig);
  var server = new webpackDevServer(compiler, {
    publicPath: '/static/js/build/',
    hot: true,
    stats: { colors: true },
    proxy: {
      '*': 'http://localhost:8080'
    }
  });

  server.listen(3000);
});

gulp.task('default', ['webpack', 'jshint:jsx']);
gulp.task('lint', ['jshint', 'jshint:jsx']);
