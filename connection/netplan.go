package connection

import (
  "os"
  "log"
  "time"
  "strings"
  "io/ioutil"

  "github.com/keita0805carp/cacis/cacis"
)

const (
  netplanConfTemplatePath = cacis.NetplanConfTemplatePath
  netplanConfPath = cacis.NetplanConfPath
)

func Connect(ssid, pw string) {
  genNetplanConfig(ssid, pw)
  log.Printf("[Info]  SSID: %s\n", ssid)
  log.Printf("[Info]  PASS: %s\n", pw)

  log.Printf("[Debug] Apply netplan config\n")
  cacis.ExecCmd("netplan apply", false)
  log.Printf("[Debug] Applied netplan config\n")
  time.Sleep(10 * time.Second)
}

func Disconnect() {
  log.Printf("[Debug] Delete netplan config\n")
  delNetplanConfig()
  log.Printf("[Debug] Apply netplan config\n")
}

func genNetplanConfig(ssid, pw string) {

  log.Printf("[Debug] Generate netplan Config...\n")
  bytes, err := ioutil.ReadFile(netplanConfTemplatePath)
  cacis.Error(err)
  config := string(bytes)

  config = strings.Replace(config, "{{SSID}}", ssid, 1)
  config = strings.Replace(config, "{{PASSWORD}}", pw, 1)

  ioutil.WriteFile(netplanConfPath, []byte(config), 0644)
  log.Printf("[Debug] Generated netplan Config\n")
}

func delNetplanConfig() {
  err := os.Remove(netplanConfPath)
  cacis.Error(err)
  cacis.ExecCmd("netplan apply", false)
  log.Printf("[Debug] Clean netplan\n")
}

