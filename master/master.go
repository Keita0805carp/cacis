package master

import (
  "log"
  "net"
  "os"
  "os/signal"
  "time"
  "strings"
  "context"
  "syscall"

  "github.com/keita0805carp/cacis/cacis"
  "github.com/keita0805carp/cacis/connection"

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

func Setup() {
  cacis.CreateTempDir()
  downloadMicrok8s()
  installMicrok8s()
  enableMicrok8s()
  exportAndPullAllImg(false)
  //fmt.Println(getImgList())
}

func Main() {
  terminate := make(chan os.Signal, 1)
  signal.Notify(terminate, syscall.SIGTERM, syscall.SIGINT, os.Interrupt)
  cancel := make(chan struct{})

  cacis.ExecCmd("microk8s start", true)
  connection.Main(cancel)
  go listening(cancel)

  go removeNotReadyNode()
  <-terminate
  close(cancel)
  log.Printf("[Debug]: Terminating Main Master Process...\n")
  time.Sleep(10 * time.Second)
  log.Printf("[Debug]: Terminated Main Master Process\n")
}

func listening(cancel chan struct{}) {
  log.Printf("[Debug] Start Main Server\n")
  // Socket
  listen, err := net.Listen("tcp", masterIP+":"+masterPort)
  cacis.Error(err)
  defer listen.Close()

  for {
    select {
    default:
      log.Printf("[Debug] Listen Request from Slave\n")
      conn2master, err := listen.Accept()
      cacis.Error(err)

      go handling(conn2master)
    case <- cancel:
      log.Println("[Debug] Terminating Main server...")
      //cacis.ExecCmd("microk8s stop", false)
      return
    }
  }
}

func handling(conn2master net.Conn) {
  // Recieve Request from slave
  buf := make([]byte, cacis.CacisLayerSize)
  packetLength, err := conn2master.Read(buf)
  cacis.Error(err)
  log.Printf("[Debug] Recieve Packet from Slave. len: %d\n", packetLength)
  cLayer := cacis.Unmarshal(buf)
  //fmt.Println(buf)
  //fmt.Println(string(rl.Payload))

  remoteIP := conn2master.RemoteAddr().String()[:strings.LastIndex(conn2master.LocalAddr().String(), ":")]

  /// Swtich Type
  if cLayer.Type == 10 {  /// request Components List

    log.Printf("[Debug] Dialing...\n")
    conn2slave, err := net.Dial("tcp", remoteIP+":27001")
    cacis.Error(err)

    conn2master.Close()
    log.Printf("[Debug] Type = 10\n")
    sendComponentsList(conn2slave)

    conn2slave.Close()

  } else if cLayer.Type == 20 {  /// request Image

    log.Printf("[Debug] Dialing...\n")
    conn2slave, err := net.Dial("tcp", remoteIP+":27001")
    cacis.Error(err)

    conn2master.Close()
    log.Print("[Debug] Type = 20\n")
    sendImg(conn2slave)

    conn2slave.Close()

  } else if cLayer.Type == 30 {  /// request microk8s snap

    log.Print("[Debug] Dialing...\n")
    conn2slave, err := net.Dial("tcp", remoteIP+":27001")
    cacis.Error(err)

    conn2master.Close()
    log.Printf("[Debug] Type = 30\n")
    sendMicrok8s(conn2slave)

    conn2slave.Close()

  } else if cLayer.Type == 40 {  /// request snapd

    log.Printf("[Debug] Dialing...\n")
    conn2slave, err := net.Dial("tcp", remoteIP+":27001")
    cacis.Error(err)

    conn2master.Close()
    log.Printf("[Debug] Type = 40\n")
    sendSnapd(conn2slave)

    conn2slave.Close()

  } else if cLayer.Type == 50 {  /// request clustering

    log.Printf("[Debug] Dialing...\n")
    conn2slave, err := net.Dial("tcp", remoteIP+":27001")
    cacis.Error(err)

    conn2master.Close()
    log.Printf("[Debug] Type = 50\n")
    clustering(conn2slave)

    conn2slave.Close()

  } else if cLayer.Type == 60 {  /// request unclustering

    log.Printf("[Debug] Type = 60\n")
    unclustering(conn2master, cLayer)
    conn2master.Close()

  } else {
    conn2master.Close()
    log.Printf("[Error] Unknown Type\n")
  }
}

func getImgList() []string {
  log.Printf("[Info] Show images list\n")

  ctx, client := ContainerdInit()
  defer client.Close()

  images, err := client.ListImages(ctx)
  cacis.Error(err)
  imagesName := make([]string, len(images))
  for i, image := range images {
    imagesName[i] = image.Name()
  }
  return imagesName
}

func pullImg(imageName string) {
  log.Printf("[Info]  Pulling   %s ...\n", imageName)

  ctx, client := ContainerdInit()
  defer client.Close()

  opts := []containerd.RemoteOpt{
    containerd.WithAllMetadata(),
  }

  contents, err := client.Fetch(ctx, imageName, opts...)
  cacis.Error(err)

  image := containerd.NewImageWithPlatform(client, contents, platforms.All)
  if image == nil {
    log.Printf("[Error] Fail to Pull\n")
    }
  log.Printf("\r[Info]  Pulled    %s Completely\n", imageName)
}

func exportImg(filePath, imageRef string){
  log.Printf("\r[Info]  Exporting %s to %s ...", imageRef, filePath)

  ctx, client := ContainerdInit()
  defer client.Close()

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

func exportAndPullAllImg(onlyExport bool){
  log.Printf("[Debug] start: Pull and Export Images\n")
  log.Printf("[Debug] Pull %d images for Kubernetes Components\n", len(componentsList))
  for exportFile, imageRef := range componentsList {
    //fmt.Printf("%s : %s\n", exporttDir + exportFile, imageRef)
    if onlyExport {
      exportImg(targetDir + exportFile, imageRef)
    } else {
      pullImg(imageRef)
      exportImg(targetDir + exportFile, imageRef)
    }
  }
  log.Printf("[Debug] end: Pull and Export Images\n")
}

func sendComponentsList(conn net.Conn) {
  /// Send Components List
  log.Printf("[Debug] start: Send Components List\n")
  cLayer := cacis.SendComponentsList(componentsList)
  packet := cLayer.Marshal()
  //fmt.Println(packet)
  conn.Write(packet)
  log.Print("[Debug] end: Send Components List\n")
}

func sendImg(conn net.Conn) {
  log.Print("[Debug] start: Send Components Images\n")
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

func readFileByte(fileName string) []byte {
  /// File
  filePath := targetDir + fileName

  log.Printf("[Debug] Read file '%s'\n", fileName)
  //filePath := "./test/hoge1.txt"
  file, err := os.Open(filePath)
  cacis.Error(err)
  fileInfo, err := file.Stat()
  cacis.Error(err)
  fileBuf := make([]byte, fileInfo.Size())
  file.Read(fileBuf)

  return fileBuf
}

func ContainerdInit() (context.Context, *containerd.Client) {
  ctx := context.Background()
  client, err := containerd.New(containerdSock, containerd.WithDefaultNamespace(containerdNameSpace))
  cacis.Error(err)
  return ctx, client
}

