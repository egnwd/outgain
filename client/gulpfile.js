'use strict';

var gulp = require('gulp');
var merge2 = require('merge2');
var sass = require('gulp-sass');
var typescript = require('gulp-typescript');

var targetDir = './dist';

var tsProject = typescript.createProject('tsconfig.json');

gulp.task('scripts', function () {
    var ts = tsProject.src()
                .pipe(typescript(tsProject)).js;
    var js = gulp.src('./src/**/*.js');

    return merge2(ts, js)
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

gulp.task('all', ['scripts', 'styles', 'html', 'images']);

gulp.task('watch', ['scripts', 'styles', 'html', 'images'], function() {
    gulp.watch('./src/**/*.ts', ['scripts']);
    gulp.watch('./src/**/*.js', ['scripts']);
    gulp.watch('./style/**/*.scss', ['styles']);
    gulp.watch('./html/**/*.html', ['html']);
    gulp.watch('./images/**/*.scss', ['images']);
});

gulp.task('default', ['all']);
