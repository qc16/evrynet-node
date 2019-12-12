#!/bin/sh
echo "------ Buiding Gev-Node docker image ------"
BASEDIR=$(dirname "$0")
git clone --single-branch -b develop --single-branch git@github.com:evrynet-official/evrynet-client.git "$BASEDIR"/project
docker build -f dockerfiles/node/Dockerfile -t kybernetwork/evrynet-node:1.0.1-dev "$BASEDIR"
rm -rf "$BASEDIR"/project
