package main

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use: "cmd --pool={address} --from={address} --to={address} --amount={amount-in}",
	Run: nil,
}

func init() {
	rootCmd.PersistentFlags().String("pool", "0x0d4a11d5eeaac28ec3f61d100daf4d40471f1852", "pool address")
	rootCmd.PersistentFlags().String("from", "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2", "token from address")
	rootCmd.PersistentFlags().String("to", "0xdac17f958d2ee523a2206206994597c13d831ec7", "token to address")
	rootCmd.PersistentFlags().String("amount", "1000000000000000000", "amount")
}
