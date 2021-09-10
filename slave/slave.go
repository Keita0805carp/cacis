package slave

import (
  "os"
  "fmt"
  "log"
  "net"
  "time"
  "context"
  "encoding/json"

  "github.com/keita0805carp/cacis/cacis"

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
  //TODO checkSnapd()
  //requestSnapd()
  //installSnapd()

  if !cacis.IsCommandAvailable("microk8s") {
    recieveMicrok8s()
    installMicrok8s()
  }

  setupMicrok8s()
  clustering()
  fmt.Println("[TEST] wait 30 seconds...")
  time.Sleep(30 * time.Second)
  unclustering()
}

func setupMicrok8s() {
  componentsList := recieveComponentsList()
  //fmt.Println(componentsList)
  // componentsList := map[string]string {
  //   "cni.img"             : "docker.io/calico/cni:v3.13.2",
  //   "pause.img"           : "k8s.gcr.io/pause:3.1",
  //   "kube-controllers.img": "docker.io/calico/kube-controllers:v3.13.2",
  //   "pod2daemon.img"      : "docker.io/calico/pod2daemon-flexvol:v3.13.2",
  //   "node.img"            : "docker.io/calico/node:v3.13.2",
  //   "coredns.img"         : "docker.io/coredns/coredns:1.8.0",
  //   "metrics-server.img"  : "k8s.gcr.io/metrics-server-arm64:v0.3.6",
  //   "dashboard.img"       : "docker.io/kubernetesui/dashboard:v2.0.0",
  // }
  sortedExportFileName := cacis.SortKeys(componentsList)
  recieveImg(sortedExportFileName)
  importAllImg(componentsList)
}

func recieveComponentsList() map[string]string {
  log.Println("[Debug] start: RECIEVE COMPONENTS LIST")
  // Socket
  conn, err := net.Dial("tcp", masterIP+":"+masterPort)
  cacis.Error(err)
  defer conn.Close()

  log.Println("[Debug] Request Components List")
  cLayer := cacis.RequestComponentsList()
  packet := cLayer.Marshal()
  log.Println(packet)
  conn.Write(packet)
  log.Printf("Requested\n")

  log.Println("[Debug] Recieve Packet")
  packet = make([]byte, cacis.CacisLayerSize)
  packetLength, err := conn.Read(packet)
  cacis.Error(err)
  log.Printf("[Debug] Read Packet HEADER. len: %d\n", packetLength)
  cLayer = cacis.Unmarshal(packet)
  log.Printf("[Debug] Read Packet PAYLOAD\n")
  cLayer.Payload = loadPayload(conn, cLayer.Length)

  log.Printf("\r[Debug] Completed  %d\n", len(cLayer.Payload))

  var tmpList map[string]string
  err = json.Unmarshal(cLayer.Payload, &tmpList)
  cacis.Error(err)

  log.Println("[Debug] end: RECIEVE COMPONENTS LIST")
  return tmpList
}

func recieveImg(s []string) {
  log.Println("[Debug] start: RECIEVE COMPONENT IMAGES")
  // Socket
  conn, err := net.Dial("tcp", masterIP+":"+masterPort)
  cacis.Error(err)
  defer conn.Close()

  log.Println("[Debug] Request Components Image")
  cLayer := cacis.RequestImage()
  packet := cLayer.Marshal()
  log.Println(packet)
  conn.Write(packet)
  log.Printf("[Debug] Requested\n\n")

  for _, fileName := range s {
    recieveFile(conn, fileName)
  }
  log.Println("[Debug] end: RECIEVE COMPONENT IMAGES")
}

func importImg(imageName, filePath string) {
  log.Println("[Debug] Importing " + imageName + " from " + filePath + "...")

  ctx := context.Background()
  client, err := containerd.New(containerdSock, containerd.WithDefaultNamespace(containerdNameSpace))
  defer client.Close()
  cacis.Error(err)

  f, err := os.Open(filePath)
  defer f.Close()
  cacis.Error(err)

  opts := []containerd.ImportOpt{
    containerd.WithIndexName(imageName),
    //containerd.WithAllPlatforms(true),
  }
  client.Import(ctx, f, opts...)
  cacis.Error(err)
  log.Println("[Debug] Imported")
}

func importAllImg(m map[string]string) {
  log.Printf("[Debug] Import %d images for Kubernetes Components", len(m))
  for importFile, imageRef := range m {
    log.Printf("%s : %s\n", targetDir + importFile, imageRef)
    log.Println("[Info]  start")
    importImg(imageRef, targetDir + importFile)
    log.Printf("[Info]  end\n")
    }
}

func clustering() {
  log.Println("[Debug] Start CLUSTERING")
  // Socket
  conn, err := net.Dial("tcp", masterIP+":"+masterPort)
  cacis.Error(err)
  defer conn.Close()

  log.Println("[Debug] Request Clustering")
  cLayer := cacis.RequestClustering()
  packet := cLayer.Marshal()
  //fmt.Println(packet)
  conn.Write(packet)
  log.Printf("Requested\n\n")

  packet = make([]byte, cacis.CacisLayerSize)
  packetLength, err := conn.Read(packet)
  cacis.Error(err)
  log.Printf("[Debug] Read Packet HEADER. len: %d\n", packetLength)
  cLayer = cacis.Unmarshal(packet)
  log.Println("Debug: Read Packet PAYLOAD")
  cLayer.Payload = loadPayload(conn, cLayer.Length)

  log.Printf("\n[Debug] Clustering...\n")
  result, err := cacis.ExecCmd(string(cLayer.Payload), true)
  fmt.Println(string(result))
  cacis.Error(err)

  log.Println("[Debug] End CLUSTERING")
}

func unclustering() {
  log.Println("[Debug] Start UNCLUSTERING")
  // Socket
  conn, err := net.Dial("tcp", masterIP+":"+masterPort)
  cacis.Error(err)
  defer conn.Close()

  hostname, err := os.Hostname()
  cacis.Error(err)

  log.Println("[Debug] Request Clustering")
  cLayer := cacis.RequestUnclustering([]byte(hostname))
  packet := cLayer.Marshal()
  //fmt.Println(packet)
  conn.Write(packet)
  fmt.Printf("Requested\n\n")
  cacis.ExecCmd("microk8s leave", true)
  log.Println("[Debug] End UNCLUSTERING")
}

func recieveFile(conn net.Conn, fileName string) {
  log.Printf("[Debug] Recieve file '%s'\n", fileName)
  packet := make([]byte, cacis.CacisLayerSize)
  packetLength, err := conn.Read(packet)
  cacis.Error(err)
  log.Printf("[Debug] Read Packet HEADER. len: %d\n", packetLength)
  cLayer := cacis.Unmarshal(packet)
  log.Println("[Debug] Read Packet PAYLOAD")
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

