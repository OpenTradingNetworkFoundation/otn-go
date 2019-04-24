package main

import (
	"encoding/csv"
	"encoding/hex"
	"io"
	"log"
	"math"
	"os"
	"strconv"
	"time"

	"github.com/opentradingnetworkfoundation/otn-go/api"
	"github.com/opentradingnetworkfoundation/otn-go/wallet"
	"go.uber.org/ratelimit"

	"github.com/spf13/cobra"
	"github.com/opentradingnetworkfoundation/otn-go/objects"
)

var distributeCmd = cobra.Command{
	Use:   "distribute",
	Short: "Distribute funds",
	Run:   distribute,
}

var distributeOptions struct {
	input     string
	output    string
	asset     string
	amount    float64
	percent   float64
	from      string
	pvtKey    string
	commit    bool
	rateLimit int
}

func init() {
	rootCmd.AddCommand(&distributeCmd)
	distributeCmd.Flags().StringVarP(&distributeOptions.input, "input", "i", "", "Input file (output from listaccounts command)")
	distributeCmd.Flags().StringVarP(&distributeOptions.output, "output", "o", "log.csv", "File to output transaction list")
	distributeCmd.Flags().StringVar(&distributeOptions.asset, "asset", "OTN", "Specify asset")
	distributeCmd.Flags().StringVar(&distributeOptions.from, "from", "", "Account to distribute funds from")
	distributeCmd.Flags().StringVar(&distributeOptions.pvtKey, "key", "", "Private key for the funding account")
	distributeCmd.Flags().Float64Var(&distributeOptions.amount, "amount", 0, "Amount to distribute")
	distributeCmd.Flags().Float64Var(&distributeOptions.percent, "percent", 0, "Percent to pay out")
	distributeCmd.Flags().IntVar(&distributeOptions.rateLimit, "rate", 10, "Maximum transactions per second")
	distributeCmd.Flags().BoolVarP(&distributeOptions.commit, "commit", "c", false, "Perform actual transfers")
}

type accountBalance struct {
	ID      objects.GrapheneID
	Name    string
	Balance uint64
}

func loadAccountBalances() ([]*accountBalance, uint64) {
	file, err := os.Open(distributeOptions.input)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	reader := csv.NewReader(file)

	accounts := make([]*accountBalance, 0)
	var total uint64

	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}

		if len(line) < 3 {
			log.Fatal("Invalid input file format")
		}

		balance, err := strconv.ParseUint(line[2], 10, 64)
		if err != nil {
			log.Fatalf("Failed to parse balance '%s': %v", line[2], err)
		}

		acc := &accountBalance{
			ID:      *objects.NewGrapheneID(objects.ObjectID(line[0])),
			Name:    line[1],
			Balance: balance,
		}

		accounts = append(accounts, acc)
		total += acc.Balance
	}

	return accounts, total
}

func distribute(cmd *cobra.Command, args []string) {
	rpc := createNodeConnection()
	dbAPI, err := rpc.DatabaseAPI()
	if err != nil {
		log.Fatalf("Failed to get database API: %v", err)
	}
	asset := getAsset(dbAPI, distributeOptions.asset)
	wallet := wallet.NewWallet()

	log.Printf("Asset: %s (%s)", asset.Symbol, asset.ID)

	prec := math.Pow10(asset.Precision)
	accounts, total := loadAccountBalances()
	totalFl := float64(total) / prec
	log.Printf("Total users balance: %f (%d users)", totalFl, len(accounts))

	if distributeOptions.percent != 0 {
		distributeOptions.amount = totalFl * distributeOptions.percent / 100.0
	}

	if distributeOptions.amount == 0 {
		return
	}

	log.Printf("Amount to distribute: %f (%.2f percent of user balances)",
		distributeOptions.amount, distributeOptions.amount*100/totalFl)

	fromAccount, err := dbAPI.GetAccountByName(distributeOptions.from)
	if err != nil {
		log.Fatalf("Failed to get account '%s': %v", distributeOptions.from, err)
	}

	fromBalance, err := dbAPI.GetAccountBalances(fromAccount.ID, asset.ID)
	if err != nil {
		log.Fatalf("Failed to get account balance: %v", err)
	}

	if distributeOptions.pvtKey != "" {
		if err := wallet.AddPrivateKeys([]string{distributeOptions.pvtKey}); err != nil {
			log.Fatalf("Failed to import key: %v", err)
		}
	}

	funds := float64(fromBalance[0].Amount) / prec
	log.Printf("Account '%s' has %f %s", fromAccount.Name, funds, asset.Symbol)
	if funds < distributeOptions.amount {
		log.Fatalf("Account don't have sufficient funds to distribute %f %s",
			distributeOptions.amount, asset.Symbol)
	}

	output, err := os.Create(distributeOptions.output)
	if err != nil {
		log.Fatalf("Failed to create output file: %v", err)
	}
	defer output.Close()

	writer := csv.NewWriter(output)
	defer writer.Flush()

	var totalPayments float64

	throttle := ratelimit.New(distributeOptions.rateLimit)

	broadcastAPI, err := rpc.BroadcastAPI()
	if err != nil {
		log.Fatalf("Unable to get broadcast API: %v", err)
	}
	for _, acc := range accounts {
		if acc.Balance == 0 {
			continue
		}
		balance := float64(acc.Balance) / prec
		userAmount := balance / totalFl * distributeOptions.amount
		log.Printf("User: %s, balance %f, payment %f", acc.Name, balance, userAmount)
		totalPayments += userAmount

		transfer := &objects.TransferOperation{
			From: fromAccount.ID,
			To:   acc.ID,
			Amount: objects.AssetAmount{
				Asset:  asset.ID,
				Amount: objects.Int64(userAmount * prec),
			},
			Extensions: objects.Extensions{},
		}

		tx, err := api.ConstructTransaction(dbAPI, &objects.CoreAssetID, transfer)
		if err != nil {
			log.Fatalf("Failed to create transaction: %v", err)
		}

		if err := api.SignTransaction(dbAPI, wallet.GetKeys(), tx); err != nil {
			log.Fatalf("Failed to sign transaction: %v", err)
		}

		// apply rate limiting
		throttle.Take()

		if distributeOptions.commit {
			if err := broadcastAPI.BroadcastTransaction(tx); err != nil {
				log.Fatalf("Failed to broadcast transacton: %v", err)
			}
		}

		txID := hex.EncodeToString(tx.ID())
		log.Printf("TxID: %s", txID)

		writer.Write([]string{
			time.Now().Format(time.RFC3339),
			acc.ID.String(),
			acc.Name,
			strconv.FormatUint(acc.Balance, 10),
			strconv.FormatUint(uint64(transfer.Amount.Amount), 10),
			txID,
		})
	}

	log.Printf("Total payout: %f (requested: %f)", totalPayments, distributeOptions.amount)
}
