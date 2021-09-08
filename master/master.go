package master

import (
  "fmt"
  "log"
  "sort"
  "net"
  "os"
  "regexp"
  "context"

  "github.com/keita0805carp/cacis/cacis"

  "github.com/containerd/containerd"
  "github.com/containerd/containerd/platforms"
  "github.com/containerd/containerd/images/archive"
)

const (
  masterIP = "192.168.56.250"
  masterPort = "27001"
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

func Main(cancel chan struct{}) {
  //downloadMicrok8s()
  //installMicrok8s()
  exportAndPullAllImg()
  server(cancel)
}


func server(cancel chan struct{}) {
  // Socket
  listen, err := net.Listen("tcp", masterIP+":"+masterPort)
  cacis.Error(err)
  defer listen.Close()

  for (cancel == nil) {
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

  } else if cLayer.Type == 30 {  /// request microk8s snap

    fmt.Println("Debug: Type = 30")
    sendMicrok8sSnap(conn)

  } else if cLayer.Type == 40 {  /// request snapd

    fmt.Println("Debug: Type = 40")
    sendSnapd(conn)

  } else if cLayer.Type == 50 {  /// request snapd

    fmt.Println("Debug: Type = 50")
    clustering(conn)

  } else if cLayer.Type == 60 {  /// request unclustering

    fmt.Println("Debug: Type = 60")
    unclustering(conn, cLayer)

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
    fileBuf := readFileByte(fileName)

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

func snapd() {
  //TODO snapd check
  //TODO snapd install
  cacis.ExecCmd("apt install snapd", false)
}

func sendSnapd(conn net.Conn) {
  fmt.Print("\nDebug: [start] Send Snapd\n")
  s := []string{"snapd.zip"}

  for _, fileName := range s {
    fileBuf := readFileByte(fileName)

    /// Send Image
    fmt.Printf("\rDebug: Sending Snapd %s ...", fileName)
    cLayer := cacis.SendSnapd(fileBuf)
    packet := cLayer.Marshal()
    //fmt.Println(cLayer)
    conn.Write(packet)
    fmt.Printf("\rDebug: Send Snapd %s Completely\n", fileName)
  }
  fmt.Printf("\nDebug: [end] Send Snapd\n")
}

func downloadMicrok8s() {
  //TODO microk8s check
  //TODO snap install microk8s --classic
  fmt.Printf("Download microk8s via snap\n")
  fmt.Printf("Downloading...")
  cacis.ExecCmd("snap download microk8s --target-directory=" + exportDir, false)
  fmt.Printf("Download Completely\n")
}

func installMicrok8s() {
  //TODO microk8s check
  //TODO snap install microk8s --classic
  fmt.Printf("Install microk8s via snap\n")
  fmt.Printf("Installing...")
  cacis.ExecCmd("snap ack " + exportDir + "microk8s_2346.assert", false)
  cacis.ExecCmd("snap install " + exportDir + "microk8s_2346.snap" + " --classic", true)
  fmt.Printf("Install Completely\n")
}

func sendMicrok8sSnap(conn net.Conn) {
  fmt.Print("\nDebug: [start] Send Snap files\n")
  s := []string{"microk8s_2346.assert", "microk8s_2346.snap", "core_11420.assert", "core_11420.snap"}

  for _, fileName := range s {
    fileBuf := readFileByte(fileName)

    /// Send Image
    fmt.Printf("\rDebug: Sending Snap files %s ...", fileName)
    cLayer := cacis.SendMicrok8sSnap(fileBuf)
    packet := cLayer.Marshal()
    //fmt.Println(cLayer)
    conn.Write(packet)
    fmt.Printf("\rDebug: Send Snap files %s Completely\n", fileName)
  }
  fmt.Printf("\nDebug: [end] Send Snap files\n")
}

func setupMicrok8s() {
  /// if RaspberryPi
  //cacis.ExecCmd("echo \n"cgroup_enable=memory cgroup_memory=1"\n >> /boot/firmware/cmdline.txt")
}

func clustering(conn net.Conn) {
  fmt.Print("\nDebug: [start] Clustering\n")
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
  fmt.Printf("\nDebug: [end] Clustering\n")
}

func enableMicrok8s() {
  //TODO microk8s enable dns dashboard
  cacis.ExecCmd("microk8s enable dns dashboard", false)
}

func getKubeconfig() {
  cacis.ExecCmd("microk8s config", true)
}

func unclustering(conn net.Conn, cLayer cacis.CacisLayer) {
  //TODO get request
  buf := make([]byte, cLayer.Length)
  packetLength, err := conn.Read(buf)
  cacis.Error(err)
  fmt.Printf("Debug: Read Packet PAYLOAD. len: %d\n", packetLength)
  fmt.Println(buf)
  fmt.Println(string(buf))

  //TODO get hostname wants to leave
  cacis.ExecCmd("microk8s remove-node cacis-vagrant-slave", false)
}


func readFileByte(fileName string) []byte {
  /// File
  filePath := exportDir + fileName

  fmt.Printf("\nDebug: Read file '%s'\n", fileName)
  //filePath := "./test/hoge1.txt"
  file, err := os.Open(filePath)
  cacis.Error(err)
  fileInfo, err := file.Stat()
  cacis.Error(err)
  fileBuf := make([]byte, fileInfo.Size())
  file.Read(fileBuf)

  return fileBuf
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
