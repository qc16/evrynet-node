go build ./cmd/gev
mkdir ~/evrynet
cp ./tests/provider_logic_test/test_genesis.json ~/evrynet/genesis.json
./gev --networkid 15 --datadir ~/evrynet init ~/evrynet/genesis.json
echo "------------------NOTICE------------------"
echo "PLEASE PUT 123 AS THE PASSWORD WHEN PROMPT"
echo "------------------======------------------"
./gev --datadir ~/evrynet account import tests/provider_logic_test/etherbase_pk.txt
