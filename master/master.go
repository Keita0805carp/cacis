package master

import (
  "fmt"
  "sort"
  //"io"
  "net"
  "os"
  //"strconv"
  "context"

  "github.com/keita0805carp/cacis/cacis"

  "github.com/containerd/containerd"
  "github.com/containerd/containerd/platforms"
  "github.com/containerd/containerd/images/archive"
)

const exportDir = "./master-vol/"

var componentsList = map[string]string {
  "cni.img"             : "docker.io/calico/cni:v3.13.2",
  "pause.img"           : "k8s.gcr.io/pause:3.1",
  "kube-controllers.img": "docker.io/calico/kube-controllers:v3.13.2",
  "pod2daemon.img"      : "docker.io/calico/pod2daemon-flexvol:v3.13.2",
  "node.img"            : "docker.io/calico/node:v3.13.2",
  "coredns.img"         : "docker.io/coredns/coredns:1.8.0",
  "metrics-server.img"  : "k8s.gcr.io/metrics-server-arm64:v0.3.6",
  "dashboard.img"       : "docker.io/kubernetesui/dashboard:v2.0.0",
}

func Main() {
  exportAndPullAllImg()
  server()
}


func server() {
  // Socket
  listen, err := net.Listen("tcp", "localhost:27001")
  Error(err)
  defer listen.Close()

  for {
    conn, err := listen.Accept()
    Error(err)
    fmt.Printf("Debug: Waiting slave\n\n")
    handling(conn)
    conn.Close()
  }
}

func handling(conn net.Conn) {
  // Recieve Request from slave
  buf := make([]byte, cacis.CacisLayerSize)
  packetLength, err := conn.Read(buf)
  Error(err)

  fmt.Printf("Recieve Packet from Slave. len: %d\n", packetLength)
  fmt.Println(buf)
  cLayer := cacis.Unmarshal(buf)
  //fmt.Println(rl)
  //fmt.Println(rl.Payload)
  //fmt.Println(string(rl.Payload))

  /// Swtich Type
  if cLayer.Type == 10 {  /// request Components List

    fmt.Println("Debug: Type = 10")
    sendComponentsList(conn)

  } else if cLayer.Type == 20 {  /// request Image

    fmt.Println("Debug: Type = 20")
    sendImg(conn)

  } else {
    fmt.Println("Err: Unknown Type")
  }
}

func pullImg(imageName string) {
  fmt.Println("Pulling " + imageName + " ...")

  ctx := context.Background()
  client, err := containerd.New("/run/containerd/containerd.sock", containerd.WithDefaultNamespace("cacis"))
  defer client.Close()
  Error(err)

  opts := []containerd.RemoteOpt{
    containerd.WithAllMetadata(),
  }

  contents, err := client.Fetch(ctx, imageName, opts...)
  Error(err)

  image := containerd.NewImageWithPlatform(client, contents, platforms.All)
  if image == nil {
    fmt.Println("Fail to Pull")
    }
  // fmt.Print("Debug: image= ")
  // fmt.Println(image)
  fmt.Println("Pulled")
}

func exportImg(filePath, imageRef string){
  fmt.Printf("Exporting  %s to %s ...\n", imageRef, filePath)

  ctx := context.Background()
  client, err := containerd.New("/run/containerd/containerd.sock", containerd.WithDefaultNamespace("cacis"))
  defer client.Close()
  Error(err)

  f, err := os.Create(filePath)
  defer f.Close()
  Error(err)

  imageStore := client.ImageService()
  opts := []archive.ExportOpt{
    archive.WithImage(imageStore, imageRef),
    archive.WithAllPlatforms(),
  }

  client.Export(ctx, f, opts...)
  Error(err)
  fmt.Println("Exported")
}

func exportAndPullAllImg(){
  fmt.Printf("Debug: [start] Pull and Export Images\n")
  fmt.Printf("Pull %d images for Kubernetes Components", len(componentsList))
  for exportFile, imageRef := range componentsList {
    fmt.Printf("%s : %s\n", exportDir + exportFile, imageRef)
    fmt.Println("start")
    pullImg(imageRef)
    exportImg(exportDir + exportFile, imageRef)
    fmt.Println("end\n")
  }
  fmt.Printf("Debug: [end] Pull and Export Images\n")
}

func sendComponentsList(conn net.Conn) {
  /// Send Components List
  fmt.Println("Debug: [start] Send Components List")
  cLayer := cacis.SendComponentsList(componentsList)
  packet := cLayer.Marshal()
  //fmt.Println(packet)
  conn.Write(packet)
  fmt.Println("Debug: [end] Send Components List")
}

//////hogehoge
func sendImg(conn net.Conn) {
  fmt.Println("Debug: [start] Send Components Images")
  s := sortKeys(componentsList)

  for _, fileName := range s {
    /// File
    filePath := exportDir + fileName

    fmt.Printf("Debug: Read file '%s'\n", fileName)
    //filePath := "./test/hoge1.txt"
    file, err := os.Open(filePath)
    Error(err)
    fileInfo, err := file.Stat()
    Error(err)
    fileBuf := make([]byte, fileInfo.Size())
    file.Read(fileBuf)

    /// Send Image
    fmt.Println("Debug: Send Image")
    cLayer := cacis.SendImage(fileBuf)
    packet := cLayer.Marshal()
    //fmt.Println("Debug: !!!!!!!!!!!!!!!!!!!!!!!!!!!!")
    //fmt.Println(cLayer)
    conn.Write(packet)
    fmt.Println("Debug: Send Image Packet to Slave.")
  }
  fmt.Println("Debug: [end] Send Components Images")
}

func microk8s_enable(){
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
