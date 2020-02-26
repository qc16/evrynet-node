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
    uint constant internal POWER110 = 2 ** 110;
    uint constant internal POWER36 = 2 ** 36;

    // fit in only 1 uint256: 36 bits (modifiedEpoch), 110 bits (curStake), 110 bits (preStake)
    struct StakeData {
        uint preStake;
        uint curStake;
        uint modifiedEpoch;
    }

    struct CandidateData {
        bool isCandidate;
        uint totalStakeData;
        address owner;
        // voter's stakes for each epoch
        mapping(address => uint) voterStakeData;
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
        uint _startBlock,
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
        // modifed epoch = 0, preStake = 0, curStake  = _minValidatorStake
        uint stakeData = encodeStakeData(0, _minValidatorStake, 0);
        for(uint i = 0; i < _candidates.length; i++) {
            candidateData[_candidates[i]] = CandidateData({
                isCandidate: true,
                owner: candidatesOwner,
                totalStakeData: stakeData
            });
            candidateData[_candidates[i]].voterStakeData[candidatesOwner] = stakeData;
        }

        initCandidates = _candidates;

        admin = _admin;
        startBlock = _startBlock;
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
        uint curEpoch = getCurrentEpoch();

        StakeData memory _totalStakeData = decodeStakeData(candidateData[candidate].totalStakeData);
        if (_totalStakeData.modifiedEpoch == curEpoch) {
            _totalStakeData.curStake = _totalStakeData.curStake.add(amount);
        } else {
            _totalStakeData.preStake = _totalStakeData.curStake;
            _totalStakeData.curStake = _totalStakeData.curStake.add(amount);
            _totalStakeData.modifiedEpoch = curEpoch;
        }
        // re-assign total stake data
        candidateData[candidate].totalStakeData = encodeStakeData(_totalStakeData);

        StakeData memory _voterStakeData = decodeStakeData(candidateData[candidate].voterStakeData[voter]);
        if (_voterStakeData.modifiedEpoch == curEpoch) {
            _voterStakeData.curStake = _voterStakeData.curStake.add(amount);
        } else {
            _voterStakeData.preStake = _voterStakeData.curStake;
            _voterStakeData.curStake = _voterStakeData.curStake.add(amount);
            _voterStakeData.modifiedEpoch = curEpoch;
        }
        candidateData[candidate].voterStakeData[voter] = encodeStakeData(_voterStakeData);

        emit Voted(voter, candidate, amount);
    }

    event Unvoted(address voter, address candidate, uint amount);
    function unvote(address candidate, uint amount) nonReentrant public {
        require(amount > 0, "withdraw amount must be positive");

        uint curEpoch = getCurrentEpoch();
        address payable voter = msg.sender;

        StakeData memory _voterStakeData = decodeStakeData(candidateData[candidate].voterStakeData[voter]);
        require(_voterStakeData.curStake >= amount, "amount too big to withdraw");

        uint remainAmount = _voterStakeData.curStake.sub(amount);

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

        if (_voterStakeData.modifiedEpoch == curEpoch) {
            _voterStakeData.curStake = remainAmount;
        } else {
            _voterStakeData.preStake = _voterStakeData.curStake;
            _voterStakeData.curStake = remainAmount;
            _voterStakeData.modifiedEpoch = curEpoch;
        }
        candidateData[candidate].voterStakeData[voter] = encodeStakeData(_voterStakeData);

        // update candidate's latest stake data
        StakeData memory _totalStakeData = decodeStakeData(candidateData[candidate].totalStakeData);
        if (_totalStakeData.modifiedEpoch == curEpoch) {
            _totalStakeData.curStake = _totalStakeData.curStake.sub(amount);
        } else {
            _totalStakeData.preStake = _totalStakeData.curStake;
            _totalStakeData.curStake = _totalStakeData.curStake.sub(amount);
            _totalStakeData.modifiedEpoch = curEpoch;
        }
        candidateData[candidate].totalStakeData = encodeStakeData(_totalStakeData);
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
           totalStakeData: encodeStakeData(0, 0, 0)
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

        uint curEpoch = getCurrentEpoch();
        StakeData memory _ownerStakeData = decodeStakeData(candidateData[_candidate].voterStakeData[sender]);
        uint ownerStake = _ownerStakeData.curStake;
        uint lockedStake = ownerStake.min(minValidatorStake);
        uint withdawableStake = ownerStake.sub(lockedStake);

        StakeData memory _totalStakeData = decodeStakeData(candidateData[_candidate].totalStakeData);
        if (_totalStakeData.modifiedEpoch == curEpoch) {
            _totalStakeData.curStake = _totalStakeData.curStake.sub(ownerStake);
        } else {
            _totalStakeData.preStake = _totalStakeData.curStake;
            _totalStakeData.curStake = _totalStakeData.curStake.sub(ownerStake);
            _totalStakeData.modifiedEpoch = curEpoch;
        }
        // update total stake data
        candidateData[_candidate].totalStakeData = encodeStakeData(_totalStakeData);
        candidateData[_candidate].voterStakeData[sender] = encodeStakeData(0, 0, 0); // just reset data

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

        StakeData memory _stakeData;
        address owner;
        uint eligibleStake;

        uint validCandiateCount;
        // only count candidate with his own stake >= minValidatorStake
        for(uint i = 0; i < candidates.length; i++) {
            owner = candidateData[candidates[i]].owner;
            _stakeData = decodeStakeData(candidateData[candidates[i]].voterStakeData[owner]);
            // modifiedEpoch <= curEpoch, either modifiedEpoch <= epoch or modifiedEpoch = epoch + 1
            if (_stakeData.modifiedEpoch <= epoch) {
                eligibleStake = _stakeData.curStake;
            } else {
                eligibleStake = _stakeData.preStake;
            }
            if (eligibleStake >= minValidatorStake) {
                validCandiateCount++;
            }
        }

        _candidates = new address[](validCandiateCount);
        _stakes = new uint[](validCandiateCount);

        uint index = 0;
        for(uint i = 0; i < candidates.length; i++) {
            owner = candidateData[candidates[i]].owner;
            _stakeData = decodeStakeData(candidateData[candidates[i]].voterStakeData[owner]);
            // modifiedEpoch <= curEpoch, either modifiedEpoch <= epoch or modifiedEpoch = epoch + 1
            if (_stakeData.modifiedEpoch <= epoch) {
                eligibleStake = _stakeData.curStake;
            } else {
                eligibleStake = _stakeData.preStake;
            }
            if (eligibleStake >= minValidatorStake) {
                _candidates[index] = candidates[i];
                _stakeData = decodeStakeData(candidateData[candidates[i]].totalStakeData);
                if (_stakeData.modifiedEpoch <= epoch) {
                    _stakes[index] = _stakeData.curStake;
                } else {
                    _stakes[index] = _stakeData.preStake;
                }
                index++;
            }
        }
    }

    // Return list of candidates, stakes, max valset and epoch number
    // using current stake data
    function getListCandidatesWithCurrentData()
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

        StakeData memory _stakeData;
        address owner;

        uint validCandiateCount;
        // only count candidate with his own stake >= minValidatorStake
        for(uint i = 0; i < candidates.length; i++) {
            owner = candidateData[candidates[i]].owner;
            _stakeData = decodeStakeData(candidateData[candidates[i]].voterStakeData[owner]);
            if (_stakeData.curStake >= minValidatorStake) {
                validCandiateCount++;
            }
        }

        _candidates = new address[](validCandiateCount);
        _stakes = new uint[](validCandiateCount);

        uint index = 0;
        for(uint i = 0; i < candidates.length; i++) {
            owner = candidateData[candidates[i]].owner;
            _stakeData = decodeStakeData(candidateData[candidates[i]].voterStakeData[owner]);
            if (_stakeData.curStake >= minValidatorStake) {
                _candidates[index] = candidates[i];
                _stakeData = decodeStakeData(candidateData[candidates[i]].totalStakeData);
                _stakes[index] = _stakeData.curStake;
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
        StakeData memory _stakeData = decodeStakeData(candidateData[_candidate].totalStakeData);
        _latestTotalStakes = _stakeData.curStake;
    }

    function getVoterStakeData(address _candidate, address _voter)
        public view
        returns(uint _preStake, uint _curStake, uint _lastModifiedEpoch)
    {
        StakeData memory _stakeData = decodeStakeData(candidateData[_candidate].voterStakeData[_voter]);
        _preStake = _stakeData.preStake;
        _curStake = _stakeData.curStake;
        _lastModifiedEpoch = _stakeData.modifiedEpoch;
    }

    function decodeStakeData(uint data) internal pure returns(StakeData memory stakeData) {
        stakeData.preStake = data & (POWER110.sub(1));
        stakeData.curStake = (data.div(POWER110)) & (POWER110.sub(1));
        stakeData.modifiedEpoch = (data.div(POWER110.mul(POWER110))) & (POWER36.sub(1));
    }

    function encodeStakeData(StakeData memory stakeData) internal pure returns(uint data) {
        data = stakeData.preStake & (POWER110.sub(1));
        data |= (stakeData.curStake & (POWER110.sub(1))).mul(POWER110);
        data |= (stakeData.modifiedEpoch & (POWER36.sub(1))).mul(POWER110).mul(POWER110);
    }

    function encodeStakeData(uint epoch, uint curStake, uint preStake) public pure returns(uint data) {
        data = preStake & (POWER110.sub(1));
        data |= (curStake & (POWER110.sub(1))).mul(POWER110);
        data |= (epoch & (POWER36.sub(1))).mul(POWER110).mul(POWER110);
    }
}