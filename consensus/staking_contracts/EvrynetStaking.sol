pragma solidity 0.5.11;

/**
 * @title Helps contracts guard against reentrancy attacks.
 */
contract ReentrancyGuard {

    /// @dev counter to allow mutex lock with only one SSTORE operation
    uint256 private guardCounter = 1;

    /**
     * @dev Prevents a function from calling itself, directly or indirectly.
     * Calling one `nonReentrant` function from
     * another is not supported. Instead, you can implement a
     * `private` function doing the actual work, and an `external`
     * wrapper marked as `nonReentrant`.
     */
    modifier nonReentrant() {
        guardCounter += 1;
        uint256 localCounter = guardCounter;
        _;
        require(localCounter == guardCounter);
    }
}

contract EvrynetStaking is ReentrancyGuard {

    struct CandidateData {
        bool isCandidate;
        uint totalStake;
        address owner;
        mapping(address => uint) voterStakes;
    }

    mapping(address => CandidateData) public candidateData;
    address[] public candidates;

    uint public startBlock;
    uint public epochPeriod;

    uint public maxValidatorSize;
    uint public minValidatorStake; // min (own) stake to be a validator
    uint public minVoterCap;

    address public admin;

    modifier onlyAdmin {
        require(msg.sender == admin);
        _;
    }

    modifier onlyValidCandidate(address candidate) {
        require(candidateData[candidate].isCandidate == true);
        _;
    }

    modifier onlyValidVoteCap {
        require(msg.value >= minVoterCap);
        _;
    }

    constructor(
        address[] memory _candidates,
        address candidatesOwner,
        uint _epochPeriod,
        uint _maxValidatorSize,
        uint _minValidatorStake,
        uint _minVoteCap) public
    {
        require(_epochPeriod > 0);

        epochPeriod = _epochPeriod;
        maxValidatorSize = _maxValidatorSize;
        minValidatorStake = _minValidatorStake;
        minVoterCap = _minVoteCap;

        candidates = _candidates;
        for(uint i = 0; i < _candidates.length; i++) {
            candidateData[_candidates[i]] = CandidateData({
                isCandidate: true,
                owner: candidatesOwner,
                totalStake: _minValidatorStake
            });
            candidateData[_candidates[i]].voterStakes[candidatesOwner] = _minValidatorStake;
        }

        admin = msg.sender;
        startBlock = block.number;
    }

    function () external payable {}

    function transferAdmin(address newAdmin) onlyAdmin public {
        require(newAdmin != address(0));
        admin = newAdmin;
    }

    function updateMinValidateStake(uint _newCap) onlyAdmin public {
        minValidatorStake = _newCap;
    }

    function updateMinVoteCap(uint _newCap) onlyAdmin public {
        minVoterCap = _newCap;
    }

    event Voted(address voter, address candidate, uint amount);
    function vote(address candidate) payable public onlyValidVoteCap onlyValidCandidate(candidate) {
        candidateData[candidate].totalStake += msg.value;
        candidateData[candidate].voterStakes[msg.sender] += msg.value;
        emit Voted(msg.sender, candidate, msg.value);
    }

    event Unvoted(address voter, address candidate, uint amount);
    function unvote(address candidate, uint amount) nonReentrant public {
        require(candidateData[candidate].voterStakes[msg.sender] >= amount);
        candidateData[candidate].voterStakes[msg.sender] -= amount;
        candidateData[candidate].totalStake -= amount;
        msg.sender.transfer(amount);
        emit Unvoted(msg.sender, candidate, amount);
    }

    event Registered(address candidate, uint stake);
    function register(address _candidate) payable public {
        // not current candidate
        require(candidateData[_candidate].isCandidate == false);
        candidateData[_candidate] = CandidateData({
           owner: msg.sender,
           isCandidate: true,
           totalStake: msg.value
        });
        candidateData[_candidate].voterStakes[msg.sender] = msg.value;
        candidates.push(_candidate);
        emit Registered(_candidate, msg.value);
    }

    event Resigned(address _candidate);
    function resign(address _candidate) nonReentrant public {
        // only owner can resign
        require(candidateData[_candidate].isCandidate == true);
        require(candidateData[_candidate].owner == msg.sender);

        for(uint i = 0; i < candidates.length; i++) {
            if (candidates[i] == _candidate) {
                candidates[i] = candidates[candidates.length - 1];
                delete candidates[candidates.length - 1];
                candidates.length--;
                break;
            }
        }

        candidateData[_candidate].isCandidate = false;
        // withdraw stake of owner
        uint ownerStake = candidateData[_candidate].voterStakes[msg.sender];
        candidateData[_candidate].voterStakes[msg.sender] = 0;
        if (ownerStake > 0) {
            msg.sender.transfer(ownerStake);
        }
        emit Resigned(_candidate);
    }

    function updateMaxValidatorSize(uint newMaxValidatorSize) onlyAdmin public {
        maxValidatorSize = newMaxValidatorSize;
    }

    function getCurrentEpoch() public view returns(uint) {
        return (block.number - startBlock) / epochPeriod;
    }

    function getAllCandidates() public view returns(address[] memory _candidates) {
        _candidates = candidates;
    }

    function getListCandidates()
        public view returns(address[] memory _candidates, uint[] memory _stakes, uint32 _maxValSize, uint32 _epochSize)
    {
        _maxValSize = uint32(maxValidatorSize);
        _epochSize = uint32(getCurrentEpoch());

        uint validCandiateCount;
        // only count candidate with his own stake >= minValidatorStake
        for(uint i = 0; i < candidates.length; i++) {
            address owner = candidateData[candidates[i]].owner;
            if (candidateData[candidates[i]].voterStakes[owner] >= minValidatorStake) {
                validCandiateCount++;
            }
        }

        _candidates = new address[](validCandiateCount);
        _stakes = new uint[](validCandiateCount);

        uint index = 0;
        for(uint i = 0; i < candidates.length; i++) {
            address owner = candidateData[candidates[i]].owner;
            if (candidateData[candidates[i]].voterStakes[owner] >= minValidatorStake) {
                _candidates[index] = candidates[i];
                _stakes[index] = candidateData[candidates[i]].totalStake;
                index++;
            }
        }
    }

    function epochSize() public view returns(uint32) {
        return uint32(getCurrentEpoch());
    }

    function getCandidateData(address _candidate)
        public view
        returns(bool _isCandidate, address _owner, uint _totalStake)
    {
        _isCandidate = candidateData[_candidate].isCandidate;
        _owner = candidateData[_candidate].owner;
        _totalStake = candidateData[_candidate].totalStake;
    }

    function getVoterStake(address _candidate, address _voter) public view returns(uint) {
        return candidateData[_candidate].voterStakes[_voter];
    }
}