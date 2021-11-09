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
    Short: "Run Slave Process",
    Run: slaveCommand,
  }
)

func slaveCommand(cmd *cobra.Command, args []string) {
  if err := slaveAction(); err != nil {
    Exit(err, 1)
  }
}

func slaveAction() (err error) {
  for {
    log.Printf("\n[Debug]: Run Main Slave Process\n")

    ssid, pw := connection.GetWifiInfo()
    connection.Connect(ssid, pw)

    cancel := make(chan struct{})
    go connection.UnstableWifiEvent(cancel)
    slave.Main()

    <- cancel

    slave.Unclustering()
    connection.Disconnect()
    slave.WaitReadyMicrok8s()
  }

  return nil
}

func init() {
  RootCmd.AddCommand(slaveCmd)
}
