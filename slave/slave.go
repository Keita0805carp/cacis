package slave

import (
  "fmt"
  //"io"
  "net"
  "os"
  //"strconv"
  //"strings"
  "context"

  "github.com/keita0805carp/cacis/cacis"

  "github.com/containerd/containerd"
)

func Main() {
  //importAllImg()
  //client()
  requestData()
}

func client() {
  // Socket
  conn, err := net.Dial("tcp", "localhost:27001")
  Error(err)
  defer conn.Close()

  // File
  filePath := "./hoge1.txt"
  file , err := os.Create(filePath)
  Error(err)

  buf := make([]byte, 22)
  conn.Read(buf)
  fmt.Println(buf)
  fmt.Println(string(buf))

  file.Write(buf)
}

func requestData() {
  conn, err := net.Dial("tcp", "localhost:27001")
  Error(err)
  defer conn.Close()

  // Request Image
  fmt.Println("Create Request Packet")
  sl := cacis.RequestImage()
  msg_s := sl.Marshal()
  fmt.Println(msg_s)
  conn.Write(msg_s)
  fmt.Println("Send\n\n")


  // Recieve Image Size Notification
  fmt.Println("Image Size Notification")
  buf := make([]byte, cacis.CacisLayerSize)
  _, err = conn.Read(buf) //TODO
  Error(err)
  rl := cacis.Unmarshal(buf)
  fmt.Println(buf)
  fmt.Println(string(buf))

  // Recieve Image
  fmt.Println("Recieve Image")
  fmt.Println(rl)
  buf = make([]byte, cacis.CacisLayerSize + rl.Length)
  _, err = conn.Read(buf)
  Error(err)
  rl = cacis.Unmarshal(buf)
  fmt.Println(buf)
  fmt.Println(string(buf))


  // File
  //filePath := "./hoge2.txt"
  filePath := "./test/alpine.img"
  file , err := os.Create(filePath)
  Error(err)

  file.Write(rl.Payload)
}

func recieveData() {
}


func importAllImg() {
  images := map[string]string {
    "cni.img": "docker.io/calico/cni:v3.13.2",
    "pause.img": "docker.io/calico/kube-controllers:v3.13.2",
    "kube-controllers.img": "docker.io/calico/pod2daemon-flexvol:v3.13.2",
    "pod2daemon.img": "docker.io/calico/node:v3.13.2",
    "node.img": "docker.io/calico/node:v3.13.2",
    "coredns.img": "docker.io/coredns/coredns:1.8.0",
    "metrics-server.img": "k8s.gcr.io/metrics-server-arm64:v0.3.6",
    "dashboard.img": "docker.io/kubernetesui/dashboard:v2.0.0",
  }

  fmt.Printf("Import %d images for Kubernetes Components", len(images))
  for file, imageRef := range images {
    fmt.Printf("%s : %s\n", "./output/" + file, imageRef)
    fmt.Println("start")
    importImg(imageRef, "./output/" + file)
    fmt.Println("end\n")
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

func install_microk8s() {
}

func notify() {
}

func Error(error error) {
  if error != nil {
    fmt.Println(error)
  }
