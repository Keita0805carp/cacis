package cmd

import (
  "log"
  "os"
  "syscall"
  "os/signal"
  "time"

  "github.com/keita0805carp/cacis/master"
  "github.com/keita0805carp/cacis/connection"

  "github.com/spf13/cobra"
)

var (
  main  bool
  setup bool

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
  log.Printf("\n[Debug]: Run Master Process\n")

  if main {
    log.Printf("\n[Debug]: Main Mode\n")
    terminate := make(chan os.Signal, 1)
    signal.Notify(terminate, syscall.SIGTERM, syscall.SIGINT, os.Interrupt)
    cancel := make(chan struct{})

    connection.Main(cancel)
    go master.Main(cancel)

    <-terminate
    close(cancel)
    log.Printf("\n[Debug]: Terminating Main Master Process...\n")
    time.Sleep(10 * time.Second)
    log.Printf("\n[Debug]: Terminate Main Master Process\n")
    return nil
  } else if setup {
    log.Printf("\n[Debug]: Setup Mode\n")
    master.Setup()
    return nil
  } else {
    log.Println("Please select option '--main' or '--setup'")
    return nil
  }
}

func init() {
  RootCmd.AddCommand(masterCmd)
  masterCmd.Flags().BoolVarP(&main, "main", "m", false, "Main Mode")
  masterCmd.Flags().BoolVarP(&setup, "setup", "s", false, "Setup Mode")
}
