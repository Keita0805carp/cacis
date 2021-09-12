package cmd

import (
  "fmt"

  "github.com/keita0805carp/cacis/master"

  "github.com/spf13/cobra"
)

var (
  get    bool
  path   string

  configCmd = &cobra.Command{
    Use: "config",
    Short: "Manage kubeconfig",
    Run: configCommand,
  }
)

func configCommand(cmd *cobra.Command, args []string) {
  if err := configAction(); err != nil {
    Exit(err, 1)
  }
}

func configAction() (err error) {
  if path != "" {
    fmt.Println("Export Config")
    master.ExportKubeconfig(path)
    return nil
  }
  if get {
    str, _ := master.GetKubeconfig()
    fmt.Println(str)
    return nil
  }
  return nil
}

func init() {
  RootCmd.AddCommand(configCmd)
  configCmd.Flags().BoolVarP(&get, "get", "g", true, "Get kubeconfig")
  configCmd.Flags().StringVarP(&path, "path", "f", "", "Export kubeconfig")
}
