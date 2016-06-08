'use strict';

require("any-promise/register/bluebird")
require("util").debuglog = require("debuglog")

var gulp = require('gulp');
var browserify = require('browserify');
var watchify = require('watchify');
var tsify = require('tsify');
var source = require('vinyl-source-stream');
var sass = require('gulp-sass');
var typings = require('gulp-typings');
var gutil = require('gulp-util');

var targetDir = __dirname + '/dist';

function bundle(main, out, watch) {
    var b = browserify();
    if (watch) {
        b = watchify(b)
    }

    b.add(__dirname + '/typings/index.d.ts')
     .add(main)
     .plugin(tsify, { "jsx": "react" })
     .on('update', function() {
         rebundle();
     })
     .on('log', gutil.log);

    function rebundle() {
        gutil.log(gutil.colors.yellow('Bundling ') + gutil.colors.blue(main));
        return b.bundle()
                .on('error', function(err) {
                    gutil.log(gutil.colors.red(err.name) + ': ' + gutil.colors.blue(err.message));
                })
                .pipe(source(out))
                .pipe(gulp.dest(targetDir + '/js'))
    }

    return rebundle()
}

function scripts(watch) {
    return [
        bundle(__dirname + '/src/main.ts', 'bundle.js', watch),
        bundle(__dirname + '/src/index.ts', 'index.bundle.js', watch),
        bundle(__dirname + '/src/lobbies.ts', 'lobbies.bundle.js', watch),
    ];
}

gulp.task('leaderboard', ['typings'], function () {
    return browserify()
        .add(__dirname + '/typings/index.d.ts')
        .add(__dirname + '/src/leaderboard.ts')
        .plugin(tsify)
        .bundle()
        .on('error', logError)
        .pipe(source('leaderboard.bundle.js'))
        .pipe(gulp.dest(targetDir + '/js'));
});

var logError = function logError(error) {
  process.stderr.write(error + '\n');
};

gulp.task('typings', function(){
    return gulp.src('./typings.json')
        .pipe(typings());
});

gulp.task('scripts', ['typings'], function () {
    return scripts();
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

gulp.task('all', ['scripts', 'popup', 'lobbies', 'styles', 'html', 'images', 'leaderboard']);

gulp.task('watch', ['scripts', 'popup', 'lobbies', 'styles', 'html', 'images', 'leaderboard'], function() {
    gulp.watch('./src/**/*.ts', ['popup']);
    gulp.watch('./src/**/lobbies.ts', ['lobbies']);
    gulp.watch('./src/**/leaderboard.ts',['leaderboard'])
    gulp.watch('./src/**/*.ts', ['scripts']);
    gulp.watch('./style/**/*.scss', ['styles']);
    gulp.watch('./html/**/*.html', ['html']);
    gulp.watch('./images/**/*.scss', ['images']);

    return scripts(true);
});

gulp.task('default', ['all']);
