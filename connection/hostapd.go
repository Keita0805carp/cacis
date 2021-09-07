package connection

import (
  "log"
  "strings"
  "io/ioutil"
  "os"
  "os/signal"

  "github.com/keita0805carp/cacis/cacis"
)

const (
  hostapdConfTemplatePath = "connection/hostapd.conf.template"
  hostapdConfPath = "connection/hostapd.conf"
)

func StartHostapd(ssid, pw string) {
  genConfig(ssid, pw)
  log.Printf("[Info]  SSID: %s\n", ssid)
  log.Printf("[Info]  PASS: %s\n", pw)

  cacis.ExecCmd("killall -q hostapd", false)

  log.Println("[Debug] Start hostapd in the Background")
  cacis.ExecCmd("hostapd -B " + hostapdConfPath, false)

  c := make(chan os.Signal, 1)
  signal.Notify(c, os.Interrupt)
  log.Println("[Debug] Running hostpad... (Press Ctrl-C to End)")

  <-c
  cacis.ExecCmd("killall -q hostapd", false)
  log.Println("[Debug] Terminated")
}

func genConfig(ssid, pw string) {
  log.Println("[Debug] Generate hostapd Config...")
  bytes, err := ioutil.ReadFile(hostapdConfTemplatePath)
  cacis.Error(err)
  config := string(bytes)

  config = strings.Replace(config, "{{SSID}}", ssid, 1)
  config = strings.Replace(config, "{{PASSWORD}}", pw, 1)

  ioutil.WriteFile(hostapdConfPath, []byte(config), 0644)
  log.Println("[Debug] Generated hostapd Config")
}

