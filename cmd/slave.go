package cmd

import (
  "log"

  "github.com/keita0805carp/cacis/slave"

  "github.com/spf13/cobra"
)

var (
  leave bool

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
  if leave {
    log.Printf("[Debug]: Manual leave\n")
    slave.Unclustering()
    return
  }

  slave.Main()

  return nil
}

func init() {
  RootCmd.AddCommand(slaveCmd)
  slaveCmd.Flags().BoolVarP(&leave, "leave", "l", false, "Manual leave from k8s cluster")
}
