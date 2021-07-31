package master

import (
  "fmt"
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
  //exportAllImg()
  server()
  //sendData()
}


func exportAllImg(){
  fmt.Printf("Pull %d images for Kubernetes Components", len(componentsList))
  for exportFile, imageRef := range componentsList {
    fmt.Printf("%s : %s\n", "./output/" + exportFile, imageRef)
    fmt.Println("start")
    pullImg(imageRef)
    exportImg("./output/" + exportFile, imageRef)
    fmt.Println("end\n")
  }
}

func server() {
  // Socket
  listen, err := net.Listen("tcp", "localhost:27001")
  Error(err)
  defer listen.Close()

  conn, err := listen.Accept()
  handling(conn)
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
  if cLayer.Type == 1 {  /// request Components List
    fmt.Println("Debug: Type = 1")
    /// Notify Image Size
    fmt.Println("Debug: Notify Components List Size")
    cLayer = cacis.NotifyComponentsListSize(componentsList)
    packet := cLayer.Marshal()
    //fmt.Println(packet)
    conn.Write(packet)
    fmt.Println("Debug: Send Notify Packet to Slave.")

    /// Send Image
    fmt.Println("Debug: Send Components List")
    cLayer = cacis.SendComponentsList(componentsList)
    packet = cLayer.Marshal()
    //fmt.Println(packet)
    conn.Write(packet)
    fmt.Println("Debug: Send Components List Packet to Slave.")
  } else if cLayer.Type == 2 {  // request Image
    fmt.Println("Debug: Type = 2")

    /// File
    filePath := "./test/hoge1.txt"
    //filePath := "./alpine.img"
    file, err := os.Open(filePath)
    Error(err)
    fileInfo, err := file.Stat()
    Error(err)
    fileBuf := make([]byte, fileInfo.Size())
    file.Read(fileBuf)

    /// Notify Image Size
    fmt.Println("Debug: Notify Image Size")
    cLayer = cacis.NotifyImageSize(fileBuf)
    packet := cLayer.Marshal()
    //fmt.Println(packet)
    conn.Write(packet)
    fmt.Println("Debug: Send Notify Packet to Slave.")

    /// Send Image
    fmt.Println("Debug: Send Image")
    cLayer = cacis.SendImage(fileBuf)
    packet = cLayer.Marshal()
    //fmt.Println(packet)
    conn.Write(packet)
    fmt.Println("Debug: Send Image Packet to Slave.")

  } else {
    fmt.Println("Err: Unknown Type")
  }

  conn.Close()

  /*
  // Wait
  conn, err := l.Accept()
  if err != nil {
    fmt.Println(err)
    }
  */
}

func sendData() {
  filePath := "./test/hoge1.txt"
  file, err := os.Open(filePath)
  Error(err)

  fileInfo, err := file.Stat()
  Error(err)

  buf := make([]byte, fileInfo.Size())
  fmt.Println(file.Read(buf))
  fmt.Println(buf)
  fmt.Println(string(buf))
  fmt.Println(fileInfo)
  fmt.Println(fileInfo.Name())
  fmt.Println(fileInfo.Size())
}

func exportImg(fileName, imageName string){
  fmt.Println("Exporting " + imageName + " to " + fileName + "...")

  ctx := context.Background()
  client, err := containerd.New("/run/containerd/containerd.sock", containerd.WithDefaultNamespace("cacis"))
  defer client.Close()
  Error(err)

  f, err := os.Create(fileName)
  defer f.Close()
  Error(err)

  imageStore := client.ImageService()
  opts := []archive.ExportOpt{
    archive.WithImage(imageStore, imageName),
    archive.WithAllPlatforms(),
  }

  client.Export(ctx, f, opts...)
  Error(err)
  fmt.Println("Exported")
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


func microk8s_enable() {
}

func notify() {
}

func Error(err error) {
  if err != nil {
    fmt.Println(err)
  }
}
