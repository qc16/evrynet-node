#!/bin/sh

for i in 1 2 3
do
  sudo docker exec -it gev-node-"$i" /bin/sh ./stop_node.sh
  sleep 3
done
