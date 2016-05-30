#!/usr/bin/env bash

# Executed on Circle once the build is succesful
# It prepares the slug by copying the necessary build artifact
# into the specified folder
#
# It creates an executable start.sh file at the root of the slug,
# which starts the server
#
# The first argument is the target folder
# This must be run from the root of the project, after building
# both the client and the server

set -eux

TARGET_DIR="$1"

mkdir -p "$TARGET_DIR"

cp -r client/dist "$TARGET_DIR/static"
cp server/server "$TARGET_DIR"
cp sandbox/sandbox "$TARGET_DIR"

cat > "$TARGET_DIR/start.sh" <<EOF
#!/usr/bin/env bash
set -eux
exec ./server \
    -redirect-plain-http \
    -static-dir=./static \
    -sandbox=trace -sandbox-bin=./sandbox
EOF
chmod +x "$TARGET_DIR/start.sh"
