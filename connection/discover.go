package connection

import (
  "log"
  "regexp"
  "strings"

  "github.com/keita0805carp/cacis/cacis"

  "github.com/muka/go-bluetooth/api"
  "github.com/muka/go-bluetooth/bluez/profile/adapter"
  "github.com/muka/go-bluetooth/bluez/profile/device"
)

func GetWifiInfo() (string, string){
  adapterID := adapter.GetDefaultAdapterID()
  ad, err := adapter.NewAdapter1FromAdapterID(adapterID)
  cacis.Error(err)
  err = ad.FlushDevices()
  cacis.Error(err)

  log.Printf("[Debug] Discovering on %s\n", adapterID)

  dev, err := discoverBLE(ad)
  cacis.Error(err)
  p := dev.Properties
  ssid := strings.Replace(p.Address, ":", "", 5)
  pw := strings.Replace(p.UUIDs[0], "-", "", 4)
  log.Printf("[Info]  Name: %s \n", p.Name)
  log.Printf("[Info]  Address: %s (=SSID) \n", ssid)
  log.Printf("[Info]  UUID: %s (=PASS) \n", pw)

  return ssid, pw
}

func discoverBLE(a *adapter.Adapter1) (*device.Device1, error) {
  discoverd, cancel, err := api.Discover(a, nil)
  cacis.Error(err)
  defer cancel()

  for ev := range discoverd {

    dev, err := device.NewDevice1(ev.Path)
    cacis.Error(err)

    if dev == nil || dev.Properties == nil {
      continue
    }

    properties := dev.Properties

    log.Printf("[Debug] Discovered (%s) %s\n", properties.Alias, properties.Address)

    isCacisNode := regexp.MustCompile(`^cacis-[0-9a-fA-F]{8}$`).MatchString(properties.Alias)

    if !isCacisNode {
      continue
    }

    return dev, nil
  }
  return nil, nil
}

