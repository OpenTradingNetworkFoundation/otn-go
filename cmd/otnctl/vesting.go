package main

import (
	"encoding/csv"
	"io"
	"log"
	"math"
	"os"
	"strconv"
	"time"

	"github.com/juju/errors"
	"github.com/spf13/cobra"

	"github.com/opentradingnetworkfoundation/otn-go/api"
	"github.com/opentradingnetworkfoundation/otn-go/objects"
	"github.com/opentradingnetworkfoundation/otn-go/wallet"
)

var vbCmd = cobra.Command{
	Use:   "vb",
	Short: "Vesting balanances",
}

var vbCreateCmd = cobra.Command{
	Use:   "create",
	Short: "Create vesting balance",
	RunE:  createVestingBalance,
}

var vbCreateOptions struct {
	From   string
	Keys   []string
	To     string
	ToList string
	Period time.Duration
	Cliff  time.Duration
	Start  string
	Asset  string
}

func init() {
	cf := vbCreateCmd.Flags()

	cf.StringVar(&vbCreateOptions.From, "from", "", "creator of the vesting balance")
	cf.StringSliceVar(&vbCreateOptions.Keys, "key", []string{}, "private key for the 'from' account")
	cf.StringVar(&vbCreateOptions.To, "to", "", "owner of the vesting balance (who will be able to withdraw)")
	cf.DurationVar(&vbCreateOptions.Period, "period", 0, "vesting period")
	cf.DurationVar(&vbCreateOptions.Cliff, "cliff", 0, "vesting cliff")
	cf.StringVar(&vbCreateOptions.Start, "start", "", "start date")
	cf.StringVar(&vbCreateOptions.Asset, "asset", "OTN", "asset used for vesting balances")

	vbCmd.AddCommand(&vbCreateCmd)
	rootCmd.AddCommand(&vbCmd)
}

func getAccount(db api.DatabaseAPI, name string) *objects.Account {
	acc, err := db.GetAccountByName(name)
	if err != nil {
		log.Fatalf("Failed to get account %s: %v", name, err)
	}

	return acc
}

func parseAnyTime(str string) (time.Time, error) {
	fmts := []string{
		time.RFC3339,
		"2006-01-02T15:04:05",
		"2006-01-02T15:04",
		"2006-01-02 15:04",
		"2006-01-02",
	}

	for _, fmt := range fmts {
		t, err := time.Parse(fmt, str)
		if err == nil {
			return t, nil
		}
	}

	return time.Time{}, errors.Errorf("failed to parse time")
}

type VestingBalanceOperation struct {
	To     string
	Amount float64
}

func loadVestingList(path string) (result []*VestingBalanceOperation, err error) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()
	reader := csv.NewReader(file)
	for {
		line, err := reader.Read()
		if err != nil {
			break
		}

		amount, err := strconv.ParseFloat(line[1], 64)
		if err != nil {
			return nil, err
		}

		result = append(result, &VestingBalanceOperation{
			To:     line[0],
			Amount: amount,
		})
	}

	if err == io.EOF {
		err = nil
	}

	return
}

func createVestingBalance(cmd *cobra.Command, args []string) error {
	startTime, err := parseAnyTime(vbCreateOptions.Start)
	if err != nil {
		log.Fatal(err)
	}

	wallet := wallet.NewWallet()
	err = wallet.AddPrivateKeys(vbCreateOptions.Keys)
	if err != nil {
		log.Fatal(err)
	}

	if vbCreateOptions.To == "" {
		return errors.Errorf("please specify receiver (--to)")
	}

	rpc := createNodeConnection()
	db, err := rpc.DatabaseAPI()
	if err != nil {
		log.Fatal(err)
	}

	var toList []*VestingBalanceOperation

	if vbCreateOptions.To[0] == '@' {
		toList, err = loadVestingList(vbCreateOptions.To[1:])
		if err != nil {
			return err
		}
	} else {
		amount, err := strconv.ParseFloat(args[0], 64)
		if err != nil {
			return errors.Annotatef(err, "invalid amount: '%s'", args[0])
		}
		toList = append(toList, &VestingBalanceOperation{
			To:     vbCreateOptions.To,
			Amount: amount,
		})
	}

	asset := getAsset(db, vbCreateOptions.Asset)
	from := getAccount(db, vbCreateOptions.From)

	for _, vb := range toList {
		amount := vb.Amount
		to := getAccount(db, vb.To)

		log.Printf("Creating vesting balance for user %s (%s) with amount = %f %s\n",
			to.Name, to.ID.String(), amount, asset.Symbol)
		log.Printf("Starting date: %s, duration: %s", startTime.String(), vbCreateOptions.Period.String())

		op := &objects.VestingBalanceCreateOperation{
			Creator: from.ID,
			Owner:   to.ID,
			Amount: objects.AssetAmount{
				Asset:  asset.ID,
				Amount: objects.Int64(math.Pow10(asset.Precision) * amount),
			},
			Policy: objects.VestingPolicyInitializer{
				Type: objects.VestingPolicyTypeLinear,
				Value: objects.LinearVestingPolicyInitializer{
					BeginTimestamp:         objects.NewTime(startTime),
					VestingCliffSeconds:    objects.UInt32(vbCreateOptions.Cliff / time.Second),
					VestingDurationSeconds: objects.UInt32(vbCreateOptions.Period / time.Second),
				},
			},
		}

		_, err = api.SignAndBroadcast(rpc, wallet.GetKeys(), objects.CoreAssetID, op)
		if err != nil {
			return errors.Annotate(err, "failed to send transaction")
		}
	}

	log.Printf("All done")
	return nil
}
