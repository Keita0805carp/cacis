package slave

import (
  "fmt"
  "sort"
  "net"
  "os"
  "os/exec"
  "context"
  "strings"
  "encoding/json"

  "github.com/keita0805carp/cacis/cacis"

  "github.com/containerd/containerd"
)

const (
  masterIP = "192.168.56.250"
  masterPort = "27001"
  importDir = "./slave-vol/"
  containerdSock = "/var/snap/microk8s/common/run/containerd.sock"
  containerdNameSpace = "cacis"
)

func Main() {
  //slave()
  //checkSnapd()

  //recieveMicrok8sSnap()
  //installMicrok8sSnap()
  //clustering()
}

func slave() {
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
  //importAllImg(componentsList)
}

func recieveComponentsList() map[string]string {
  fmt.Println("Debug: [start] RECIEVE COMPONENTS LIST")
  // Socket
  conn, err := net.Dial("tcp", masterIP+":"+masterPort)
  Error(err)
  defer conn.Close()

  // Request Image
  fmt.Println("Debug: Request Components List")
  cLayer := cacis.RequestComponentsList()
  packet := cLayer.Marshal()
  fmt.Println(packet)
  conn.Write(packet)
  fmt.Println("Requested\n\n")


  // Recieve Packet
  fmt.Println("Debug: Recieve Packet")
  packet = make([]byte, cacis.CacisLayerSize)
  packetLength, err := conn.Read(packet)
  Error(err)
  fmt.Printf("Debug: Read Packet HEADER. len: %d\n", packetLength)
  cLayer = cacis.Unmarshal(packet)
  //fmt.Println(cLayer)

  fmt.Println("Debug: Read Packet PAYLOAD")
  //fmt.Println(cLayer)
  cLayer.Payload = loadPayload(conn, cLayer.Length)

  fmt.Printf("\rCompleted  %d\n", len(cLayer.Payload))

  var tmp map[string]string
  err = json.Unmarshal(cLayer.Payload, &tmp)
  Error(err)

  fmt.Println("Debug: [end] RECIEVE COMPONENTS LIST")
  return tmp
}

func recieveImg(s []string) {
  fmt.Println("Debug: [start] RECIEVE COMPONENT IMAGES")
  // Socket
  conn, err := net.Dial("tcp", masterIP+":"+masterPort)
  Error(err)
  defer conn.Close()

  // Request Image
  fmt.Println("Debug: Request Components Image")
  cLayer := cacis.RequestImage()
  packet := cLayer.Marshal()
  fmt.Println(packet)
  conn.Write(packet)
  fmt.Println("Requested\n\n")

  for _, fileName := range s {
    recieveFile(conn, fileName)
  }
  fmt.Println("Debug: [end] RECIEVE COMPONENT IMAGES")
}

func importImg(imageName, filePath string) {
  fmt.Println("Importing " + imageName + " from " + filePath + "...")

  ctx := context.Background()
  client, err := containerd.New(containerdSock, containerd.WithDefaultNamespace(containerdNameSpace))
  defer client.Close()
  Error(err)

  f, err := os.Open(filePath)
  defer f.Close()
  Error(err)

  opts := []containerd.ImportOpt{
    containerd.WithIndexName(imageName),
    //containerd.WithAllPlatforms(true),
  }
  client.Import(ctx, f, opts...)
  Error(err)
  fmt.Println("Imported")
}

func importAllImg(m map[string]string) {
  fmt.Printf("Import %d images for Kubernetes Components", len(m))
  for importFile, imageRef := range m {
    fmt.Printf("%s : %s\n", importDir + importFile, imageRef)
    fmt.Println("start")
    importImg(imageRef, importDir + importFile)
    fmt.Println("end\n")
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
  Error(err)
  defer conn.Close()

  // Request Image
  fmt.Println("Debug: Request snapd package")
  cLayer := cacis.RequestSnapd()
  packet := cLayer.Marshal()
  fmt.Println(packet)
  conn.Write(packet)
  fmt.Println("Requested\n\n")

  recieveFile(conn, "snapd.zip")

  fmt.Println("Debug: [end] RECIEVE COMPONENT IMAGES")
  myexec("dpkg -i ./*.deb")
  //reboot
  myexec("snap install core")
}

func recieveMicrok8sSnap() {
  fmt.Println("Debug: [start] RECIEVE SNAP FILES")
  s := []string{"microk8s_2346.assert", "microk8s_2346.snap", "core_11420.assert", "core_11420.snap"}
  // Socket
  conn, err := net.Dial("tcp", masterIP+":"+masterPort)
  Error(err)
  defer conn.Close()

  // Request Snap files
  fmt.Println("Debug: Request Snap files")
  cLayer := cacis.RequestMicrok8sSnap()
  packet := cLayer.Marshal()
  fmt.Println(packet)
  conn.Write(packet)
  fmt.Println("Requested\n\n")

  for _, fileName := range s {
    recieveFile(conn, fileName)
  }
  fmt.Println("Debug: [end] RECIEVE SNAP FILES")
}

func installMicrok8sSnap() {
  //TODO microk8s check
  fmt.Printf("Install microk8s via snap\n")
  fmt.Printf("Installing...")
  //myexec("snap install microk8s --classic")
  myexec("snap ack " + importDir + "microk8s_2346.assert")
  myexec("snap install " + importDir + "microk8s_2346.snap" + " --classic")
}

func clustering() {
  fmt.Println("Debug: [start] CLUSTERING")
  // Socket
  conn, err := net.Dial("tcp", masterIP+":"+masterPort)
  Error(err)
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
  Error(err)
  fmt.Printf("Debug: Read Packet HEADER. len: %d\n", packetLength)
  cLayer = cacis.Unmarshal(packet)

  // Recieve Packet PAYLOAD
  fmt.Println("Debug: Read Packet PAYLOAD")
  cLayer.Payload = loadPayload(conn, cLayer.Length)

  result, err := myexec(string(cLayer.Payload))
  fmt.Println(string(result))
  Error(err)

  fmt.Println("Debug: [end] CLUSTERING")
}

func leaveMicrok8s() {
  myexec("microk8s leave")
}

func removeMicrok8s() {
  myexec("microk8s stop")
  myexec("microk8s reset --destroy-storage")
  myexec("snap remove microk8s")
  myexec("apt purge microk8s")
}


func recieveFile(conn net.Conn, fileName string) {
  fmt.Printf("Debug: Recieve file '%s'\n", fileName)
  packet := make([]byte, cacis.CacisLayerSize)
  packetLength, err := conn.Read(packet)
  Error(err)
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
  Error(err)

  file.Write(cLayer.Payload)
}

func loadPayload(conn net.Conn, targetBytes uint64) []byte {
  packet := []byte{}
  recievedBytes := 0

  for len(packet) < int(targetBytes){
    buf := make([]byte, targetBytes - uint64(recievedBytes))
    packetLength, err := conn.Read(buf)
    Error(err)
    recievedBytes += packetLength
    packet = append(packet, buf[:packetLength]...)
    fmt.Printf("\rDebug: recieving...")
    fmt.Printf("\rCompleted  %d  of %d", len(packet), int(targetBytes))
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

func myexec(cmd string) ([]byte, error) {
  slice := strings.Split(cmd, " ")
  stdout, err := exec.Command(slice[0], slice[1:]...).Output()
  //fmt.Printf("exec: %s\noutput:\n%s", cmd, stdout)
  Error(err)
  return stdout, err
}

func Error(err error) {
  if err != nil {
    fmt.Println(err)
  }
}
