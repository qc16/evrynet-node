#!/bin/bash

echo "***********************************************"
echo "***********************************************"
echo "    INIT YOUR NODE MUST ENTER SOME INFOR BELOW"
listAcc=$(./gev account list)
hasAccount=0
address=''
password=''
datadir=~/evrynet

if [[ -z "$listAcc" ]]; then
    echo "    In the localhost there are not accounts, you must create or import an account"
else
    echo "    In the localhost there are accounts here: $listAcc "
    hasAccount=1
fi
echo "***********************************************"

read -p "    Do you have an account? (y|n): " yn
if [[ $yn = 'y' ]]; then
    echo    "    starting import your account"
    read -p "    Enter your private key: "  privateKey
    read -p "    Enter your passphrase for unlock: " -s password
    result=$(./gev account import --datadir "$datadir" --password <(echo "$password") <(echo "$privateKey"))
    result="$(cut -d'{' -f2 <<<"$result")"
    address="$(cut -d'}' -f1 <<<"$result")"
else
    read -p "    Enter your passphrase for new account: " -s password
    ./gev account new --datadir "$datadir" --password <(echo "$password")
    read -p "    Please enter the address have been created above to unlock and run your node: " address
fi

echo "***********************************************"
read -p "    Do you want to set the "$address" as coinbase? (y|n): " yn
echo "***********************************************"

if [[ $yn = 'y' ]]; then
    ./gev --datadir "$datadir" --port 30311 --rpc --rpcaddr 'localhost' --networkid 15 --gasprice '0' --password <(echo "$password") --unlock "$address" --etherbase "$address" --mine --allow-insecure-unlock
else
    ./gev --datadir "$datadir" --port 30311 --rpc --rpcaddr 'localhost' --networkid 15 --gasprice '0' --password <(echo "$password") --unlock "$address" --mine --allow-insecure-unlock
fi
