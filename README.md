# outgain
AI based evolution simulation

## Dependencies
- Go
- Node
- Gulp (`npm install -g gulp`)

Run `npm install` inside the `node_modules` directory.

## Building
Running `./build_all.sh` from the root of the project will build everything.

### Client
From the `client` directory, run `gulp`.
Alternatively, `gulp watch` will watch for changes to the source and rebuild
automatically.

### Server
From the `server` directory, run `go build`.

## Running
To run the outgain server, run the following from the root of the project :

```shell
./server/server
```

This will listen on port 8080 by default, use the `PORT` environment variable
to override.

It will serve the files for the client from the `client/dist` directory.

## Deploying
Pushing to master or merging a pull request into it will build both the server and
the client on Circle CI.
If the build is succesful, it will be deployed automatically to Heroku.

Only the files needed to run the server are pushed to Heroku.
Check the `build_slug.sh` if you need to add some files.

## Manual deployment
Unless you have a good reason to, you shouldn't do this, but rely on the CI to deploy
automatically.

```
./build_slug.sh app
tar czvf slug.tgz ./app
HEROKU_OAUTH_TOKEN="<CHANGEME>" ./deploy.rb outgain slug.tgz
```
