package connection

import (
  "fmt"
  "regexp"

  "github.com/keita0805carp/cacis/cacis"

  "github.com/muka/go-bluetooth/api"
  "github.com/muka/go-bluetooth/bluez/profile/adapter"
  "github.com/muka/go-bluetooth/bluez/profile/device"
)

func Discover() {
  fmt.Println("ble slave")
  slave()
  //sockettest()
}

func slave() {
  adapterID := adapter.GetDefaultAdapterID()
  ad, err := adapter.NewAdapter1FromAdapterID(adapterID)
  cacis.Error(err)
  err = ad.FlushDevices()
  cacis.Error(err)

  fmt.Printf("Discovering on %s\n", adapterID)

  dev, err := discover(ad)
  cacis.Error(err)
  p := dev.Properties
  fmt.Printf("Name: %s \n", p.Name)
  fmt.Printf("Address: %s (=SSID) \n", p.Address)
  fmt.Printf("UUID: %s (=PASS) \n", p.UUIDs[0])

}

func discover(a *adapter.Adapter1) (*device.Device1, error) {
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

    fmt.Printf("[Debug] Discovered (%s) %s\n", properties.Alias, properties.Address)

    isCacisNode := regexp.MustCompile(`^cacis-[0-9a-fA-F]{8}$`).MatchString(properties.Alias)

    //TODO regex cacis-xxxxxxxx
    if !isCacisNode {
      continue
    }

    return dev, nil
  }
  return nil, nil
}

