package gas

import (
	"fmt"
	. "github.com/a41-official/sui-gas-cli/validator"
	"log"
	"math/big"
	"os/exec"
	"sort"
	"strconv"
)

func calcMin(validators []Validator) int {
	min := int(^uint(0) >> 1) // set to max int value
	for _, v := range validators {
		val, _ := strconv.Atoi(v.NextEpochGasPrice)
		if val < min {
			min = val
		}
	}
	return min
}

func calcMax(validators []Validator) int {
	max := -int(^uint(0) >> 1) // set to min int value
	for _, v := range validators {
		val, _ := strconv.Atoi(v.NextEpochGasPrice)
		if val > max {
			max = val
		}
	}
	return max
}

func calcMean(validators []Validator) float64 {
	total := 0
	for _, v := range validators {
		val, _ := strconv.Atoi(v.NextEpochGasPrice)
		total += val
	}
	return float64(total) / float64(len(validators))
}

func calcMedian(validators []Validator) float64 {
	vals := make([]int, len(validators))
	for i, v := range validators {
		vals[i], _ = strconv.Atoi(v.NextEpochGasPrice)
	}
	sort.Ints(vals)

	if len(vals)%2 == 1 {
		return float64(vals[len(vals)/2])
	} else {
		return float64(vals[len(vals)/2-1]+vals[len(vals)/2]) / 2
	}
}

func calcWeightedMean(validators []Validator) *big.Float {
	gasMultiples := new(big.Int)
	totalDelegation := new(big.Int)

	for _, v := range validators {
		gasPrice := new(big.Int)
		gasPrice.SetString(v.NextEpochGasPrice, 10)

		stake := new(big.Int)
		stake.SetString(v.NextEpochStake, 10)

		gasMultiples.Add(gasMultiples, new(big.Int).Mul(gasPrice, stake))
		totalDelegation.Add(totalDelegation, stake)
	}

	// Convert to big.Float to perform the division
	result := new(big.Float).Quo(new(big.Float).SetInt(gasMultiples), new(big.Float).SetInt(totalDelegation))
	return result
}

func nextReferenceGasPrice(validators []Validator) int {
	var quorum, cumulativePower, referenceGasPrice int
	quorum = 6667

	sort.Slice(validators, func(i, j int) bool {
		a, _ := strconv.Atoi(validators[i].NextEpochGasPrice)
		b, _ := strconv.Atoi(validators[j].NextEpochGasPrice)
		return a < b
	})

	for _, v := range validators {
		if cumulativePower < quorum {
			referenceGasPrice, _ = strconv.Atoi(v.NextEpochGasPrice)
			votingPower, _ := strconv.Atoi(v.VotingPower)
			cumulativePower += votingPower
		}
	}
	return referenceGasPrice
}

func CalcGasPrice(validators []Validator) int {
	min := calcMin(validators)
	max := calcMax(validators)
	mean := calcMean(validators)
	median := calcMedian(validators)
	weightedMean := calcWeightedMean(validators)
	nextReferenceGasPrice := nextReferenceGasPrice(validators)

	// Print
	fmt.Println("===== Gas Price Calculation =====")
	fmt.Println("Total Validators: ", len(validators))
	fmt.Println("Min Reference Gas Price: ", min)
	fmt.Println("Max Reference Gas Price: ", max)
	fmt.Println("Mean Reference Gas Price: ", mean)
	fmt.Println("Stake Weighted Mean Reference Gas Price: ", weightedMean)
	fmt.Println("Median Reference Gas Price: ", median)
	fmt.Println("Estimated Next Reference Gas Price: ", nextReferenceGasPrice)

	return nextReferenceGasPrice
}

func SubmitGasPrice(gasPrice int) {

	cmd := exec.Command("sui", "validator", "update-gas-price", strconv.Itoa(gasPrice))

	output, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(output))
}
