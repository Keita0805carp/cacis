package ble

import (
  "fmt"
  "time"
  "strings"
  "github.com/google/uuid"
  "github.com/muka/go-bluetooth/hw"
  "github.com/muka/go-bluetooth/api/service"
  "github.com/muka/go-bluetooth/bluez/profile/agent"
)

const (
  adapterId = "hci0"
  UUID = "12345678-9012-3456-7890-abcdefabcdef"
)

func Main() {
  //UUID := genUUID()

  fmt.Println("ble master")
  addr := initialize()
  fmt.Printf("UUID: %s \nmyaddr: %s\n", UUID, addr)
  advertise(UUID, addr)
}

func genUUID() string {
  return uuid.New().String()
}

func initialize() string {
  btmgmt := hw.NewBtMgmt(adapterId)
  btmgmt.SetPowered(false)
  btmgmt.SetLe(true)
  btmgmt.SetBredr(false)
  btmgmt.SetPowered(true)

  adapter, err := hw.GetAdapter(adapterId)
  Error(err)
  return adapter.Address
}

func advertise(UUID, addr string) {
  //TODO UUID gen
  serviceID := UUID[4:8]
  options := service.AppOptions {
    AdapterID: adapterId,
    AgentCaps: agent.CapNoInputNoOutput,
    UUIDSuffix: UUID[8:],
    UUID:       UUID[:4],
  }

  ssid := strings.Replace(addr, ":", "", 5)
  pass := strings.Replace(UUID, "-", "", 4)

  app, err := service.NewApp(options)
  Error(err)
  defer app.Close()

  app.SetName("cacis-" + options.UUID + serviceID)

  service, err := app.NewService(serviceID)
  Error(err)
  err = app.AddService(service)
  Error(err)
  err = app.Run()
  Error(err)

  fmt.Printf("[INFO] SSID: %s\n", ssid)
  fmt.Printf("[INFO] PASS: %s\n", pass)

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
