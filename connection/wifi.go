package connection

import (
  "log"
  "strings"
  "io/ioutil"

  "github.com/keita0805carp/cacis/cacis"
)

const (
  netplanConfTemplatePath = "connection/netplan.conf.template"
  netplanConfPath = "/etc/netplan/60-cacis.yaml"
)

func Connect(ssid, pw string) {
  genNetplanConfig(ssid, pw)
  log.Printf("[Info]  SSID: %s\n", ssid)
  log.Printf("[Info]  PASS: %s\n", pw)

  log.Println("[Debug] Apply netplan config")
  cacis.ExecCmd("netplan apply", false)
  log.Println("[Debug] Applied netplan config")
}

func genNetplanConfig(ssid, pw string) {

  log.Println("[Debug] Generate netplan Config...")
  bytes, err := ioutil.ReadFile(netplanConfTemplatePath)
  cacis.Error(err)
  config := string(bytes)

  config = strings.Replace(config, "{{SSID}}", ssid, 1)
  config = strings.Replace(config, "{{PASSWORD}}", pw, 1)

  ioutil.WriteFile(netplanConfPath, []byte(config), 0644)
  log.Println("[Debug] Generated netplan Config")
}
