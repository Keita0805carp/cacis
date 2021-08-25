package connection

import (
  "log"
  "time"
  "strings"

  "github.com/keita0805carp/cacis/cacis"

  "github.com/google/uuid"
  "github.com/muka/go-bluetooth/hw"
  "github.com/muka/go-bluetooth/api/service"
  "github.com/muka/go-bluetooth/bluez/profile/agent"
)

func Advertise() {
  adapterAddr, adapterId := initialize()
  //UUID := genUUID()
  UUID := "12345678-9012-3456-7890-abcdefabcdef"
  ssid := strings.Replace(adapterAddr, ":", "", 5)
  pass := strings.Replace(UUID, "-", "", 4)

  advertise(UUID, adapterAddr, adapterId, ssid, pass)

  StartHostapd(ssid, pass)
}

func genUUID() string {
  log.Println("[Debug] Generate UUID")
  return uuid.New().String()
}

func initialize() (string, string) {
  log.Println("[Debug] Initialize Bluetooth")

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

func advertise(UUID, adapterAddr, adapterId, ssid, pass string) {
  log.Println("[Debug] Start Advertise via Bluetooth")
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
  defer app.Close()

  app.SetName("cacis-" + options.UUID + serviceID)

  service, err := app.NewService(serviceID)
  cacis.Error(err)
  err = app.AddService(service)
  cacis.Error(err)
  err = app.Run()
  cacis.Error(err)

  log.Printf("[Info]  SSID: %s\n", ssid)
  log.Printf("[Info]  PASS: %s\n", pass)

  timeout := uint32(10) // 10s
  //timeout := uint32(6 * 3600) // 6h
  log.Printf("[Debug] Advertising for %ds...\n", timeout)
  cancel, err := app.Advertise(timeout)
  cacis.Error(err)

  defer cancel()
  time.Sleep(time.Duration(timeout) * time.Second)

  log.Println("[Debug] Stop Advertise via Bluetooth")
}
