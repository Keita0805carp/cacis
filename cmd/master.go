package cmd

import (
  "log"

  "github.com/keita0805carp/cacis/master"

  "github.com/spf13/cobra"
)

var (
  main  bool
  setup bool

  masterCmd = &cobra.Command{
    Use: "master",
    Short: "Run Master Process",
    Run: masterCommand,
  }
)

func masterCommand(cmd *cobra.Command, args []string) {
  if err := masterAction(); err != nil {
    Exit(err, 1)
  }
}

func masterAction() (err error) {
  log.Printf("[Debug] Run Master Process\n")

  if main {
    log.Printf("[Debug] Main Mode\n")
    master.Main()
  } else if setup {
    log.Printf("[Debug] Setup Mode\n")
    master.Setup()
  } else {
    log.Printf("[Error] Please select option '--main' or '--setup'\n")
  }
  return nil
}

func init() {
  RootCmd.AddCommand(masterCmd)
  masterCmd.Flags().BoolVarP(&main, "main", "m", false, "Main Mode")
  masterCmd.Flags().BoolVarP(&setup, "setup", "s", false, "Setup Mode")
}
