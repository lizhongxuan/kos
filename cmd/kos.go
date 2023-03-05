package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"kos/pkg/api"
)

var kosCmd = &cobra.Command{
	Use:   "kos",
	Short: "kkk object storage",
	Long:  "This command running object storage server",
	Run:   serverMain,
}

func serverMain(cmd *cobra.Command, args []string) {
	fmt.Println("run kos.")
	api.RegisterApi(":9090")
}
