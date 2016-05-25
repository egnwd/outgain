'use strict';

var gulp = require('gulp');
var browserify = require('browserify');
var tsify = require('tsify');
var source = require('vinyl-source-stream');
var sass = require('gulp-sass');


var targetDir = __dirname + '/dist';

gulp.task('scripts', function () {
    return browserify(__dirname + '/src/main.ts')
        .plugin(tsify)
        .bundle()
        .pipe(source('bundle.js'))
        .pipe(gulp.dest(targetDir + '/js'));
});

gulp.task('sweetalert', function () {
  return gulp.src('./src/sweetalert.min.js')
      .pipe(gulp.dest(targetDir + '/js'));
});

gulp.task('styles', function () {
    return gulp.src('./style/**/*.scss')
        .pipe(sass().on('error', sass.logError))
        .pipe(gulp.dest(targetDir + '/css'));
});

gulp.task('html', function () {
    return gulp.src('./html/*.html')
        .pipe(gulp.dest(targetDir));
});

gulp.task('images', function () {
    return gulp.src('./images/*.{png,jpg,svg}')
        .pipe(gulp.dest(targetDir + '/images'));
});

gulp.task('all', ['scripts', 'styles', 'html', 'images', 'sweetalert']);

gulp.task('watch', ['scripts', 'styles', 'html', 'images'], function() {
    gulp.watch('./src/**/*.ts', ['scripts']);
    gulp.watch('./src/**/*.js', ['scripts']);
    gulp.watch('./style/**/*.scss', ['styles']);
    gulp.watch('./html/**/*.html', ['html']);
    gulp.watch('./images/**/*.scss', ['images']);
});

gulp.task('demo', function () {
  return gulp.src('./src/demo.js')
      .pipe(gulp.dest(targetDir + '/js'));
});

gulp.task('default', ['all']);
