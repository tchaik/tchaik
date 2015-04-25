var gulp = require('gulp');
var _ = require('lodash');

var browserify = require('browserify');
var buffer = require('vinyl-buffer');
var envify = require('envify');
var gutil = require('gulp-util');
var notify = require('gulp-notify');
var jshint = require('gulp-jshint');
var reactify = require('reactify');
var sass = require('gulp-sass');
var source = require('vinyl-source-stream');
var sourcemaps = require('gulp-sourcemaps');
var watchify = require('watchify');

var paths = {
    sass: {
        src: ['static/sass/**/*.scss'],
        dest: 'static/css'
    },
    js: {
        entry: './static/js/src/app.js',
        bundleName: 'tchaik.js',
        dest: 'static/js/build'
    }
};

gulp.task('sass', function() {
    return gulp.src(paths.sass.src)
               .pipe(sourcemaps.init())
               .pipe(sass())
               .pipe(sourcemaps.write(
                   paths.sass.dest,
                   {sourceRoot: '/static/sass'}
               ))
               .pipe(gulp.dest(paths.sass.dest));
});

function bundle(watch) {
    var bundler, rebundle;

    bundler = browserify(paths.js.entry, {
        debug: true,
        cache: {},
        packageCache: {},
        fullPaths: !watch
    });

    if (watch) {
        bundler = watchify(bundler);
    }

    bundler.transform(reactify);
    bundler.transform(envify);

    rebundle = function() {
        return bundler.bundle()
            .on('error', gutil.log.bind(gutil, 'Browserify Error'))
            .pipe(source(paths.js.bundleName))
            .pipe(buffer())
            .pipe(sourcemaps.init({loadMaps: true}))
            .pipe(sourcemaps.write(
                './',
                {sourceRoot: '/'}
            ))
            .pipe(gulp.dest(paths.js.dest))
            .pipe(notify({message: function() { gutil.log("Built JS"); }, onLast: true}));
    };

    bundler.on('update', rebundle);
    return rebundle();
}

gulp.task('js', function() {
    return bundle(false);
});

gulp.task('jshint:jsx', function() {
    return gulp.src(['static/js/src/**/*.js'])
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

gulp.task('watch', function() {
    gulp.watch(paths.sass.src, ['sass']);
    bundle(true);
});

gulp.task('default', ['sass', 'js', 'jshint:jsx']);
