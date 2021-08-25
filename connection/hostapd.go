package connection

import (
  "fmt"
  "strings"
  "io/ioutil"

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

  fmt.Println("[DEBUG] Start hostapd")
  hoge, err := cacis.ExecCmd("hostapd " + hostapdConfPath)
  cacis.Error(err)
  fmt.Println(string(hoge))
  fmt.Println("[DEBUG] Terminating...")
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

