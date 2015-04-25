var gulp = require('gulp');
var _ = require('lodash');

var browserify = require('browserify');
var buffer = require('vinyl-buffer');
var envify = require('envify');
var gutil = require('gulp-util');
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
            .pipe(gulp.dest(paths.js.dest));
    };

    bundler.on('update', rebundle);
    return rebundle();
}

gulp.task('js', function() {
    return bundle(false);
});

gulp.task('watch', function() {
    gulp.watch(paths.sass.src, ['sass']);
    bundle(true);
});

gulp.task('default', []);
