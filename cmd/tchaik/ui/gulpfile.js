var gulp = require('gulp');

var sass = require('gulp-sass');
var sourcemaps = require('gulp-sourcemaps');

var paths = {
    sass: {
        src: ['static/sass/**/*.scss'],
        dest: 'static/css'
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

gulp.task('watch', function() {
    gulp.watch(paths.sass.src, ['sass']);
});

gulp.task('default', []);
