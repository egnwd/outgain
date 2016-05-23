'use strict';

var gulp = require('gulp');
var sass = require('gulp-sass');
var ts = require('gulp-typescript');

var targetDir = './dist';

var tsProject = ts.createProject('tsconfig.json');

gulp.task('scripts', function () {
    return tsProject.src()
        .pipe(ts(tsProject))
        .js.pipe(gulp.dest(targetDir + '/js'));
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

gulp.task('all', ['scripts', 'styles', 'html', 'images']);

gulp.task('watch', ['scripts', 'styles', 'html', 'images'], function() {
    gulp.watch('./src/**/*.ts', ['scripts']);
    gulp.watch('./style/**/*.scss', ['styles']);
    gulp.watch('./html/**/*.html', ['html']);
    gulp.watch('./images/**/*.scss', ['images']);
});

gulp.task('default', ['all']);
