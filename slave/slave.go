package slave

import (
  "os"
  "fmt"
  "log"
  "net"
  "encoding/json"

  "github.com/keita0805carp/cacis/cacis"
  "github.com/keita0805carp/cacis/master"
  "github.com/keita0805carp/cacis/connection"

  "github.com/containerd/containerd"
)

const (
  masterIP = cacis.MasterIP
  masterPort = cacis.MasterPort
  targetDir = cacis.TargetDir
  containerdSock = cacis.ContainerdSock
  containerdNameSpace = cacis.ContainerdNameSpace
)

func Main() {
  for {
    log.Printf("[Debug]: Run Main Slave Process\n")
    cancel := make(chan struct{})
    listen, err := net.Listen("tcp", ":27001")
    cacis.Error(err)

    ssid, pw := connection.GetWifiInfo()
    connection.Connect(ssid, pw)

    if !cacis.IsCommandAvailable("microk8s") {
      cacis.CreateTempDir()
      recieveMicrok8s(listen)
      installMicrok8s()
      setupMicrok8s(listen)
    }

    WaitReadyMicrok8s()
    clustering(listen)
    labelNode()

    go connection.UnstableWifiEvent(cancel)

    <- cancel

    Unclustering()
    connection.Disconnect()
    WaitReadyMicrok8s()
  }
}

func setupMicrok8s(listen net.Listener) {
  componentsList := recieveComponentsList(listen)
  sortedExportFileName := cacis.SortKeys(componentsList)
  recieveImg(listen, sortedExportFileName)
  importAllImg(componentsList)
}

func recieveComponentsList(listen net.Listener) map[string]string {
  log.Printf("[Debug] Start Recieve Components List\n")
  // Socket
  conn2master, err := net.Dial("tcp", masterIP+":"+masterPort)
  cacis.Error(err)
  log.Printf("[Debug] Request Components List\n")
  cLayer := cacis.RequestComponentsList()
  packet := cLayer.Marshal()
  //fmt.Println(packet)
  conn2master.Write(packet)
  log.Printf("Requested\n")
  conn2master.Close()

  conn2slave, err := listen.Accept()
  cacis.Error(err)

  log.Printf("[Debug] Recieve Packet\n")
  packet = make([]byte, cacis.CacisLayerSize)
  packetLength, err := conn2slave.Read(packet)
  cacis.Error(err)
  log.Printf("[Debug] Read Packet HEADER. len: %d\n", packetLength)
  cLayer = cacis.Unmarshal(packet)
  log.Printf("[Debug] Read Packet PAYLOAD\n")
  cLayer.Payload = loadPayload(conn2slave, cLayer.Length)

  log.Printf("\r\n[Debug] Completed  %d\n", len(cLayer.Payload))

  var tmpList map[string]string
  err = json.Unmarshal(cLayer.Payload, &tmpList)
  cacis.Error(err)

  log.Printf("[Debug] End Recieve Components List\n")
  conn2slave.Close()
  return tmpList
}

func recieveImg(listen net.Listener, s []string) {
  log.Printf("[Debug] Start Recieve Component Images\n")
  // Socket
  conn2master, err := net.Dial("tcp", masterIP+":"+masterPort)
  cacis.Error(err)
  log.Printf("[Debug] Request Components Image\n")
  cLayer := cacis.RequestImage()
  packet := cLayer.Marshal()
  //fmt.Println(packet)
  conn2master.Write(packet)
  log.Printf("Requested\n")
  conn2master.Close()

  conn2slave, err := listen.Accept()
  cacis.Error(err)

  for _, fileName := range s {
    recieveFile(conn2slave, fileName)
  }
  log.Printf("[Debug] End Recieve Component Images\n")
  conn2slave.Close()
}

func importImg(imageName, filePath string) {
  log.Printf("[Debug] Importing " + imageName + " from " + filePath + "...\n")

  ctx, client := master.ContainerdInit()

  f, err := os.Open(filePath)
  defer f.Close()
  cacis.Error(err)

  opts := []containerd.ImportOpt{
    containerd.WithIndexName(imageName),
    //containerd.WithAllPlatforms(true),
  }
  client.Import(ctx, f, opts...)
  cacis.Error(err)
  log.Printf("[Debug] Imported\n")
}

func importAllImg(m map[string]string) {
  log.Printf("[Debug] Import %d images for Kubernetes Components", len(m))
  for importFile, imageRef := range m {
    log.Printf("%s : %s\n", targetDir + importFile, imageRef)
    log.Printf("[Info]  start\n")
    importImg(imageRef, targetDir + importFile)
    log.Printf("[Info]  end\n")
    }
}

func recieveFile(conn net.Conn, fileName string) {
  log.Printf("[Debug] Recieve file '%s'\n", fileName)
  packet := make([]byte, cacis.CacisLayerSize)
  packetLength, err := conn.Read(packet)
  cacis.Error(err)
  log.Printf("[Debug] Read Packet HEADER. len: %d\n", packetLength)
  cLayer := cacis.Unmarshal(packet)
  log.Printf("[Debug] Read Packet PAYLOAD\n")
  cLayer.Payload = loadPayload(conn, cLayer.Length)

  log.Printf("\rCompleted  %d\n", len(cLayer.Payload))

  log.Printf("Debug: Write file '%s'\n\n", fileName)
  // File
  filePath := targetDir + fileName
  file , err := os.Create(filePath)
  cacis.Error(err)

  file.Write(cLayer.Payload)
}

func loadPayload(conn net.Conn, targetBytes uint64) []byte {
  packet := []byte{}
  recievedBytes := 0

  for len(packet) < int(targetBytes){
    buf := make([]byte, targetBytes - uint64(recievedBytes))
    packetLength, err := conn.Read(buf)
    cacis.Error(err)
    recievedBytes += packetLength
    packet = append(packet, buf[:packetLength]...)
    //log.Printf("\r[Debug] recieving...")
    fmt.Printf("\r[Info]  Completed  %d  of %d", len(packet), int(targetBytes))
  }
  fmt.Printf("\r[Info]  Completed  %d  of %d\n", len(packet), int(targetBytes))
  return packet
}

