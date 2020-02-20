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

/**
 * Math operations with safety checks
 */
library SafeMath {
    function mul(uint a, uint b) internal pure returns (uint) {
        uint c = a * b;
        require(a == 0 || c / a == b);
        return c;
    }

    function div(uint a, uint b) internal pure returns (uint) {
        require(b > 0);
        uint c = a / b;
        require(a == b * c + a % b);
        return c;
    }

    function sub(uint a, uint b) internal pure returns (uint) {
        require(b <= a);
        return a - b;
    }

    function add(uint a, uint b) internal pure returns (uint) {
        uint c = a + b;
        require(c >= a);
        return c;
    }

    function max(uint a, uint b) internal pure returns (uint) {
        return a >= b ? a : b;
    }

    function min(uint a, uint b) internal pure returns (uint) {
        return a < b ? a : b;
    }
}


contract EvrynetStaking is ReentrancyGuard {

    using SafeMath for uint;

    // maximum number of candidates
    uint constant internal MAX_CANDIDATES = 128;
    // 2 epochs
    uint constant internal STAKER_LOCKING_PERIOD = 2;

    struct CandidateData {
        bool isCandidate;
        // total stakes at each epoch
        mapping(uint => uint) totalStakes;
        uint latestTotalStakes;
        address owner;
        // voter's stakes for each epoch
        mapping(address => mapping(uint => uint)) voterStakes;
        mapping(address => uint) latestVoterStakes;
    }

    mapping(address => CandidateData) public candidateData;
    address[] public candidates;
    address[] public initCandidates;

    // cap to withdraw for a staker at an epoch
    mapping(address => mapping(uint => uint)) public withdrawalCap;

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

    modifier onlyNotCandidate(address candidate) {
        require(candidateData[candidate].isCandidate == false);
        _;
    }

    modifier onlyValidVoteCap {
        require(msg.value >= minVoterCap);
        _;
    }

    // note: this list candidates should be the validators for epoch
    // other validators should be added after deployed
    constructor(
        address[] memory _candidates,
        address candidatesOwner,
        uint _epochPeriod,
        uint _maxValidatorSize,
        uint _minValidatorStake,
        uint _minVoteCap,
        address _admin) public
    {
        require(_epochPeriod > 0);

        epochPeriod = _epochPeriod;
        maxValidatorSize = _maxValidatorSize;
        minValidatorStake = _minValidatorStake;
        minVoterCap = _minVoteCap;

        require(_maxValidatorSize >= _candidates.length);

        candidates = _candidates;
        for(uint i = 0; i < _candidates.length; i++) {
            candidateData[_candidates[i]] = CandidateData({
                isCandidate: true,
                owner: candidatesOwner,
                latestTotalStakes: _minValidatorStake
            });
            candidateData[_candidates[i]].voterStakes[candidatesOwner][0] = _minValidatorStake;
            candidateData[_candidates[i]].latestVoterStakes[candidatesOwner] = _minValidatorStake;
            candidateData[_candidates[i]].totalStakes[0] = _minValidatorStake;
        }

        initCandidates = _candidates;

        admin = _admin;
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

        uint amount = msg.value;
        address voter = msg.sender;

        candidateData[candidate].latestTotalStakes = candidateData[candidate].latestTotalStakes.add(amount);

        uint curEpoch = getCurrentEpoch();
        candidateData[candidate].totalStakes[curEpoch] = candidateData[candidate].totalStakes[curEpoch].add(amount);

        candidateData[candidate].latestVoterStakes[voter] = candidateData[candidate].latestVoterStakes[voter].add(amount);

        candidateData[candidate].voterStakes[voter][curEpoch] = candidateData[candidate].latestVoterStakes[voter];
        emit Voted(voter, candidate, amount);
    }

    event Unvoted(address voter, address candidate, uint amount);
    function unvote(address candidate, uint amount) nonReentrant public {
        require(amount > 0, "withdraw amount must be positive");

        uint curEpoch = getCurrentEpoch();
        address payable voter = msg.sender;

        uint lVoterStake = candidateData[candidate].latestVoterStakes[voter];
        require(lVoterStake >= amount, "amount too big to withdraw");

        uint remainAmount = lVoterStake.sub(amount);

        if (voter == candidateData[candidate].owner) {
            // owner, remainAmount must be >= minValidatorStake, otherwise need to use resign & withdraw
            require(remainAmount >= minValidatorStake, "remain amount of owner is too low");
        } else {
            // normal voter, either withdraw all or remain amount must be >= minVoterCap
            require(
                remainAmount == 0 || remainAmount >= minVoterCap,
                "remain amount must be either 0 or at least min voter cap"
            );
        }

        // update voter's latest stake and current epoch stake
        candidateData[candidate].latestVoterStakes[voter] = remainAmount;
        candidateData[candidate].voterStakes[voter][curEpoch] = remainAmount;

        // update candidate's latest stake and current epoch stake
        candidateData[candidate].latestTotalStakes = candidateData[candidate].latestTotalStakes.sub(amount);
        candidateData[candidate].totalStakes[curEpoch] = candidateData[candidate].latestTotalStakes;

        // transfer funds back to user
        voter.transfer(amount);

        emit Unvoted(voter, candidate, amount);
    }

    event Registered(address candidate, address owner);
    function register(address _candidate, address _owner) onlyAdmin onlyNotCandidate(_candidate) public {
        require(_candidate != address(0), "_candidate address is missing");
        require(_owner != address(0), "_owner address is missing");
        // not current candidate
        candidateData[_candidate] = CandidateData({
           owner: _owner,
           isCandidate: true,
           latestTotalStakes: 0
        });
        candidates.push(_candidate);
        emit Registered(_candidate, _owner);
    }

    event Resigned(address _candidate);
    // when a candidate resigns, at least minValidatorStake will be locked
    // after 2 epochs candidate can withdraw
    function resign(address _candidate) nonReentrant onlyValidCandidate(_candidate) public {
        // only owner can resign
        address payable sender = msg.sender;
        require(candidateData[_candidate].owner == sender);

        // remove from candidate list
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
        uint curEpoch = getCurrentEpoch();

        uint ownerStake = candidateData[_candidate].latestVoterStakes[sender];
        uint lockedStake = ownerStake.min(minValidatorStake);
        uint withdawableStake = ownerStake.sub(lockedStake);

        candidateData[_candidate].latestTotalStakes = candidateData[_candidate].latestTotalStakes.sub(ownerStake);
        candidateData[_candidate].totalStakes[curEpoch] = candidateData[_candidate].latestTotalStakes;
        candidateData[_candidate].latestVoterStakes[sender] = 0;
        candidateData[_candidate].voterStakes[sender][curEpoch] = 0;

        // locked this fund for 2 epochs
        uint unlockEpoch = curEpoch.add(STAKER_LOCKING_PERIOD);
        withdrawalCap[sender][unlockEpoch] = withdrawalCap[sender][unlockEpoch].add(lockedStake);

        if (withdawableStake > 0) {
            // transfer funds back to owner
            sender.transfer(withdawableStake);
        }
        emit Resigned(_candidate);
    }

    event Withdraw(address _staker, uint _amount);
    function withdraw(uint epoch) nonReentrant public returns(bool) {
        uint curEpoch = getCurrentEpoch();
        require(curEpoch >= epoch, "can not withdraw for future epoch");

        address payable sender = msg.sender;
        uint amount = withdrawalCap[sender][epoch];
        withdrawalCap[sender][epoch] = 0;

        require(amount > 0, "withdraw cap is 0");

        // transfer funds back to owner
        sender.transfer(amount);

        return true;
    }

    function updateMaxValidatorSize(uint newMaxValidatorSize) onlyAdmin public {
        maxValidatorSize = newMaxValidatorSize;
    }

    function getCurrentEpoch() public view returns(uint) {
        return (block.number.sub(startBlock)).div(epochPeriod);
    }

    function getAllCandidates() public view returns(address[] memory _candidates) {
        _candidates = candidates;
    }

    function getListCandidates()
        public view returns(address[] memory _candidates, uint[] memory _stakes, uint32 _maxValSize, uint32 _epochSize)
    {
        _maxValSize = uint32(maxValidatorSize);
        _epochSize = uint32(getCurrentEpoch());

        uint curEpoch = getCurrentEpoch();
        if (curEpoch == 0) {
             _candidates = new address[](initCandidates.length);
             _stakes = new uint[](initCandidates.length);
             for(uint i = 0; i < initCandidates.length; i++) {
                 _candidates[i] = initCandidates[i];
                 _stakes[i] = minValidatorStake;
             }
             return (_candidates, _stakes, _maxValSize, _epochSize);
        }
        // using previous epoch data to compute list validators
        uint epoch = curEpoch - 1;

        uint validCandiateCount;
        // only count candidate with his own stake >= minValidatorStake
        for(uint i = 0; i < candidates.length; i++) {
            address owner = candidateData[candidates[i]].owner;
            if (candidateData[candidates[i]].voterStakes[owner][epoch] >= minValidatorStake) {
                validCandiateCount++;
            }
        }

        _candidates = new address[](validCandiateCount);
        _stakes = new uint[](validCandiateCount);

        uint index = 0;
        for(uint i = 0; i < candidates.length; i++) {
            address owner = candidateData[candidates[i]].owner;
            if (candidateData[candidates[i]].voterStakes[owner][epoch] >= minValidatorStake) {
                _candidates[index] = candidates[i];
                _stakes[index] = candidateData[candidates[i]].totalStakes[epoch];
                index++;
            }
        }
    }

    function epochSize() public view returns(uint32) {
        return uint32(getCurrentEpoch());
    }

    function getCandidateData(address _candidate)
        public view
        returns(bool _isCandidate, address _owner, uint _latestTotalStakes)
    {
        _isCandidate = candidateData[_candidate].isCandidate;
        _owner = candidateData[_candidate].owner;
        _latestTotalStakes = candidateData[_candidate].latestTotalStakes;
    }

    function getVoterStake(address _candidate, address _voter, uint _epoch) public view returns(uint) {
        return candidateData[_candidate].voterStakes[_voter][_epoch];
    }

    function getVoterLatestStake(address _candidate, address _voter) public view returns(uint) {
        return candidateData[_candidate].latestVoterStakes[_voter];
    }

    function getTotalStakes(address _candidate, uint epoch) public view returns(uint) {
        return candidateData[_candidate].totalStakes[epoch];
    }
}