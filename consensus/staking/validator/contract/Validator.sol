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