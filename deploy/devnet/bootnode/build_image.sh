#!/bin/sh
echo "------ Buiding Gev-Bootnode docker image ------"
BASEDIR=$(dirname "$0")
git clone --single-branch -b develop --single-branch git@github.com:evrynet-official/evrynet-client.git "$BASEDIR"/project
docker build -t kybernetwork/evrynet-bootnode:1.0.1-dev "$BASEDIR"
rm -rf "$BASEDIR"/project