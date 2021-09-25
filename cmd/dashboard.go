package cmd

import (
	"fmt"
  "os/exec"

	"github.com/keita0805carp/cacis/cacis"

	"github.com/spf13/cobra"
)

var (
  token       bool
  patch       bool
  portfoward  bool

  dashboardCmd = &cobra.Command{
    Use: "dashboard",
    Short: "Controll Kubrentes dashboard",
    Run: dashboardCommand,
  }
)

func dashboardCommand(cmd *cobra.Command, args []string) {
  if err := dashboardAction(); err != nil {
    Exit(err, 1)
  }
}

func dashboardAction() (err error) {
  fmt.Println("Dashboard command")
  if token {
    fmt.Println("Token")
    output, err := exec.Command("sh", "-c", "kubectl -n kube-system describe $(kubectl -n kube-system get secret -o name | grep kubernetes-dashboard-token)").CombinedOutput()
    cacis.Error(err)
    fmt.Println(string(output))
  } else if patch {
    fmt.Println("Patch")
    output, err := exec.Command("sh", "-c", "kubectl -n kube-system patch service kubernetes-dashboard -p '{\"spec\":{\"type\":\"NodePort\"}}'").CombinedOutput()
    cacis.Error(err)
    fmt.Println(string(output))
  } else if portfoward {
    fmt.Println("Portforwarding...")
    output, err := exec.Command("sh", "-c", "kubectl -n kube-system port-forward --address 0.0.0.0 service/kubernetes-dashboard 8443:443").CombinedOutput()
    cacis.Error(err)
    fmt.Println(string(output))
  } else {
    fmt.Println("Error: Need some optinos")
    Exit(err, 1)
  }
  return nil
}

func init() {
  RootCmd.AddCommand(dashboardCmd)
  dashboardCmd.Flags().BoolVarP(&patch, "patch", "p", false, "Patch svc/kubernetes-dashboard: ClusterIP => NodePort")
  dashboardCmd.Flags().BoolVarP(&token, "token", "t", false, "Show dashboard token")
  dashboardCmd.Flags().BoolVarP(&portfoward, "portfoward", "P", false, "Portfoward dashboard to 0.0.0.0:8443")
}
