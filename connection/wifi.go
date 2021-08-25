package connection

import (
  "fmt"
  "strings"
  "io/ioutil"

  "github.com/keita0805carp/cacis/cacis"

  //"github.com/mdlayher/wifi"
  //"github.com/theojulienne/go-wireless"
)

const (
  netplanConfTemplatePath = "connection/netplan-template.conf"
  //netplanConfPath = "connection/60-cacis.yaml"
  netplanConfPath = "/etc/netplan/60-cacis.yaml"
)

func Connect(ssid, pw string) {
  genNetplanConfig(ssid, pw)
  fmt.Printf("[INFO] SSID: %s\n", ssid)
  fmt.Printf("[INFO] PASS: %s\n", pw)

  fmt.Println("[DEBUG] Apply netplan config")
  cacis.ExecCmd("netplan apply", false)
  fmt.Println("[DEBUG] Applied netplan config")
}

func genNetplanConfig(ssid, pw string) {

  fmt.Println("[DEBUG] Generate netplan Config...")
  bytes, err := ioutil.ReadFile(netplanConfTemplatePath)
  cacis.Error(err)
  config := string(bytes)

  config = strings.Replace(config, "{{SSID}}", ssid, 1)
  config = strings.Replace(config, "{{PASSWORD}}", pw, 1)

  ioutil.WriteFile(netplanConfPath, []byte(config), 0644)
  fmt.Println("[DEBUG] Generated netplan Config")
}

func test() {
  //ssid := "DCA6326E3A9D"
  //pw := "12345678901234567890abcdefabcdef"

  //ifaces := wireless.Interfaces()
  //fmt.Println(ifaces)

  //wc, err := wireless.NewClient("wlan0")
  //cacis.Error(err)
  //defer wc.Close()
  //fmt.Println(wc)

  //dial, err := wireless.Dial("wlan0")
  //cacis.Error(err)

  // c, err := wifi.New()
  // cacis.Error(err)
  // defer c.Close()
  //
  // ifi, err := c.Interfaces()
  // cacis.Error(err)
  //
  // err = c.Connect(ifi[0], ssid)
  // cacis.Error(err)
  //
  //
  //err = c.Connect(
}
