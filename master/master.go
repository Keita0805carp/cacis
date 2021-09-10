package master

import (
  "os"
  "fmt"
  "log"
  "net"
  "regexp"
  "context"

  "github.com/keita0805carp/cacis/cacis"

  "github.com/containerd/containerd"
  "github.com/containerd/containerd/platforms"
  "github.com/containerd/containerd/images/archive"
)

const (
  masterIP = cacis.MasterIP
  masterPort = cacis.MasterPort
  targetDir = cacis.TargetDir
  containerdSock = cacis.ContainerdSock
  containerdNameSpace = cacis.ContainerdNameSpace
)
var (
  componentsList = cacis.ComponentsList
)

func Main(cancel chan struct{}) {
  downloadMicrok8s()
  installMicrok8s()
  exportAndPullAllImg()
  go server(cancel)
}

func server(cancel chan struct{}) {
  log.Println("[Debug] Starting Main Server")
  // Socket
  listen, err := net.Listen("tcp", masterIP+":"+masterPort)
  cacis.Error(err)
  defer listen.Close()

  for {
    select {
    default:
      log.Printf("[Debug] Waiting slave\n\n")
      conn, err := listen.Accept()
      cacis.Error(err)
      handling(conn)
      conn.Close()
    case <- cancel:
      log.Println("[Debug] Terminating Main server...")
      return
    }
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

  } else if cLayer.Type == 30 {  /// request microk8s snap

    log.Println("[Debug] Type = 30")
    sendMicrok8sSnap(conn)

  } else if cLayer.Type == 40 {  /// request snapd

    log.Println("[Debug] Type = 40")
    sendSnapd(conn)

  } else if cLayer.Type == 50 {  /// request snapd

    log.Println("[Debug] Type = 50")
    clustering(conn)

  } else if cLayer.Type == 60 {  /// request unclustering

    log.Println("[Debug] Type = 60")
    unclustering(conn, cLayer)

  } else {
    log.Println("[Error] Unknown Type")
  }
}

func pullImg(imageName string) {
  log.Printf("\n[Info]  Pulling   %s ...", imageName)

  ctx := context.Background()
  client, err := containerd.New(containerdSock, containerd.WithDefaultNamespace(containerdNameSpace))
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
  client, err := containerd.New(containerdSock, containerd.WithDefaultNamespace(containerdNameSpace))
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
    //fmt.Printf("%s : %s\n", exporttDir + exportFile, imageRef)
    pullImg(imageRef)
    exportImg(targetDir + exportFile, imageRef)
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
  s := cacis.SortKeys(componentsList)

  for _, fileName := range s {
    fileBuf := readFileByte(fileName)
    /// Send Image
    log.Printf("\r[Debug] Sending Image %s ...", fileName)
    cLayer := cacis.SendImage(fileBuf)
    packet := cLayer.Marshal()
    //fmt.Println(cLayer)
    conn.Write(packet)
    log.Printf("\r[Debug] Send Image %s Completely\n", fileName)
  }
  log.Printf("\n[Debug] end: Send Components Images\n")
}

func clustering(conn net.Conn) {
  log.Print("\nDebug: [start] Clustering\n")
  //TODO microk8s enable dns dashboard
  output, err := cacis.ExecCmd("microk8s add-node", true)
  cacis.Error(err)
  //fmt.Println(string(output))
  //TODO regex getc command to join node
  regex := regexp.MustCompile("microk8s join " + masterIP + ".*")
  joinCmd := regex.FindAllStringSubmatch(string(output), 1)[0][0]

  /// Send CLuster Info
  cLayer := cacis.SendClusterInfo([]byte(joinCmd))
  packet := cLayer.Marshal()
  //fmt.Println(cLayer)
  conn.Write(packet)
  log.Printf("\nDebug: [end] Clustering\n")
}

func unclustering(conn net.Conn, cLayer cacis.CacisLayer) {
  log.Printf("\nDebug: [start] Unclustering\n")
  //TODO get request
  buf := make([]byte, cLayer.Length)
  packetLength, err := conn.Read(buf)
  cacis.Error(err)
  log.Printf("[Debug] Read Packet PAYLOAD. len: %d\n", packetLength)
  //fmt.Println(string(buf))

  cacis.ExecCmd("microk8s remove-node " + string(buf), true)
  log.Printf("\nDebug: [end] Unclustering\n")
}

func readFileByte(fileName string) []byte {
  /// File
  filePath := targetDir + fileName

  log.Printf("\n[Debug] Read file '%s'\n", fileName)
  //filePath := "./test/hoge1.txt"
  file, err := os.Open(filePath)
  cacis.Error(err)
  fileInfo, err := file.Stat()
  cacis.Error(err)
  fileBuf := make([]byte, fileInfo.Size())
  file.Read(fileBuf)

  return fileBuf
}

