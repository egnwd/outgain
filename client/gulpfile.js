'use strict';

var gulp        = require('gulp'),
    browserify  = require('browserify'),
    tsify       = require('tsify'),
    source      = require('vinyl-source-stream'),
    sass        = require('gulp-sass'),
    gulpTypings = require('gulp-typings'),
    es          = require('event-stream');

var targetDir = __dirname + '/dist';

gulp.task('scripts', ['typings'], function () {
    return browserify()
        .add(__dirname + '/typings/index.d.ts')
        .add(__dirname + '/src/main.ts')
        .plugin(tsify)
        .bundle()
        .on('error', logError)
        .pipe(source('bundle.js'))
        .pipe(gulp.dest(targetDir + '/js'));
});

gulp.task('popup', ['typings'], function () {
    return browserify()
        .add(__dirname + '/typings/index.d.ts')
        .add(__dirname + '/src/index.ts')
        .plugin(tsify)
        .bundle()
        .on('error', logError)
        .pipe(source('index.bundle.js'))
        .pipe(gulp.dest(targetDir + '/js'));
});

var logError = function logError(error) {
  process.stderr.write(error + '\n');
};

gulp.task('typings', function(){
    return gulp.src('./typings.json')
        .pipe(gulpTypings());
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

gulp.task('all', ['scripts', 'popup', 'styles', 'html', 'images']);

gulp.task('watch', ['scripts', 'popup', 'styles', 'html', 'images'], function() {
    gulp.watch('./src/**/*.ts', ['popup']);
    gulp.watch('./src/**/*.ts', ['scripts']);
    gulp.watch('./src/**/*.js', ['scripts']);
    gulp.watch('./style/**/*.scss', ['styles']);
    gulp.watch('./html/**/*.html', ['html']);
    gulp.watch('./images/**/*.scss', ['images']);
});

gulp.task('default', ['all']);
