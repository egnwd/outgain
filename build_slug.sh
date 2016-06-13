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
cp runner/target/release/runner "$TARGET_DIR"
cp default_ai.rb "$TARGET_DIR"
cp bot_ai.rb "$TARGET_DIR"
cp "/usr/lib/x86_64-linux-gnu/libseccomp.so.2" "$TARGET_DIR"

cat > "$TARGET_DIR/start.sh" <<EOF
#!/usr/bin/env bash
set -eux
export LD_LIBRARY_PATH=.
exec ./server \
    -redirect-plain-http \
    -static-dir=./static \
    -sandbox=trace \
    -runner-bin=./runner \
    -default-ai=./default_ai.rb \
    -bot-ai=./bot_ai.rb
EOF
chmod +x "$TARGET_DIR/start.sh"
