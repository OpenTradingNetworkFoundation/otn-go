package main

import "log"

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	rootCmd.Execute()
}
