package main

import "github.com/spf13/cobra"

var rootCmd = cobra.Command{
	Use:   "otnctl",
	Short: "OTN command line tool",
}

var cfg struct {
	nodeAddress string
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfg.nodeAddress, "node", "wss://wallet.otn.org/ws", "Node rpc address")
}
