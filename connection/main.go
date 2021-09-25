package connection

import (
  "time"
  "strings"
)

func Main(cancel chan struct{}) {
  adapterAddr, adapterId := Initialize()
  //UUID := connection.genUUID()
  UUID := "12345678-9012-3456-7890-abcdefabcdef"
  ssid := strings.Replace(adapterAddr, ":", "", 5)
  pass := strings.Replace(UUID, "-", "", 4)

  go StartHostapd(cancel, ssid, pass)

  go Advertise(cancel, UUID, adapterAddr, adapterId, ssid, pass)

  go DHCP(cancel)

  time.Sleep(3*time.Second)
  return
}
