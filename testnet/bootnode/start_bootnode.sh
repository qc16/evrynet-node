echo "Running Bootnode 1: $NODE_KEY_HEX"
ls -la
# shellcheck disable=SC2230
which bootnode
bootnode -nodekeyhex "$NODE_KEY_HEX"