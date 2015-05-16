var gulp = require('gulp');

var _ = require('lodash');
var gutil = require('gulp-util');
var jshint = require('gulp-jshint');
var karma = require('gulp-karma');
var webpack = require('webpack');
var webpackDevServer = require('webpack-dev-server');

var paths = {
  sass: {
    src: ['static/sass/**/*.scss'],
    dest: 'static/css'
  },
  js: {
    tests: 'static/js/src/**/*-test.js',
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
  globals: {
    'beforeEach': false,
    'describe': false,
    'expect': false,
    'it': false,
  }
};

gulp.task('jshint', function() {
  gulp.src(paths.jshint.src)
    .on('error',gutil.log.bind(gutil, 'JSHint Error'))
    .pipe(jshint(jshintConfig))
    .pipe(jshint.reporter('jshint-stylish'));
});

gulp.task('jshint:jsx', function() {
  var config = _.assign({}, jshintConfig, {
    linter: require('jshint-jsx').JSXHINT,
  });

  return gulp.src(['static/js/src/components/*.js'])
    .pipe(jshint(config))
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


function setupKarma(options) {
  return gulp.src([
    // Polyfill so we can use react in phantomjs
    './node_modules/phantomjs-polyfill/bind-polyfill.js',
    './node_modules/babel-core/browser-polyfill.js',
    // Test files
    paths.js.tests,
  ])
  .pipe(karma(_.assign({
    configFile: 'karma.conf.js',
    browsers: ['PhantomJS'],
  }, options)));
}

gulp.task('karma:ci', function() {
  return setupKarma({
    action: 'run',
  });
});

gulp.task('karma:dev', function() {
  return setupKarma({
    action: 'watch',
  });
});

gulp.task('default', ['webpack', 'jshint:jsx']);
gulp.task('lint', ['jshint', 'jshint:jsx']);
gulp.task('test', ['karma:ci']);
