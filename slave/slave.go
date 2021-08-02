package slave

import (
  "fmt"
  "sort"
  //"io"
  "net"
  "os"
  //"time"
  //"strconv"
  //"strings"
  "context"
  "encoding/json"

  "github.com/keita0805carp/cacis/cacis"

  "github.com/containerd/containerd"
)

const (
  MASTER = "localhost:27001"
  //importDir = "./master-vol/"
  importDir = "./slave-vol/"
)

func Main() {
  slave()
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
  importAllImg(componentsList)
}

func recieveComponentsList() map[string]string {
  // Socket
  conn, err := net.Dial("tcp", MASTER)
  Error(err)
  defer conn.Close()

  // Request Image
  fmt.Println("Debug: Create Request Packet")
  cLayer := cacis.RequestComponentsList()
  packet := cLayer.Marshal()
  fmt.Println(packet)
  conn.Write(packet)
  fmt.Println("Send\n\n")


  // Recieve Packet HEAD
  fmt.Println("Debug: Recieve Packet HEAD")
  packet = make([]byte, cacis.CacisLayerSize)
  packetLength, err := conn.Read(packet)
  Error(err)
  cLayer = cacis.Unmarshal(packet)
  //fmt.Println(cLayer)
  //fmt.Println(string(cLayer.Payload))

  // Recieve Packet PAYLOAD
  fmt.Println("Debug: Recieve Packet PAYLOAD")
  //fmt.Println(cLayer)
  packet = []byte{}
  recievedBytes := 0
  targetBytes := cLayer.Length

  for {
    buf := make([]byte, targetBytes - uint64(recievedBytes))
    packetLength, err = conn.Read(buf)
    Error(err)
    recievedBytes += packetLength
    packet = append(packet, buf[:packetLength]...)
    fmt.Println("Debug: recieving...")
    fmt.Printf("Complete %d of %d\n", len(packet), int(targetBytes))
    if len(packet) == int(targetBytes) {
      break
    }
  }
  //cLayer.Payload = packet

  var tmp map[string]string
  err = json.Unmarshal(packet, &tmp)
  Error(err)

  /// Debug
  // for exportFile, imageRef := range tmp {
  //   fmt.Printf("%-20s : %s\n", exportFile, imageRef)
  // }
  return tmp
}

/////hogehoge
func recieveImg(s []string) {
  // Socket
  conn, err := net.Dial("tcp", MASTER)
  Error(err)
  defer conn.Close()

  // Request Image
  fmt.Println("Debug: Create Request Packet")
  cLayer := cacis.RequestImage()
  packet := cLayer.Marshal()
  fmt.Println(packet)
  conn.Write(packet)
  fmt.Println("Send\n\n")

  for _, fileName := range s {
    fmt.Printf("Debug: Recieve file '%s'\n", fileName)
    // Recieve Packet HEAD
    packet := make([]byte, cacis.CacisLayerSize)
    packetLength, err := conn.Read(packet)
    Error(err)
    fmt.Printf("Debug: Recieve only CacisLayer HEAD. len: %d", packetLength)
    //fmt.Println(packet)
    cLayer = cacis.Unmarshal(packet)
    //fmt.Println(cLayer)
    //fmt.Println(string(cLayer.Payload))

    // Recieve Packet PAYLOAD
    fmt.Println("Debug: Recieve CacisLayer PAYLOAD")
    //fmt.Println(cLayer.Length)
    packet = []byte{}
    recievedBytes := 0
    targetBytes := cLayer.Length

    for {
      buf := make([]byte, targetBytes - uint64(recievedBytes))
      packetLength, err = conn.Read(buf)
      Error(err)
      recievedBytes += packetLength
      packet = append(packet, buf[:packetLength]...)
      fmt.Println("Debug: recieving...")
      fmt.Printf("Complete %d of %d\n", len(packet), targetBytes)
      if len(packet) == int(targetBytes) {
        break
      }
    }

    fmt.Printf("Debug: Image Packet from Master. len: %d\n", len(packet))
    //fmt.Println("Debug:!!!!!!!!!!!!!!!!!!!!!!!!!!!")
    //cLayer.Payload = packet
    //fmt.Println(cLayer[:10])
    //fmt.Println(string(cLayer.Payload))

    fmt.Printf("Debug: Write file '%s'\n", fileName)
    // File
    filePath := importDir + fileName
    //filePath := "./hoge1.txt"
    file , err := os.Create(filePath)
    Error(err)

    file.Write(packet)
  }
}

func importImg(imageName, filePath string) {
  fmt.Println("Importing " + imageName + " from " + filePath + "...")

  ctx := context.Background()
  client, err := containerd.New("/run/containerd/containerd.sock", containerd.WithDefaultNamespace("cacis"))
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

func install_microk8s() {
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

func Error(err error) {
  if err != nil {
    fmt.Println(err)
  }
}
