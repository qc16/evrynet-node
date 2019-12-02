#!/bin/sh
echo "Stop node $NODE_ID"
# shellcheck disable=SC2046
kill -INT $(pgrep -f ./gev)