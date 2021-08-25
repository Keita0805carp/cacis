package master

import (
  "fmt"
  "log"
  "sort"
  "net"
  "os"
  "context"

  "github.com/keita0805carp/cacis/cacis"

  "github.com/containerd/containerd"
  "github.com/containerd/containerd/platforms"
  "github.com/containerd/containerd/images/archive"
)

const (
  exportDir = "./master-vol/"
  //containerdSock = "/run/containerd/containerd.sock"
  containerdSock = "/var/snap/microk8s/common/run/containerd.sock"
  containerdNameSpace = "cacis"
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
  exportAndPullAllImg()
  server()
}


func server() {
  // Socket
  listen, err := net.Listen("tcp", "localhost:27001")
  cacis.Error(err)
  defer listen.Close()

  for {
    log.Printf("[Debug] Waiting slave\n\n")
    conn, err := listen.Accept()
    cacis.Error(err)
    handling(conn)
    conn.Close()
  }
}

func handling(conn net.Conn) {
  // Recieve Request from slave
  buf := make([]byte, cacis.CacisLayerSize)
  packetLength, err := conn.Read(buf)
  cacis.Error(err)

  log.Printf("[Debug] Recieve Packet from Slave. len: %d\n", packetLength)
  fmt.Println(buf)
  cLayer := cacis.Unmarshal(buf)
  //fmt.Println(string(rl.Payload))

  /// Swtich Type
  if cLayer.Type == 10 {  /// request Components List

    log.Println("[Debug] Type = 10")
    sendComponentsList(conn)

  } else if cLayer.Type == 20 {  /// request Image

    log.Println("[Debug] Type = 20")
    sendImg(conn)

  } else {
    log.Println("[Error] Unknown Type")
  }
}

func pullImg(imageName string) {
  log.Printf("\n[Info]  Pulling   %s ...", imageName)

  ctx := context.Background()
  client, err := containerd.New(containerdSock, containerd.WithDefaultNamespace("cacis"))
  defer client.Close()
  cacis.Error(err)

  opts := []containerd.RemoteOpt{
    containerd.WithAllMetadata(),
  }

  contents, err := client.Fetch(ctx, imageName, opts...)
  cacis.Error(err)

  image := containerd.NewImageWithPlatform(client, contents, platforms.All)
  if image == nil {
    log.Println("[Error] Fail to Pull")
    }
  log.Printf("\r[Info]  Pulled    %s Completely\n", imageName)
}

func exportImg(filePath, imageRef string){
  log.Printf("\r[Info]  Exporting %s to %s ...", imageRef, filePath)

  ctx := context.Background()
  client, err := containerd.New("/run/containerd/containerd.sock", containerd.WithDefaultNamespace(containerdNameSpace))
  defer client.Close()
  cacis.Error(err)

  f, err := os.Create(filePath)
  defer f.Close()
  cacis.Error(err)

  imageStore := client.ImageService()
  opts := []archive.ExportOpt{
    archive.WithImage(imageStore, imageRef),
    archive.WithAllPlatforms(),
  }

  client.Export(ctx, f, opts...)
  cacis.Error(err)
  log.Printf("\r[Info]  Exported  %s to %s Completely\n", imageRef, filePath)
}

func exportAndPullAllImg(){
  log.Printf("[Debug] start: Pull and Export Images\n")
  log.Printf("\n[Debug] Pull %d images for Kubernetes Components\n", len(componentsList))
  for exportFile, imageRef := range componentsList {
    //fmt.Printf("%s : %s\n", exportDir + exportFile, imageRef)
    pullImg(imageRef)
    exportImg(exportDir + exportFile, imageRef)
  }
  log.Printf("\n[Debug] end: Pull and Export Images\n\n")
}

func sendComponentsList(conn net.Conn) {
  /// Send Components List
  log.Printf("\n[Debug] start: Send Components List\n")
  cLayer := cacis.SendComponentsList(componentsList)
  packet := cLayer.Marshal()
  //fmt.Println(packet)
  conn.Write(packet)
  log.Print("\n[Debug] end: Send Components List\n")
}

func sendImg(conn net.Conn) {
  log.Print("\n[Debug] start: Send Components Images\n")
  s := sortKeys(componentsList)

  for _, fileName := range s {
    /// File
    filePath := exportDir + fileName

    log.Printf("\n[Debug] Read: file '%s'\n", fileName)
    //filePath := "./test/hoge1.txt"
    file, err := os.Open(filePath)
    cacis.Error(err)
    fileInfo, err := file.Stat()
    cacis.Error(err)
    fileBuf := make([]byte, fileInfo.Size())
    file.Read(fileBuf)

    /// Send Image
    log.Printf("\r[Debug] Sending Image %s ...", fileName)
    cLayer := cacis.SendImage(fileBuf)
    packet := cLayer.Marshal()
    //fmt.Println(cLayer)
    conn.Write(packet)
    log.Printf("\r[Debug] Send Image %s Completely\n", fileName)
  }
  fmt.Printf("\n[Debug] end: Send Components Images\n")
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
