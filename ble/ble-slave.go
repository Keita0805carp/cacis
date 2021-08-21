package ble

import (
  "fmt"
  "time"
  "strings"

  "github.com/godbus/dbus/v5"
  "github.com/muka/go-bluetooth/api"
  "github.com/muka/go-bluetooth/bluez/profile/adapter"
  "github.com/muka/go-bluetooth/bluez/profile/agent"
  "github.com/muka/go-bluetooth/bluez/profile/device"
)

const (
  adapterId = "hci0"
  me = "DC:A6:32:6E:43:D9"
  hwaddr = "DC:A6:32:6E:3A:9D"
)

func Main() {
  fmt.Println("ble slave")
  slave()
}

func slave() {
  fmt.Printf("Discovering %s on %s\n", hwaddr, adapterId)

  ad, ag := initialize()

  dev, err := discover(ad, hwaddr)
  Error(err)

  err = connect(dev, ag, adapterId)
  Error(err)

  time.Sleep(5 * time.Second)

  disconnect(dev)
  dev.Close()
}

func initialize() (*adapter.Adapter1, *agent.SimpleAgent) {
  ad, err := adapter.NewAdapter1FromAdapterID(adapterId)
  Error(err)

  //Connect DBus System bus
  conn, err := dbus.SystemBus()
  Error(err)

  // do not reuse agent0 from service
  agent.NextAgentPath()

  ag := agent.NewSimpleAgent()
  err = agent.ExposeAgent(conn, ag, agent.CapNoInputNoOutput, true)
  Error(err)

  return ad, ag
}

func discover(a *adapter.Adapter1, hwaddr string) (*device.Device1, error) {
  err := a.FlushDevices()
  Error(err)

  discovery, cancel, err := api.Discover(a, nil)
  Error(err)
  defer cancel()

  for ev := range discovery {

    dev, err := device.NewDevice1(ev.Path)
    Error(err)

    if dev == nil || dev.Properties == nil {
      continue
    }

    p := dev.Properties

    n := p.Alias
    if p.Name != "" {
      n = p.Name
    }
    fmt.Printf("[Debug] Discovered (%s) %s\n", n, p.Address)

    if p.Address != hwaddr {
      continue
    }

    return dev, nil
  }

  return nil, nil
}

func connect(dev *device.Device1, ag *agent.SimpleAgent, adapterID string) error {

  props, err := dev.GetProperties()
  if err != nil {
    return fmt.Errorf("Failed to load props: %s", err)
  }

  fmt.Printf("Found device name=%s addr=%s rssi=%d\n", props.Name, props.Address, props.RSSI)

  if props.Connected {
    fmt.Println("Device is connected")
    return nil
  }

  if !props.Paired || !props.Trusted {
    fmt.Println("Pairing device")

    err := dev.Pair()
    if err != nil {
      return fmt.Errorf("Pair failed: %s", err)
    }

    fmt.Printf("[Debug] Pair succeed, connecting...\n")
    agent.SetTrusted(adapterID, dev.Path())
  }

  if !props.Connected {
    fmt.Println("Connecting device")
    err = dev.Connect()
    if err != nil {
      if !strings.Contains(err.Error(), "Connection refused") {
        return fmt.Errorf("Connect failed: %s", err)
      }
    }
    fmt.Println("Connected")
  }

  return nil
}

func disconnect(dev *device.Device1) {
  fmt.Println("Disconnecting...")
  dev.Disconnect()
  fmt.Println("Disconnected")
}

func Error(err error) {
  if err != nil {
    fmt.Println(err)
  }
}
