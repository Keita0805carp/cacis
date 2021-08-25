package cmd

import (
  "log"
  "github.com/keita0805carp/cacis/connection"
  //"github.com/keita0805carp/cacis/slave"

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
    log.Printf("[Debug] This is slave subcommand\n")
    connection.Discover()
    //connection.Connect()
    //slave.Main()
    return nil
}

func init() {
    RootCmd.AddCommand(slaveCmd)
}
