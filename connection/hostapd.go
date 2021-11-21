package connection

import (
  "log"
  "strings"
  "io/ioutil"

  "github.com/keita0805carp/cacis/cacis"
)

const (
  hostapdConfTemplatePath = cacis.HostapdConfTemplatePath
  hostapdConfPath = cacis.HostapdConfPath
)

func StartHostapd(cancel chan struct{}, ssid, pw string) {
  ipSet(cacis.MasterInterface)
  genConfig(ssid, pw)
  log.Printf("[Info]  SSID: %s\n", ssid)
  log.Printf("[Info]  PASS: %s\n", pw)

  cacis.ExecCmd("killall -q hostapd", false)

  log.Println("[Debug] Start hostapd in the Background")
  cacis.ExecCmd("hostapd -B " + hostapdConfPath, false)

  <-cancel
  cacis.ExecCmd("killall -q hostapd", false)
  log.Printf("[Debug] Terminated hostapd\n")
}

func genConfig(ssid, pw string) {
  log.Printf("[Debug] Generate hostapd Config...\n")
  bytes, err := ioutil.ReadFile(hostapdConfTemplatePath)
  cacis.Error(err)
  config := string(bytes)

  config = strings.Replace(config, "{{SSID}}", ssid, 1)
  config = strings.Replace(config, "{{PASSWORD}}", pw, 1)

  ioutil.WriteFile(hostapdConfPath, []byte(config), 0644)
  log.Printf("[Debug] Generated hostapd Config\n")
}

func ipSet(iface string) {
  cacis.ExecCmd("ifconfig " + iface + " " + cacis.MasterIP + "/24" , false)
  cacis.ExecCmd("route add default gw 172.30.0.254" , false) // Containerd needs default gateway to start
}

