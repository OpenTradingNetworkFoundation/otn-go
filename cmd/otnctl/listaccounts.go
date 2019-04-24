package main

import (
	"encoding/csv"
	"log"
	"math"
	"os"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/opentradingnetworkfoundation/otn-go/api"
	"github.com/opentradingnetworkfoundation/otn-go/objects"
)

var listaccountsCmd = cobra.Command{
	Use:   "listaccounts",
	Short: "List accounts",
	Run:   listAccounts,
}

var listOptions struct {
	asset           string
	output          string
	minBalanceFloat float64

	// processed values
	minBalanceUint uint64
	assetID        *objects.GrapheneID
}

func init() {
	rootCmd.AddCommand(&listaccountsCmd)
	listaccountsCmd.Flags().StringVarP(&listOptions.output, "file", "f", "", "Write accounts list to the file")
	listaccountsCmd.Flags().Float64Var(&listOptions.minBalanceFloat, "minbalance", 0, "Minimal balance to include in list")
	listaccountsCmd.Flags().StringVar(&listOptions.asset, "asset", "OTN", "Report balance for this asset")
}

func assetBalance(assetID *objects.GrapheneID, balances []objects.AccountBalance) uint64 {
	for i := range balances {
		if balances[i].AssetID == *assetID {
			return uint64(balances[i].Balance)
		}
	}
	return 0
}

func csvAccountWriter(output *csv.Writer, acc *objects.FullAccount) error {
	balance := assetBalance(listOptions.assetID, acc.Balances)

	if balance < listOptions.minBalanceUint {
		return nil
	}

	return output.Write([]string{
		acc.Account.ID.String(),
		acc.Account.Name,
		strconv.FormatUint(balance, 10),
	})
}

func listAccounts(cmd *cobra.Command, args []string) {
	rpc := createNodeConnection()
	dbAPI, err := rpc.DatabaseAPI()
	if err != nil {
		log.Fatalf("Unable to get database API: %v", err)
	}
	asset := getAsset(dbAPI, listOptions.asset)

	listOptions.assetID = &asset.ID
	listOptions.minBalanceUint = uint64(math.Pow10(asset.Precision) * listOptions.minBalanceFloat)

	log.Printf("Asset: %s (%s), minimal balance: %f (%d)",
		listOptions.asset, listOptions.assetID,
		listOptions.minBalanceFloat, listOptions.minBalanceUint)

	output := os.Stdout

	if listOptions.output != "" {
		file, err := os.Create(listOptions.output)
		if err != nil {
			log.Fatal(err)
		}

		defer file.Close()
		output = file
	}

	csvOut := csv.NewWriter(output)
	api.LoadAccounts(dbAPI, func(acc *objects.FullAccount) error {
		return csvAccountWriter(csvOut, acc)
	})

	csvOut.Flush()
}
