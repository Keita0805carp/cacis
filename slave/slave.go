package slave

import (
  "os"
  "fmt"
  "log"
  "net"
  "time"
  "encoding/json"

  "github.com/keita0805carp/cacis/cacis"
  "github.com/keita0805carp/cacis/master"

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
  //TODO? checkSnapd()
  //requestSnapd()
  //installSnapd()

  listen, err := net.Listen("tcp", ":27001")
  cacis.Error(err)

  if !cacis.IsCommandAvailable("microk8s") {
    recieveMicrok8s(listen)
    installMicrok8s()
    setupMicrok8s(listen)
  }

  waitReadyMicrok8s()
  clustering(listen)
  fmt.Printf("[TEST] wait 60 seconds...\n")
  time.Sleep(60 * time.Second)
  unclustering()
  //TODO remove microk8s
}

func setupMicrok8s(listen net.Listener) {
  componentsList := recieveComponentsList(listen)
  sortedExportFileName := cacis.SortKeys(componentsList)
  recieveImg(listen, sortedExportFileName)
  importAllImg(componentsList)
}

func recieveComponentsList(listen net.Listener) map[string]string {
  log.Printf("[Debug] start: RECIEVE COMPONENTS LIST\n")
  // Socket
  conn2master, err := net.Dial("tcp", masterIP+":"+masterPort)
  cacis.Error(err)
  log.Printf("[Debug] Request Components List\n")
  cLayer := cacis.RequestComponentsList()
  packet := cLayer.Marshal()
  //fmt.Println(packet)
  conn2master.Write(packet)
  log.Printf("Requested\n\n")
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

  log.Printf("[Debug] end: RECIEVE COMPONENTS LIST\n")
  conn2slave.Close()
  return tmpList
}

func recieveImg(listen net.Listener, s []string) {
  log.Printf("[Debug] start: RECIEVE COMPONENT IMAGES\n")
  // Socket
  conn2master, err := net.Dial("tcp", masterIP+":"+masterPort)
  cacis.Error(err)
  log.Printf("[Debug] Request Components Image\n")
  cLayer := cacis.RequestImage()
  packet := cLayer.Marshal()
  //fmt.Println(packet)
  conn2master.Write(packet)
  log.Printf("Requested\n\n")
  conn2master.Close()

  conn2slave, err := listen.Accept()
  cacis.Error(err)

  for _, fileName := range s {
    recieveFile(conn2slave, fileName)
  }
  log.Printf("[Debug] end: RECIEVE COMPONENT IMAGES\n")
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

func clustering(listen net.Listener) {
  log.Printf("[Debug] Start CLUSTERING\n")
  // Socket
  conn2master, err := net.Dial("tcp", masterIP+":"+masterPort)
  cacis.Error(err)
  log.Printf("[Debug] Request Clustering\n")
  cLayer := cacis.RequestClustering()
  packet := cLayer.Marshal()
  //fmt.Println(packet)
  conn2master.Write(packet)
  log.Printf("Requested\n\n")
  conn2master.Close()

  conn2slave, err := listen.Accept()
  cacis.Error(err)

  packet = make([]byte, cacis.CacisLayerSize)
  packetLength, err := conn2slave.Read(packet)
  cacis.Error(err)
  log.Printf("[Debug] Read Packet HEADER. len: %d\n", packetLength)
  cLayer = cacis.Unmarshal(packet)
  log.Printf("Debug: Read Packet PAYLOAD\n")
  cLayer.Payload = loadPayload(conn2slave, cLayer.Length)

  log.Printf("\n[Debug] Clustering...\n")
  result, err := cacis.ExecCmd(string(cLayer.Payload), true)
  fmt.Println(string(result))
  cacis.Error(err)

  log.Printf("[Debug] End CLUSTERING\n")
  conn2slave.Close()
}

func unclustering() {
  log.Printf("[Debug] Start UNCLUSTERING\n")
  log.Printf("[Debug] Leaving...\n")
  cacis.ExecCmd("microk8s leave", true)
  // Socket
  conn, err := net.Dial("tcp", masterIP+":"+masterPort)
  cacis.Error(err)
  defer conn.Close()

  hostname, err := os.Hostname()
  cacis.Error(err)

  log.Printf("[Debug] Request Unclustering\n")
  cLayer := cacis.RequestUnclustering([]byte(hostname))
  packet := cLayer.Marshal()
  //fmt.Println(packet)
  conn.Write(packet)
  log.Printf("[Debug] End UNCLUSTERING\n")
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
  return packet
}

