## build 
Navigate to `./cmd/indexGenerator` and run the command below:

```
$ go build .
```

## generate index of state variables

Copy `indexGenerator` and `getStorageLayput.sh` files then place its to the directory where there are the contract solidity files

The currently, we are using the configuration below to generate the index of storage layout:

``` js
{
  "language": "Solidity",
  "sources": {
    "EvrynetStaking.sol": {
      "urls": [
        "EvrynetStaking.sol"
      ]
    }
  },
  "settings": {
    "optimizer": {
      "enabled": true,
      "runs": 200
    },
    "evmVersion": "homestead",
    "outputSelection": {
      "*": {
        "*": [
          "storageLayout"
        ]
      }
    }
  }
}
```

and using the command below:

```

echo '{"language":"Solidity","sources":{"EvrynetStaking.sol":{"urls":["EvrynetStaking.sol"]}},"settings":{"optimizer":{"enabled":true,"runs":200},"evmVersion":"homestead","outputSelection":{"*":{"*":["storageLayout"]}}}}' | solc --optimize --standard-json --allow-paths *,

```

Run command below:

```
$ ./indexGenerator run -h

NAME:
   indexGenerator run - use run command to generates the storage layout of contract's state variables

USAGE:
   indexGenerator run [command options] [arguments...]

DESCRIPTION:
   this tool to generates the storage layout of contract's state variables

OPTIONS:
   --shfilepath value  The path of file to generates storage layout (there are commands to generates storage data layout in this file) (default: "./getStorageLayput.sh")

```

- Then we have got a data file's name is `indexGenerator.json`. When run a node let use this file to set data for index state variables with syntax `./gev --tendermint.index-generator-path [PATH_OF_indexGenerator.json]` (see more via a command `./gev -h` with the flag `--tendermint.index-generator-path`)