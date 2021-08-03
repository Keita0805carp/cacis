package cmd

import (
  "fmt"
  "github.com/keita0805carp/cacis/master"

  "github.com/spf13/cobra"
)

var (
    masterCmd = &cobra.Command{
        Use: "master",
        Run: masterCommand,
    }
)

func masterCommand(cmd *cobra.Command, args []string) {
    if err := masterAction(); err != nil {
        Exit(err, 1)
    }
}

func masterAction() (err error) {
    fmt.Println("Debug: This is master subcommand\n")
    master.Main()
    return nil
}

func init() {
    RootCmd.AddCommand(masterCmd)
}
