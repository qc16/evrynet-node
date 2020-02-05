pragma solidity ^0.5.11;

contract Validator {
	address[] public validators;

	constructor (address[] memory _validators) public {
		for (uint256 i = 0; i < _validators.length; i++) {
			validators.push(_validators[i]);
		}
	}

	function getValidators() public view returns(address[] memory) {
		return validators;
	}
}

// NOTICE: 
// to generate file validator.go let do
// - go build ./cmd/abigen
// - run command ./abigen  --sol ./consensus/staking/validator/contract/Validator.sol --pkg contract --out "./consensus/staking/validator/contract/validator.go"