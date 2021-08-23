package ble

import (
  "fmt"
  "regexp"

  "github.com/muka/go-bluetooth/api"
  "github.com/muka/go-bluetooth/bluez/profile/adapter"
  "github.com/muka/go-bluetooth/bluez/profile/device"
)

const (
  adapterID = "hci0"
)

func Main() {
  fmt.Println("ble slave")
  slave()
  //sockettest()
}

func slave() {
  ad, err := adapter.NewAdapter1FromAdapterID(adapterID)
  Error(err)
  err = ad.FlushDevices()
  Error(err)

  fmt.Printf("Discovering on %s\n", adapterID)

  dev, err := discover(ad)
  Error(err)
  p := dev.Properties
  fmt.Printf("Name: %s \n", p.Name)
  fmt.Printf("Address: %s (=SSID) \n", p.Address)
  fmt.Printf("UUID: %s (=PASS) \n", p.UUIDs[0])

}

func discover(a *adapter.Adapter1) (*device.Device1, error) {
  discoverd, cancel, err := api.Discover(a, nil)
  Error(err)
  defer cancel()

  for ev := range discoverd {

    dev, err := device.NewDevice1(ev.Path)
    Error(err)

    if dev == nil || dev.Properties == nil {
      continue
    }

    properties := dev.Properties

    fmt.Printf("[Debug] Discovered (%s) %s\n", properties.Alias, properties.Address)
    fmt.Println(regexp.MustCompile(`^cacis-[0-9a-fA-F]{8}$`).MatchString(properties.Alias))

    //TODO regex cacis-xxxxxxxx
    if properties.Name != "cacis-12345678" {
      continue
    }

    return dev, nil
  }
  return nil, nil
}

func Error(err error) {
  if err != nil {
    fmt.Println(err)
  }
}
