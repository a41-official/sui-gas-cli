package main

import (
	"github.com/a41-official/sui-gas-cli/gas"
	"github.com/a41-official/sui-gas-cli/validator"
	"github.com/spf13/cobra"
)

func main() {

	// Validator commands
	var cmdValidator = &cobra.Command{
		Use:   "validator",
		Short: "Validator commands",
	}

	var cmdGetValidator = &cobra.Command{
		Use:   "get-all",
		Short: "Get validators metadata",
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			validator.GetAllValidatorsAndPrint()
		},
	}

	var cmdGetOneValidator = &cobra.Command{
		Use:   "get [validator name]",
		Short: "Get specific validator metadata",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			validator.GetOneValidatorInfo(args[0])
		},
	}

	// Gas commands
	var cmdGas = &cobra.Command{
		Use:   "gas",
		Short: "Gas commands",
	}

	var cmdGasCalc = &cobra.Command{
		Use:   "calc",
		Short: "Calculate gas price",
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			validators, err := validator.GetAllValidators()
			if err != nil {
				panic(err)
			}
			gas.CalcGasPrice(validators)
		},
	}

	var cmdGasPriceSubmit = &cobra.Command{
		Use:   "auto-submit",
		Short: "Submit gas price by auto calculation",
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			validators, err := validator.GetAllValidators()
			if err != nil {
				panic(err)
			}
			nextReferenceGasPrice := gas.CalcGasPrice(validators)
			gas.SubmitGasPrice(nextReferenceGasPrice - 5)
		},
	}

	var rootCmd = &cobra.Command{Use: "sui-tool"}

	rootCmd.AddCommand(cmdValidator)
	cmdValidator.AddCommand(cmdGetValidator)
	cmdValidator.AddCommand(cmdGetOneValidator)

	rootCmd.AddCommand(cmdGas)
	cmdGas.AddCommand(cmdGasCalc)
	cmdGas.AddCommand(cmdGasPriceSubmit)
	rootCmd.Execute()

}
