package connection

import (
  "fmt"
  "time"
  "strings"

  "github.com/keita0805carp/cacis/cacis"

  "github.com/google/uuid"
  "github.com/muka/go-bluetooth/hw"
  "github.com/muka/go-bluetooth/api/service"
  "github.com/muka/go-bluetooth/bluez/profile/agent"
)

func Advertise() {
  //UUID := genUUID()
  UUID := "12345678-9012-3456-7890-abcdefabcdef"

  fmt.Println("ble master")
  adapterAddr, adapterId := initialize()
  fmt.Printf("Addr: %s \nUUID: %s\n", adapterAddr, UUID)
  advertise(UUID, adapterAddr, adapterId)
}

func genUUID() string {
  return uuid.New().String()
}

func initialize() (string, string) {
  adaptersInfo, err := hw.GetAdapters()
  adapterInfo := adaptersInfo[0]
  cacis.Error(err)
  adapterId := adapterInfo.AdapterID
  adapterAddr := adapterInfo.Address

  btmgmt := hw.NewBtMgmt(adapterId)
  btmgmt.SetPowered(false)
  btmgmt.SetLe(true)
  btmgmt.SetBredr(false)
  btmgmt.SetPowered(true)

  return adapterAddr, adapterId
}

func advertise(UUID, adapterAddr, adapterId string) {
  //TODO UUID gen
  serviceID := UUID[4:8]
  options := service.AppOptions {
    AdapterID: adapterId,
    AgentCaps: agent.CapNoInputNoOutput,
    UUIDSuffix: UUID[8:],
    UUID:       UUID[:4],
  }

  ssid := strings.Replace(adapterAddr, ":", "", 5)
  pass := strings.Replace(UUID, "-", "", 4)

  app, err := service.NewApp(options)
  cacis.Error(err)
  defer app.Close()

  app.SetName("cacis-" + options.UUID + serviceID)

  service, err := app.NewService(serviceID)
  cacis.Error(err)
  err = app.AddService(service)
  cacis.Error(err)
  err = app.Run()
  cacis.Error(err)

  fmt.Printf("[INFO] SSID: %s\n", ssid)
  fmt.Printf("[INFO] PASS: %s\n", pass)

  timeout := uint32(6 * 3600) // 6h
  fmt.Printf("[DEBUG] Advertising for %ds...\n", timeout)
  cancel, err := app.Advertise(timeout)
  cacis.Error(err)

  defer cancel()
  time.Sleep(time.Duration(timeout) * time.Second)
}

func Error(err error) {
  if err != nil {
    fmt.Println(err)
  }
}