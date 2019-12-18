#!/bin/sh
echo "Stop node $ID"
# shellcheck disable=SC2046
kill -INT $(pgrep -f ./gev)
