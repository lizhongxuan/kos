package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var kosCmd = &cobra.Command{
	Use:   "kos",
	Short: "kkk object storage",
	Long:  "This command running object storage server",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("run kos.")
	},
}

func main() {
	rootCmd := &cobra.Command{
		Use:  "root",
		Args: cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("run root.")
		},
	}
	rootCmd.AddCommand(kosCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
