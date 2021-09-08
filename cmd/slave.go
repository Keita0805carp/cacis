package cmd

import (
  "log"

  "github.com/keita0805carp/cacis/slave"
  "github.com/keita0805carp/cacis/connection"

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
    log.Printf("\n[Debug]: Run Main Slave Process\n")

    ssid, pw := connection.GetWifiInfo()
    connection.Connect(ssid, pw)
    slave.Main()

    return nil
}

func init() {
    RootCmd.AddCommand(slaveCmd)
}
