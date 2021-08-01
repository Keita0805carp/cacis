package slave

import (
  "fmt"
  //"io"
  "net"
  "os"
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
  client()
  //recieveImages(images)
  //recieveComponentsList()
}

func client() {
  // Socket
  conn, err := net.Dial("tcp", MASTER)
  Error(err)
  defer conn.Close()

  componentsList := recieveComponentsList(conn)
  recieveImages(conn, componentsList)
  // Request Image
  //fmt.Println(componentsList)
}

func recieveComponentsList(conn net.Conn) map[string]string {
  // Request Image
  fmt.Println("Debug: Create Request Packet")
  cLayer := cacis.RequestComponentsList()
  packet := cLayer.Marshal()
  fmt.Println(packet)
  conn.Write(packet)
  fmt.Println("Send\n\n")


  // Recieve Image Size Notification
  fmt.Println("Debug: Components List Size Notification")
  buf := make([]byte, cacis.CacisLayerSize)
  packetLength, err := conn.Read(buf)
  Error(err)
  fmt.Printf("Debug: Notifying Packet from Master. len: %d\n", packetLength)
  cLayer = cacis.Unmarshal(buf)
  //fmt.Println(cLayer)
  //fmt.Println(string(cLayer.Payload))

  // Recieve Image
  fmt.Println("Debug: Recieve Components List")
  fmt.Println(cLayer)
  buf = make([]byte, cacis.CacisLayerSize + cLayer.Length)
  packetLength, err = conn.Read(buf)
  Error(err)
  fmt.Printf("Debug: Image Packet from Master. len: %d\n", packetLength)
  cLayer = cacis.Unmarshal(buf)
  //fmt.Println(cLayer)

  var tmp map[string]string
  err = json.Unmarshal((cLayer.Payload), &tmp)
  Error(err)

  for exportFile, imageRef := range tmp {
    fmt.Printf("%-20s : %s\n", exportFile, imageRef)
  }
  return tmp
}

func recieveImages(conn net.Conn, componentsList map[string]string) {

  // Request Image
  fmt.Println("Debug: Create Request Packet")
  cLayer := cacis.RequestImage()
  packet := cLayer.Marshal()
  fmt.Println(packet)
  conn.Write(packet)
  fmt.Println("Send\n\n")


  for fileName, imageRef := range componentsList {
    // Recieve Image Size Notification
    fmt.Println("Debug: Image Size Notification")
    buf := make([]byte, cacis.CacisLayerSize)
    packetLength, err := conn.Read(buf)
    Error(err)
    fmt.Printf("Debug: Notifying Packet from Master. len: %d\n", packetLength)
    cLayer = cacis.Unmarshal(buf)
    //fmt.Println(cLayer)
    //fmt.Println(string(cLayer.Payload))

    // Recieve Image
    fmt.Println("Debug: Recieve Image")
    fmt.Println(cLayer)
    buf = make([]byte, cacis.CacisLayerSize + cLayer.Length)
    packetLength, err = conn.Read(buf)
    Error(err)
    fmt.Printf("Debug: Image Packet from Master. len: %d\n", packetLength)
    cLayer = cacis.Unmarshal(buf)
    //fmt.Println(cLayer)
    //fmt.Println(string(cLayer.Payload))


    fmt.Printf("Debug: Write file '%s'(%s)\n", fileName, imageRef)
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
  // images := map[string]string {
  //   "cni.img": "docker.io/calico/cni:v3.13.2",
  //   "pause.img": "docker.io/calico/kube-controllers:v3.13.2",
  //   "kube-controllers.img": "docker.io/calico/pod2daemon-flexvol:v3.13.2",
  //   "pod2daemon.img": "docker.io/calico/node:v3.13.2",
  //   "node.img": "docker.io/calico/node:v3.13.2",
  //   "coredns.img": "docker.io/coredns/coredns:1.8.0",
  //   "metrics-server.img": "k8s.gcr.io/metrics-server-arm64:v0.3.6",
  //   "dashboard.img": "docker.io/kubernetesui/dashboard:v2.0.0",
  // }

  fmt.Printf("Import %d images for Kubernetes Components", len(componentsList))
  for file, imageRef := range componentsList {
    fmt.Printf("%s : %s\n", "./output/" + file, imageRef)
    fmt.Println("start")
    importImg(imageRef, "./output/" + file)
    fmt.Println("end\n")
    }
}

func install_microk8s() {
}

func Error(err error) {
  if err != nil {
    fmt.Println(err)
  }
}
