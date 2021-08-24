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

func StartHostapd() {
  genConfig("ssidhoge", "passhoge")

  hoge, err := cacis.ExecCmd("hostapd " + hostapdConfPath)
  cacis.Error(err)
  fmt.Println(string(hoge))
}

func genConfig(ssid, pw string) {
  bytes, err := ioutil.ReadFile(hostapdConfTemplatePath)
  cacis.Error(err)

  rep := regexp.MustCompile(`ssid=.*`)
  bytes = rep.ReplaceAll(bytes, []byte("ssid="+ssid))
  rep = regexp.MustCompile(`wpa_passphrase=.*`)
  bytes = rep.ReplaceAll(bytes, []byte("wpa_passphrase="+pw))

  fmt.Println("Export config")

  ioutil.WriteFile(hostapdConfPath, bytes, 0644)
}

