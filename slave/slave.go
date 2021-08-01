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
  importDir = "./slave-vol/"
)

func Main() {
  //importAllImg()
  slave()
}

func slave() {
  // Socket
  conn, err := net.Dial("tcp", MASTER)
  Error(err)
  defer conn.Close()

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
  // Socket
  componentsList := recieveComponentsList(conn)

  conn, err = net.Dial("tcp", MASTER)
  Error(err)
  defer conn.Close()

  //fmt.Println(componentsList)
  sortedExportFileName := sortKeys(componentsList)
  recieveImages(conn, sortedExportFileName)
  // Request Image
  //fmt.Println(sortedExportFileName)
}

func recieveComponentsList(conn net.Conn) map[string]string {
  // Request Image
  fmt.Println("Debug: Create Request Packet")
  cLayer := cacis.RequestComponentsList()
  packet := cLayer.Marshal()
  fmt.Println(packet)
  conn.Write(packet)
  fmt.Println("Send\n\n")


  // Recieve Components List Size Notification
  fmt.Println("Debug: Components List Size Notification")
  packet = make([]byte, cacis.CacisLayerSize)
  packetLength, err := conn.Read(packet)
  Error(err)
  fmt.Printf("Debug: Notifying Packet from Master. len: %d\n", packetLength)
  cLayer = cacis.Unmarshal(packet)
  //fmt.Println(cLayer)
  //fmt.Println(string(cLayer.Payload))

  // Recieve Components List
  fmt.Println("Debug: Recieve Components List")
  fmt.Println(cLayer)
  packet = make([]byte, cacis.CacisLayerSize + cLayer.Length)
  packetLength, err = conn.Read(packet)
  Error(err)
  fmt.Printf("Debug: Image Packet from Master. len: %d\n", packetLength)
  cLayer = cacis.Unmarshal(packet)
  //fmt.Println(cLayer)

  var tmp map[string]string
  err = json.Unmarshal((cLayer.Payload), &tmp)
  Error(err)

  /// Debug
  // for exportFile, imageRef := range tmp {
  //   fmt.Printf("%-20s : %s\n", exportFile, imageRef)
  // }
  return tmp
}

/////hogehoge
func recieveImages(conn net.Conn, s []string) {
  // Request Image
  fmt.Println("Debug: Create Request Packet")
  cLayer := cacis.RequestImage()
  packet := cLayer.Marshal()
  fmt.Println(packet)
  conn.Write(packet)
  fmt.Println("Send\n\n")

  for _, fileName := range s {
    fmt.Printf("Debug: Recieve file '%s'\n", fileName)
    // Recieve Image Size Notification
    fmt.Println("Debug: Image Size Notification")
    packet := make([]byte, cacis.CacisLayerSize)
    packetLength, err := conn.Read(packet)
    Error(err)
    fmt.Printf("Debug: Notifying Packet from Master. len: %d\n", packetLength)
    fmt.Println(packet)
    cLayer = cacis.Unmarshal(packet)
    //fmt.Println(cLayer)
    //fmt.Println(string(cLayer.Payload))

    // Recieve Image
    fmt.Println("Debug: Recieve Image")
    fmt.Println(cLayer.Length)
    fmt.Printf("Debug: makeslice len: %d\n", cacis.CacisLayerSize + cLayer.Length)
    packet = make([]byte, cacis.CacisLayerSize + cLayer.Length)
    recievedBytes := uint64(0)
    targetBytes := cacis.CacisLayerSize + cLayer.Length

    for {
      buf := make([]byte, targetBytes - recievedBytes)
      packetLength, err = conn.Read(buf)
      Error(err)
      recievedBytes += uint64(packetLength)
      _ = append(packet, buf...)
      fmt.Println("Debug: recieving...")
      fmt.Printf("Complete %d of %d\n", len(buf), len(packet))
      if len(buf) == 0 {
        break
      }
    }

    fmt.Printf("Debug: Image Packet from Master. len: %d\n", len(packet))
    cLayer = cacis.Unmarshal(packet)
    //fmt.Println(cLayer)
    //fmt.Println(string(cLayer.Payload))

    fmt.Printf("Debug: Write file '%s'\n", fileName)
    // File
    filePath := importDir + fileName
    //filePath := "./hoge1.txt"
    file , err := os.Create(filePath)
    Error(err)

    file.Write(cLayer.Payload)
  }
}

func importImg(imageName, fileName string) {
  fmt.Println("Importing " + imageName + " from " + fileName + "...")

  ctx := context.Background()
  client, err := containerd.New("/run/containerd/containerd.sock", containerd.WithDefaultNamespace("cacis"))
  defer client.Close()
  Error(err)

  f, err := os.Open(fileName)
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

func importAllImg(componentsList map[string]string) {
  fmt.Printf("Import %d images for Kubernetes Components", len(componentsList))
  for file, imageRef := range componentsList {
    fmt.Printf("%s : %s\n", importDir + file, imageRef)
    fmt.Println("start")
    importImg(imageRef, importDir + file)
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
