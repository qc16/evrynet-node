## build 
Navigate to `./cmd/indexGenerator` and run the command below:

```
$ go build .
```

## generate index of state variables
- copy `indexGenerator` and `getStorageLayput.sh` files then place its to the directory where there are the contract solidity files
- run command below:

```
$ ./indexGenerator
```
- Then we have got a data file's name is `indexGenerator.json`. When run a node let use this file to set data for index state variables with syntax `./gev --sc.index-generator-path [PATH_OF_indexGenerator.json]` (see more via a command `./gev -h` with the flag `--sc.index-generator-path`)