package connection

import (
  "fmt"
  "regexp"
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

  rep := regexp.MustCompile(`ssid=.*`)
  bytes = rep.ReplaceAll(bytes, []byte("ssid="+ssid))
  rep = regexp.MustCompile(`wpa_passphrase=.*`)
  bytes = rep.ReplaceAll(bytes, []byte("wpa_passphrase="+pw))

  ioutil.WriteFile(hostapdConfPath, bytes, 0644)
  fmt.Println("[DEBUG] Generated hostapd Config")
}

