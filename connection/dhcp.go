package connection

import (
  "log"
  "os"
  "os/signal"
  "github.com/keita0805carp/cacis/cacis"

  "github.com/coredhcp/coredhcp/server"
  "github.com/coredhcp/coredhcp/config"
  "github.com/coredhcp/coredhcp/plugins"
  pl_leasetime "github.com/coredhcp/coredhcp/plugins/leasetime"
  pl_netmask "github.com/coredhcp/coredhcp/plugins/netmask"
  pl_range "github.com/coredhcp/coredhcp/plugins/range"
  pl_router "github.com/coredhcp/coredhcp/plugins/router"
  pl_serverid "github.com/coredhcp/coredhcp/plugins/serverid"
)

const configPath = "./connection/dhcp.conf"

var desiredPlugins = []*plugins.Plugin{
  &pl_leasetime.Plugin,
  &pl_netmask.Plugin,
  &pl_range.Plugin,
  &pl_router.Plugin,
  &pl_serverid.Plugin,
}

func DHCP() {
  c := make(chan os.Signal, 1)
  signal.Notify(c, os.Interrupt)
  log.Println("[Debug] Starting dhcpd...")

	config, err := config.Load(configPath)
  cacis.Error(err)
	// register plugins
  for _, plugin := range desiredPlugins {
    err := plugins.RegisterPlugin(plugin)
    cacis.Error(err)
  }

  srv, err := server.Start(config)
  cacis.Error(err)
  log.Println("[Debug] Running dhcpd... (Press Ctrl-C to End)")

  <-c

  log.Println("[Debug] Terminating dhcpd...")
  srv.Close()
}
