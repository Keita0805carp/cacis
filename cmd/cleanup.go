package cmd

import (
  "fmt"

  "github.com/keita0805carp/cacis/cacis"
  "github.com/keita0805carp/cacis/slave"

  "github.com/spf13/cobra"
)

var (
  cleanupCmd = &cobra.Command{
    Use: "cleanup",
    Short: "Uninstall and Cleanup",
    Run: cleanupCommand,
  }
)

func cleanupCommand(cmd *cobra.Command, args []string) {
  if err := cleanupAction(); err != nil {
    Exit(err, 1)
  }
}

func cleanupAction() (err error) {
  fmt.Println("cleanup command")
  slave.RemoveMicrok8s()
  cacis.RemoveTempDir()
  return nil
}

func init() {
  RootCmd.AddCommand(cleanupCmd)
}
