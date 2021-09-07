package cmd

import (
  "log"
  "os"
  "os/signal"
  "time"
  "strings"
  "github.com/keita0805carp/cacis/connection"
  //"github.com/keita0805carp/cacis/master"

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
    log.Printf("\n[Debug]: Run Main Master Process\n")

    terminate := make(chan os.Signal, 1)
    signal.Notify(terminate, os.Interrupt)

    adapterAddr, adapterId := connection.Initialize()
    //UUID := connection.genUUID()
    UUID := "12345678-9012-3456-7890-abcdefabcdef"
    ssid := strings.Replace(adapterAddr, ":", "", 5)
    pass := strings.Replace(UUID, "-", "", 4)

    cancel := make(chan struct{})
    go connection.Advertise(cancel, UUID, adapterAddr, adapterId, ssid, pass)

    go connection.StartHostapd(cancel, ssid, pass)

    go connection.DHCP(cancel)

    //master.Main()

    <-terminate
    close(cancel)
    log.Printf("\n[Debug]: Terminating Main Master Process...\n")
    time.Sleep(10 * time.Second)
    log.Printf("\n[Debug]: Terminate Main Master Process\n")

    return nil
}

func init() {
    RootCmd.AddCommand(masterCmd)
}
