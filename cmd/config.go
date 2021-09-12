package cmd

import (
  "fmt"

  "github.com/keita0805carp/cacis/master"

  "github.com/spf13/cobra"
)

var (
    configCmd = &cobra.Command{
        Use: "config",
        Run: configCommand,
    }
)
var Source string

func configCommand(cmd *cobra.Command, args []string) {
    if err := configAction(); err != nil {
        Exit(err, 1)
    }
}

func configAction() (err error) {
  fmt.Println("Export Config")
  master.ExportKubeconfig("/home/ubuntu/.kube/config")
  return nil
}

func init() {
    RootCmd.AddCommand(configCmd)
    configCmd.Flags().StringVarP(&Source, "source", "s", "", "Source directory to read from")
}
