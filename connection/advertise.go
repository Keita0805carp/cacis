package connection

import (
  "log"

  "github.com/keita0805carp/cacis/cacis"

  "github.com/google/uuid"
  "github.com/muka/go-bluetooth/hw"
  "github.com/muka/go-bluetooth/api/service"
  "github.com/muka/go-bluetooth/bluez/profile/agent"
)

func Advertise(cancel chan struct{}, UUID, adapterAddr, adapterId, ssid, pass string) {
  log.Printf("[Debug] Start Advertise via Bluetooth\n")
  log.Printf("[Info]  Addr: %s\n", adapterAddr)
  log.Printf("[Info]  UUID: %s\n", UUID)

  serviceID := UUID[4:8]
  options := service.AppOptions {
    AdapterID: adapterId,
    AgentCaps: agent.CapNoInputNoOutput,
    UUIDSuffix: UUID[8:],
    UUID:       UUID[:4],
  }

  app, err := service.NewApp(options)
  cacis.Error(err)

  app.SetName("cacis-" + options.UUID + serviceID)

  service, err := app.NewService(serviceID)
  cacis.Error(err)
  err = app.AddService(service)
  cacis.Error(err)
  err = app.Run()
  cacis.Error(err)

  log.Printf("[Info]  SSID: %s\n", ssid)
  log.Printf("[Info]  PASS: %s\n", pass)

  timeout := uint32(65536 * 3600)
  log.Printf("[Debug] Advertising...\n")
  stop, err := app.Advertise(timeout)
  cacis.Error(err)

  <-cancel
  stop()
  app.Close()

  log.Printf("[Debug] Stop Advertise via Bluetooth\n")
}

func Initialize() (string, string) {
  log.Printf("[Debug] Initialize Bluetooth\n")

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

func GenUUID() string {
  log.Println("[Debug] Generate UUID")
  return uuid.New().String()
}
