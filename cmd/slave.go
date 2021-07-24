package cmd

import (
  "fmt"
  "github.com/keita0805carp/cacis/slave"

  "github.com/spf13/cobra"
)

var (
    slaveCmd = &cobra.Command{
        Use: "slave",
        Run: slaveCommand,
    }
)

func slaveCommand(cmd *cobra.Command, args []string) {
    if err := slaveAction(); err != nil {
        Exit(err, 1)
    }
}

func slaveAction() (err error) {
    fmt.Println("This is slave subcommand")
    slave.Main()
    return nil
}

func init() {
    RootCmd.AddCommand(slaveCmd)
}
