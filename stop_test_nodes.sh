#!/bin/sh
echo "--- Stop all test nodes ..."
kill -9 $(pgrep -f ./gev)
for i in 1 2 3 4
do
  lsof -ti:3030"$i" | xargs kill
done
