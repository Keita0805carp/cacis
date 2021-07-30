package slave

import (
  "fmt"
  //"io"
  "net"
  "os"
  //"strconv"
  //"strings"
  "context"

  "github.com/containerd/containerd"
)

func Main() {
  //importAllImg()
  client()
}

func client() {
  // Socket
  conn, err := net.Dial("tcp", "localhost:27001")
  if err != nil {
    fmt.Println(err)
    }
  defer conn.Close()

  // File
  filePath := "./hoge1.txt"
  file , err := os.Create(filePath)
  if err != nil {
    fmt.Println(err)
    }

  buf := make([]byte, 22)
  conn.Read(buf)
  fmt.Println(buf)
  fmt.Println(string(buf))

  file.Write(buf)
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
  if err != nil {
    fmt.Println(err)
    }

  f, err := os.Open(fileName)
  defer f.Close()
  if err != nil {
    fmt.Println(err)
    }

  opts := []containerd.ImportOpt{
    containerd.WithIndexName(imageName),
    //containerd.WithAllPlatforms(true),
  }

  client.Import(ctx, f, opts...)
  if err != nil {
    fmt.Println(err)
    }
  fmt.Println("Imported")
}

func install_microk8s() {
}

func notify() {
}
