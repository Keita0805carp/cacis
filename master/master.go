package master

import (
  "fmt"
  "sort"
  "net"
  "time"
  "os"
  "os/exec"
  "regexp"
  "context"
  "strings"

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

func Main() {
  //downloadMicrok8s()
  //installMicrok8s()
  //fmt.Println("wait 5 seconds")
  time.Sleep(0 * time.Second)
  //exportAndPullAllImg()
  server()
}


func server() {
  // Socket
  listen, err := net.Listen("tcp", masterIP+":"+masterPort)
  Error(err)
  defer listen.Close()

  for {
    fmt.Printf("Debug: Waiting slave\n\n")
    conn, err := listen.Accept()
    Error(err)
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
  //fmt.Println(string(rl.Payload))

  /// Swtich Type
  if cLayer.Type == 10 {  /// request Components List

    fmt.Println("Debug: Type = 10")
    sendComponentsList(conn)

  } else if cLayer.Type == 20 {  /// request Image

    fmt.Println("Debug: Type = 20")
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

  } else {
    fmt.Println("Err: Unknown Type")
  }
}

func pullImg(imageName string) {
  fmt.Printf("\nPulling   %s ...", imageName)

  ctx := context.Background()
  client, err := containerd.New(containerdSock, containerd.WithDefaultNamespace("cacis"))
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
  fmt.Printf("\rPulled    %s Completely\n", imageName)
}

func exportImg(filePath, imageRef string){
  fmt.Printf("\rExporting %s to %s ...", imageRef, filePath)

  ctx := context.Background()
  client, err := containerd.New(containerdSock, containerd.WithDefaultNamespace(containerdNameSpace))
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
  fmt.Printf("\rExported  %s to %s Completely\n", imageRef, filePath)
}

func exportAndPullAllImg(){
  fmt.Printf("Debug: [start] Pull and Export Images\n")
  fmt.Printf("\nPull %d images for Kubernetes Components\n", len(componentsList))
  for exportFile, imageRef := range componentsList {
    //fmt.Printf("%s : %s\n", exportDir + exportFile, imageRef)
    pullImg(imageRef)
    exportImg(exportDir + exportFile, imageRef)
  }
  fmt.Printf("\nDebug: [end] Pull and Export Images\n\n")
}

func sendComponentsList(conn net.Conn) {
  /// Send Components List
  fmt.Printf("\nDebug: [start] Send Components List\n")
  cLayer := cacis.SendComponentsList(componentsList)
  packet := cLayer.Marshal()
  //fmt.Println(packet)
  conn.Write(packet)
  fmt.Print("\nDebug: [end] Send Components List\n")
}

func sendImg(conn net.Conn) {
  fmt.Print("\nDebug: [start] Send Components Images\n")
  s := sortKeys(componentsList)

  for _, fileName := range s {
    fileBuf := readFileByte(fileName)

    /// Send Image
    fmt.Printf("\rDebug: Sending Image %s ...", fileName)
    cLayer := cacis.SendImage(fileBuf)
    packet := cLayer.Marshal()
    //fmt.Println(cLayer)
    conn.Write(packet)
    fmt.Printf("\rDebug: Send Image %s Completely\n", fileName)
  }
  fmt.Printf("\nDebug: [end] Send Components Images\n")
}

func snapd() {
  //TODO snapd check
  //TODO snapd install
  myexec("apt install snapd")
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
  myexec("snap download microk8s --target-directory=" + exportDir)
  fmt.Printf("Download Completely\n")
}

func installMicrok8s() {
  //TODO microk8s check
  //TODO snap install microk8s --classic
  fmt.Printf("Install microk8s via snap\n")
  fmt.Printf("Installing...")
  myexec("snap ack " + exportDir + "microk8s_2346.assert")
  myexec("snap install " + exportDir + "microk8s_2346.snap" + " --classic")
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
  //myexec("echo \n"cgroup_enable=memory cgroup_memory=1"\n >> /boot/firmware/cmdline.txt")
}

func clustering(conn net.Conn) {
  fmt.Print("\nDebug: [start] Clustering\n")
  //TODO microk8s enable dns dashboard
  output, err := myexec("microk8s add-node")
  Error(err)
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
  myexec("microk8s enable dns dashboard")
}

func getKubeconfig() {
  myexec("microk8s config")
}

func unclustering() {
  //TODO get request
  //TODO get hostname wants to leave
  myexec("microk8s remove-node cacis-vagrant-slave")
}


func readFileByte(fileName string) []byte {
  /// File
  filePath := exportDir + fileName

  fmt.Printf("\nDebug: Read file '%s'\n", fileName)
  //filePath := "./test/hoge1.txt"
  file, err := os.Open(filePath)
  Error(err)
  fileInfo, err := file.Stat()
  Error(err)
  fileBuf := make([]byte, fileInfo.Size())
  file.Read(fileBuf)

  return fileBuf
}

func myexec(cmd string) ([]byte, error) {
  slice := strings.Split(cmd, " ")
  stdout, err := exec.Command(slice[0], slice[1:]...).Output()
  //fmt.Printf("exec: %s\noutput:\n%s", cmd, stdout)
  Error(err)
  return stdout, err
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
