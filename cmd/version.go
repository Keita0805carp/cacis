package cmd

import (
  "fmt"
  "github.com/spf13/cobra"
)

var (
  versionCmd = &cobra.Command{
    Use: "version",
    Short: "Show version",
    Run: versionCommand,
  }
)

func versionCommand(cmd *cobra.Command, args []string) {
  if err := versionAction(); err != nil {
    Exit(err, 1)
  }
}

func versionAction() (err error) {
  fmt.Println("version: 0.0.1")
  return nil
}

func init() {
  RootCmd.AddCommand(versionCmd)
}
