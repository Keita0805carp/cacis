package slave

import (
  "fmt"
  "log"
  "sort"
  "net"
  "os"
  "context"
  "encoding/json"

  "github.com/keita0805carp/cacis/cacis"

  "github.com/containerd/containerd"
)

const (
  masterIP = "192.168.56.250"
  masterPort = "27001"
  importDir = "./slave-vol/"
  containerdSock = "/var/snap/microk8s/common/run/containerd.sock"
  containerdNameSpace = "k8s.io"
  //containerdNameSpace = "cacis"
)

func Main() {
  //TODO
  //checkSnapd()
  //requestSnapd()
  //installSnapd()

  //recieveMicrok8sSnap()
  //installMicrok8sSnap()

  //fmt.Println("wait 3 seconds...")
  //time.Sleep(5 * time.Second)
  //setupMicrok8s()
  //clustering()
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
  sortedExportFileName := sortKeys(componentsList)
  recieveImg(sortedExportFileName)
  importAllImg(componentsList)
}

func recieveComponentsList() map[string]string {
  log.Println("[Debug] start: RECIEVE COMPONENTS LIST")
  // Socket
  conn, err := net.Dial("tcp", masterIP+":"+masterPort)
  cacis.Error(err)
  defer conn.Close()

  // Request Image
  log.Println("[Debug] Request Components List")
  cLayer := cacis.RequestComponentsList()
  packet := cLayer.Marshal()
  log.Println(packet)
  conn.Write(packet)
  log.Println("Requested\n\n")


  // Recieve Packet
  log.Println("[Debug] Recieve Packet")
  packet = make([]byte, cacis.CacisLayerSize)
  packetLength, err := conn.Read(packet)
  cacis.Error(err)
  log.Printf("[Debug] Read Packet HEADER. len: %d\n", packetLength)
  cLayer = cacis.Unmarshal(packet)
  //fmt.Println(cLayer)

  log.Println("[Debug] Read Packet PAYLOAD")
  //fmt.Println(cLayer)
  cLayer.Payload = loadPayload(conn, cLayer.Length)

  log.Printf("\r[Debug] Completed  %d", len(cLayer.Payload))

  var tmp map[string]string
  err = json.Unmarshal(cLayer.Payload, &tmp)
  cacis.Error(err)

  log.Println("[Debug] end: RECIEVE COMPONENTS LIST")
  return tmp
}

func recieveImg(s []string) {
  log.Println("[Debug] start: RECIEVE COMPONENT IMAGES")
  // Socket
  conn, err := net.Dial("tcp", masterIP+":"+masterPort)
  cacis.Error(err)
  defer conn.Close()

  // Request Image
  log.Println("[Debug] Request Components Image")
  cLayer := cacis.RequestImage()
  packet := cLayer.Marshal()
  log.Println(packet)
  conn.Write(packet)
  log.Println("[Debug] Requested\n\n")

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
    log.Printf("%s : %s\n", importDir + importFile, imageRef)
    log.Println("[Info]  start")
    importImg(imageRef, importDir + importFile)
    log.Println("[Info]  end\n")
    }
}


func snapd() {
  //TODO
  /// if !snap
    /// if debian(raspberry pi os)
      //recieve zip
      //unzip

  fmt.Println("Debug: [start] RECIEVE snapd and install")
  // Socket
  conn, err := net.Dial("tcp", masterIP+":"+masterPort)
  cacis.Error(err)
  defer conn.Close()

  // Request Image
  fmt.Println("Debug: Request snapd package")
  cLayer := cacis.RequestSnapd()
  packet := cLayer.Marshal()
  fmt.Println(packet)
  conn.Write(packet)
  fmt.Printf("Requested\n\n")

  recieveFile(conn, "snapd.zip")

  fmt.Println("Debug: [end] RECIEVE COMPONENT IMAGES")
  cacis.ExecCmd("dpkg -i ./*.deb", false)
  //reboot
  cacis.ExecCmd("snap install core", false)
}

func recieveMicrok8sSnap() {
  fmt.Println("Debug: [start] RECIEVE SNAP FILES")
  s := []string{"microk8s_2407.assert", "microk8s_2407.snap", "core_11420.assert", "core_11420.snap"}
  // Socket
  conn, err := net.Dial("tcp", masterIP+":"+masterPort)
  cacis.Error(err)
  defer conn.Close()

  // Request Snap files
  fmt.Println("Debug: Request Snap files")
  cLayer := cacis.RequestMicrok8sSnap()
  packet := cLayer.Marshal()
  fmt.Println(packet)
  conn.Write(packet)
  fmt.Printf("Requested\n\n")

  for _, fileName := range s {
    recieveFile(conn, fileName)
  }
  fmt.Println("Debug: [end] RECIEVE SNAP FILES")
}

func installMicrok8sSnap() {
  //TODO microk8s check
  fmt.Printf("Install microk8s via snap\n")
  fmt.Printf("Installing...")
  //cacis.ExecCmd("snap install microk8s --classic")
  cacis.ExecCmd("snap ack " + importDir + "microk8s_2407.assert", false)
  cacis.ExecCmd("snap install " + importDir + "microk8s_2407.snap" + " --classic", false)
}

func clustering() {
  fmt.Println("Debug: [start] CLUSTERING")
  // Socket
  conn, err := net.Dial("tcp", masterIP+":"+masterPort)
  cacis.Error(err)
  defer conn.Close()

  // Request Snap files
  fmt.Println("Debug: Request Clustering")
  cLayer := cacis.RequestClustering()
  packet := cLayer.Marshal()
  fmt.Println(packet)
  conn.Write(packet)
  fmt.Println("Requested\n\n")

  packet = make([]byte, cacis.CacisLayerSize)
  packetLength, err := conn.Read(packet)
  cacis.Error(err)
  fmt.Printf("Debug: Read Packet HEADER. len: %d\n", packetLength)
  cLayer = cacis.Unmarshal(packet)

  // Recieve Packet PAYLOAD
  fmt.Println("Debug: Read Packet PAYLOAD")
  cLayer.Payload = loadPayload(conn, cLayer.Length)

  fmt.Printf("\nDebug: Clustering...\n")
  result, err := cacis.ExecCmd(string(cLayer.Payload), false)
  fmt.Println(string(result))
  cacis.Error(err)

  fmt.Println("Debug: [end] CLUSTERING")
}

func unclustering() {
  //TODO request leave
  fmt.Println("Debug: [start] UNCLUSTERING")
  // Socket
  conn, err := net.Dial("tcp", masterIP+":"+masterPort)
  cacis.Error(err)
  defer conn.Close()

  hostname, err := os.Hostname()
  cacis.Error(err)

  // Request Snap files
  fmt.Println("Debug: Request Clustering")
  cLayer := cacis.RequestUnclustering([]byte(hostname))
  packet := cLayer.Marshal()
  fmt.Println(packet)
  conn.Write(packet)
  fmt.Printf("Requested\n\n")
  //TODO tell hostname to master
  cacis.ExecCmd("microk8s leave", true)
  fmt.Println("Debug: [end] UNCLUSTERING")
}

func removeMicrok8s() {
  cacis.ExecCmd("microk8s stop", false)
  cacis.ExecCmd("microk8s reset --destroy-storage", false)
  cacis.ExecCmd("snap remove --purge microk8s", false)
  cacis.ExecCmd("apt purge snap", false)
}


func recieveFile(conn net.Conn, fileName string) {
  fmt.Printf("Debug: Recieve file '%s'\n", fileName)
  packet := make([]byte, cacis.CacisLayerSize)
  packetLength, err := conn.Read(packet)
  cacis.Error(err)
  fmt.Printf("Debug: Read Packet HEADER. len: %d\n", packetLength)
  cLayer := cacis.Unmarshal(packet)

  // Recieve Packet PAYLOAD
  fmt.Println("Debug: Read Packet PAYLOAD")
  cLayer.Payload = loadPayload(conn, cLayer.Length)

  fmt.Printf("\rCompleted  %d\n", len(cLayer.Payload))

  fmt.Printf("Debug: Write file '%s'\n\n", fileName)
  // File
  filePath := importDir + fileName
  //filePath := "./hoge1.txt"
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

func sortKeys(m map[string]string) []string {
  ///sort
  sorted := make([]string, len(m))
  index := 0
  for key := range m {
        sorted[index] = key
        index++
    }
    sort.Strings(sorted)
  /*
  for _, exportFile := range exportFileNameSort {
    fmt.Printf("%-20s : %s\n", exportFile, componentsList[exportFile])
  }
  */
  return sorted
}

