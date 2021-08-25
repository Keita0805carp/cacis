package connection

import (
  "fmt"
  "strings"
  "io/ioutil"
  "os"
  "os/signal"

  "github.com/keita0805carp/cacis/cacis"
)

const (
  hostapdConfTemplatePath = "connection/hostapd-template.conf"
  hostapdConfPath = "connection/hostapd.conf"
)

func StartHostapd(ssid, pw string) {
  genConfig(ssid, pw)
  fmt.Printf("[INFO] SSID: %s\n", ssid)
  fmt.Printf("[INFO] PASS: %s\n", pw)

  cacis.ExecCmd("killall -q hostapd", false)

  fmt.Println("[DEBUG] Start hostapd in the Background")
  cacis.ExecCmd("hostapd -B " + hostapdConfPath, false)

  c := make(chan os.Signal, 1)
  signal.Notify(c, os.Interrupt)
  fmt.Println("Running hostpad... (Press Ctrl-C to End)")

  <-c
  cacis.ExecCmd("killall -q hostapd", false)
  fmt.Println("[DEBUG] Terminated")
}

func genConfig(ssid, pw string) {
  fmt.Println("[DEBUG] Generate hostapd Config...")
  bytes, err := ioutil.ReadFile(hostapdConfTemplatePath)
  cacis.Error(err)
  config := string(bytes)

  config = strings.Replace(config, "{{SSID}}", ssid, 1)
  config = strings.Replace(config, "{{PASSWORD}}", pw, 1)

  ioutil.WriteFile(hostapdConfPath, []byte(config), 0644)
  fmt.Println("[DEBUG] Generated hostapd Config")
}

