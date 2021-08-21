package ble

import (
  "fmt"
  "time"
  "github.com/muka/go-bluetooth/hw"
  "github.com/muka/go-bluetooth/api/service"
  "github.com/muka/go-bluetooth/bluez/profile/agent"
)

const (
  adapterId = "hci0"
)

func Main() {
  fmt.Println("ble master")
  initialize()
  advertise()
}

func initialize() {
  btmgmt := hw.NewBtMgmt(adapterId)
  btmgmt.SetPowered(false)
  btmgmt.SetLe(true)
  btmgmt.SetBredr(false)
  btmgmt.SetPowered(true)
}

func advertise() {
  options := service.AppOptions {
    AdapterID: adapterId,
    AgentCaps: agent.CapNoInputNoOutput,
    UUIDSuffix: "-0000-1000-8000-00805F9B34FB",
    UUID:       "1234",
  }

  app, err := service.NewApp(options)
  Error(err)
  defer app.Close()

  app.SetName("master-go")

  service1, err := app.NewService("2233")
  Error(err)
  err = app.AddService(service1)
  Error(err)

  err = app.Run()
  Error(err)
  fmt.Printf("[DEBUG] Exposed service %s\n", service1.Properties.UUID)

  timeout := uint32(6 * 3600) // 6h
  fmt.Printf("[DEBUG] Advertising for %ds...\n", timeout)
  cancel, err := app.Advertise(timeout)
  Error(err)

  defer cancel()
  time.Sleep(time.Duration(timeout) * time.Second)
}

func Error(err error) {
  if err != nil {
    fmt.Println(err)
  }
}
