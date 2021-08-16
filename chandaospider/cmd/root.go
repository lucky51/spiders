package cmd

import "github.com/spf13/cobra"

const timeLayout = "2006-01-02 15:04:05"

var rootCmd = &cobra.Command{}

func init() {
	rootCmd.AddCommand(crawlCmd)
	rootCmd.AddCommand(jobCmd)
	rootCmd.AddCommand(serverCmd)
}

func Execute() error {
	return rootCmd.Execute()
}
