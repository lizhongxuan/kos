package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"kos/command"
)

var (
	s3Options S3Options
)

func Main() error {
	rootCmd := &cobra.Command{
		Use:  "kos",
		Args: cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("run kos")
		},
	}
	rootCmd.AddCommand(serverCmd)
	if err := rootCmd.Execute(); err != nil {
		return err
	}
	return nil
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "kkk object storage",
	Long:  "This command running object storage server",
	Run:   serverMain,
}

func serverMain(cmd *cobra.Command, args []string) {
	fmt.Println("run kos.")
	s3Options := command.NewS3()
	s3Options.startS3Server()
	//api.RegisterApi(":9090")
}
