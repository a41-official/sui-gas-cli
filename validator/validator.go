package validator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type SuiSystemState struct {
	Jsonrpc string `json:"jsonrpc"`
	Result  struct {
		Epoch                                 string      `json:"epoch"`
		ProtocolVersion                       string      `json:"protocolVersion"`
		SystemStateVersion                    string      `json:"systemStateVersion"`
		StorageFundTotalObjectStorageRebates  string      `json:"storageFundTotalObjectStorageRebates"`
		StorageFundNonRefundableBalance       string      `json:"storageFundNonRefundableBalance"`
		ReferenceGasPrice                     string      `json:"referenceGasPrice"`
		SafeMode                              bool        `json:"safeMode"`
		SafeModeStorageRewards                string      `json:"safeModeStorageRewards"`
		SafeModeComputationRewards            string      `json:"safeModeComputationRewards"`
		SafeModeStorageRebates                string      `json:"safeModeStorageRebates"`
		SafeModeNonRefundableStorageFee       string      `json:"safeModeNonRefundableStorageFee"`
		EpochStartTimestampMs                 string      `json:"epochStartTimestampMs"`
		EpochDurationMs                       string      `json:"epochDurationMs"`
		StakeSubsidyStartEpoch                string      `json:"stakeSubsidyStartEpoch"`
		MaxValidatorCount                     string      `json:"maxValidatorCount"`
		MinValidatorJoiningStake              string      `json:"minValidatorJoiningStake"`
		ValidatorLowStakeThreshold            string      `json:"validatorLowStakeThreshold"`
		ValidatorVeryLowStakeThreshold        string      `json:"validatorVeryLowStakeThreshold"`
		ValidatorLowStakeGracePeriod          string      `json:"validatorLowStakeGracePeriod"`
		StakeSubsidyBalance                   string      `json:"stakeSubsidyBalance"`
		StakeSubsidyDistributionCounter       string      `json:"stakeSubsidyDistributionCounter"`
		StakeSubsidyCurrentDistributionAmount string      `json:"stakeSubsidyCurrentDistributionAmount"`
		StakeSubsidyPeriodLength              string      `json:"stakeSubsidyPeriodLength"`
		StakeSubsidyDecreaseRate              int         `json:"stakeSubsidyDecreaseRate"`
		TotalStake                            string      `json:"totalStake"`
		ActiveValidators                      []Validator `json:"activeValidators"`
	} `json:"result"`
}

type Validator struct {
	SuiAddress                   string      `json:"suiAddress"`
	ProtocolPubkeyBytes          string      `json:"protocolPubkeyBytes"`
	NetworkPubkeyBytes           string      `json:"networkPubkeyBytes"`
	WorkerPubkeyBytes            string      `json:"workerPubkeyBytes"`
	ProofOfPossessionBytes       string      `json:"proofOfPossessionBytes"`
	Name                         string      `json:"name"`
	Description                  string      `json:"description"`
	ImageUrl                     string      `json:"imageUrl"`
	ProjectUrl                   string      `json:"projectUrl"`
	NetAddress                   string      `json:"netAddress"`
	P2PAddress                   string      `json:"p2pAddress"`
	PrimaryAddress               string      `json:"primaryAddress"`
	WorkerAddress                string      `json:"workerAddress"`
	NextEpochProtocolPubkeyBytes interface{} `json:"nextEpochProtocolPubkeyBytes"`
	NextEpochProofOfPossession   interface{} `json:"nextEpochProofOfPossession"`
	NextEpochNetworkPubkeyBytes  interface{} `json:"nextEpochNetworkPubkeyBytes"`
	NextEpochWorkerPubkeyBytes   interface{} `json:"nextEpochWorkerPubkeyBytes"`
	NextEpochNetAddress          interface{} `json:"nextEpochNetAddress"`
	NextEpochP2PAddress          interface{} `json:"nextEpochP2pAddress"`
	NextEpochPrimaryAddress      interface{} `json:"nextEpochPrimaryAddress"`
	NextEpochWorkerAddress       interface{} `json:"nextEpochWorkerAddress"`
	VotingPower                  string      `json:"votingPower"`
	OperationCapId               string      `json:"operationCapId"`
	GasPrice                     string      `json:"gasPrice"`
	CommissionRate               string      `json:"commissionRate"`
	NextEpochStake               string      `json:"nextEpochStake"`
	NextEpochGasPrice            string      `json:"nextEpochGasPrice"`
	NextEpochCommissionRate      string      `json:"nextEpochCommissionRate"`
	StakingPoolId                string      `json:"stakingPoolId"`
	StakingPoolActivationEpoch   string      `json:"stakingPoolActivationEpoch"`
	StakingPoolDeactivationEpoch interface{} `json:"stakingPoolDeactivationEpoch"`
	StakingPoolSuiBalance        string      `json:"stakingPoolSuiBalance"`
	RewardsPool                  string      `json:"rewardsPool"`
	PoolTokenBalance             string      `json:"poolTokenBalance"`
	PendingStake                 string      `json:"pendingStake"`
	PendingTotalSuiWithdraw      string      `json:"pendingTotalSuiWithdraw"`
	PendingPoolTokenWithdraw     string      `json:"pendingPoolTokenWithdraw"`
	ExchangeRatesId              string      `json:"exchangeRatesId"`
	ExchangeRatesSize            string      `json:"exchangeRatesSize"`
}

func GetAllValidatorsAndPrint() error {
	var suiSystemState SuiSystemState
	err := getValidators(&suiSystemState)
	if err != nil {
		fmt.Println("Error in getting validators:", err)
		return err
	}

	//print all validators
	for i, validator := range suiSystemState.Result.ActiveValidators {
		jsonData, _ := json.MarshalIndent(validator, "", "    ")
		if i < len(suiSystemState.Result.ActiveValidators)-1 {
			fmt.Println(string(jsonData) + ",")
		} else {
			fmt.Println(string(jsonData))
		}
	}
	return nil
}

func GetOneValidatorInfo(name string) error {
	var suiSystemState SuiSystemState
	err := getValidators(&suiSystemState)
	if err != nil {
		fmt.Println("Error in getting validators:", err)
		return err
	}

	//print all validators
	for _, validator := range suiSystemState.Result.ActiveValidators {
		if validator.Name == name {
			jsonData, _ := json.MarshalIndent(validator, "", "    ")
			fmt.Println(string(jsonData))
			break
		}
	}
	return nil
}

func GetAllValidators() ([]Validator, error) {
	var suiSystemState SuiSystemState
	err := getValidators(&suiSystemState)
	if err != nil {
		fmt.Println("Error in getting validators:", err)
		return nil, err
	}

	fmt.Println("Total validators:", len(suiSystemState.Result.ActiveValidators))
	return suiSystemState.Result.ActiveValidators, nil
}

func getValidators(suiSystemState *SuiSystemState) error {

	// Define request URL
	url := "https://rpc-mainnet.suiscan.xyz:443/"

	// Define request body
	jsonReq := `{"jsonrpc":"2.0","id":1,"method":"suix_getLatestSuiSystemState"}`

	// Create a new request using http
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(jsonReq)))

	// If there was an error in creating the request, handle it
	if err != nil {
		fmt.Println("Error in creating the request:", err)
		return err
	}

	// Add required headers to the request
	req.Header.Add("Content-Type", "application/json")

	// Send the request via a client
	client := &http.Client{}
	resp, err := client.Do(req)

	// If there was an error in sending the request, handle it
	if err != nil {
		fmt.Println("Error in sending the request:", err)
		return err
	}

	// Defer the closing of the response body
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error in reading the response body:", err)
		return err
	}

	// Unmarshal the JSON body into a map
	err = json.Unmarshal(body, &suiSystemState)
	if err != nil {
		fmt.Println("Error in unmarshalling the response body:", err)
		return err
	}

	return nil
}
