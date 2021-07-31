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

func Main() {
  //exportAllImg()
  server()
  //sendData()
}


func exportAllImg(){
  images := map[string]string {
    "cni.img": "docker.io/calico/cni:v3.13.2",
    "pause.img": "k8s.gcr.io/pause:3.1",
    "kube-controllers.img": "docker.io/calico/kube-controllers:v3.13.2",
    "pod2daemon.img": "docker.io/calico/pod2daemon-flexvol:v3.13.2",
    "node.img": "docker.io/calico/node:v3.13.2",
    "coredns.img": "docker.io/coredns/coredns:1.8.0",
    "metrics-server.img": "k8s.gcr.io/metrics-server-arm64:v0.3.6",
    "dashboard.img": "docker.io/kubernetesui/dashboard:v2.0.0",
  }

  fmt.Printf("Pull %d images for Kubernetes Components", len(images))
  for file, imageRef := range images {
    fmt.Printf("%s : %s\n", "./output/" + file, imageRef)
    fmt.Println("start")
    pullImg(imageRef)
    exportImg("./output/" + file, imageRef)
    fmt.Println("end\n")
    }
}

func server() {
  // Socket
  listen, err := net.Listen("tcp", "localhost:27001")
  if err != nil {
    fmt.Println(err)
    }
  defer listen.Close()

  conn, err := listen.Accept()
  // Recieve Request from slave
  buf := make([]byte, cacis.CacisLayerSize)
  _, err = conn.Read(buf)
  if err != nil {
    fmt.Println(err)
  }
  fmt.Println("Recieve Request")
  fmt.Println(buf)
  rl := cacis.Unmarshal(buf)
  //fmt.Println(rl)
  //fmt.Println(rl.Payload)
  //fmt.Println(string(rl.Payload))

  // Swtich Type
  if rl.Type == 1 {  // request Image
    fmt.Println("Type = 1")
    // File
    //filePath := "./test/hoge2.txt"
    filePath := "./alpine.img"
    file, err := os.Open(filePath)
    if err != nil {
	    fmt.Println(err)
    }

    fileInfo, err := file.Stat()
    if err != nil {
	    fmt.Println(err)
    }
    fileBuf := make([]byte, fileInfo.Size())
    file.Read(fileBuf)

    // Notify Image Size
    fmt.Println("Notify Image Size")
    nl := cacis.NotifyImageSize(fileBuf)
    msg_n := nl.Marshal()
    fmt.Println(msg_n)
    conn.Write(msg_n)

    // Send Image
    fmt.Println("Send Image")
    sl := cacis.SendImage(fileBuf)
    msg_s := sl.Marshal()
    fmt.Println(msg_s)
    conn.Write(msg_s)

    conn.Close()
  } else {
    fmt.Println("Unknown Type")
    conn.Close()
  }


  /*

  // buffer
  buf := make([]byte, fileInfo.Size())
  if err != nil {
    fmt.Println(err)
    }

  // Wait
  conn, err := l.Accept()
  if err != nil {
    fmt.Println(err)
    }

  //fmt.Println(buf)
  file.Read(buf)
  fmt.Println(buf)
  //conn.Write(buf)
  fmt.Println(string(buf))
  */

}

func sendData() {
  filePath := "./test/hoge1.txt"
  file, err := os.Open(filePath)
  if err != nil {
    fmt.Println(err)
    }

  fileInfo, err := file.Stat()
  if err != nil {
    fmt.Println(err)
    }

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
  if err != nil {
    fmt.Println(err)
    }

  f, err := os.Create(fileName)
  defer f.Close()
  if err != nil {
    fmt.Println(err)
    }

  imageStore := client.ImageService()
  opts := []archive.ExportOpt{
    archive.WithImage(imageStore, imageName),
    archive.WithAllPlatforms(),
  }

  client.Export(ctx, f, opts...)
  if err != nil {
    fmt.Println(err)
    }
  fmt.Println("Exported")
}

func pullImg(imageName string) {
  fmt.Println("Pulling " + imageName + " ...")

  ctx := context.Background()
  client, err := containerd.New("/run/containerd/containerd.sock", containerd.WithDefaultNamespace("cacis"))
  defer client.Close()
  if err != nil {
    fmt.Println(err)
    }

  opts := []containerd.RemoteOpt{
    containerd.WithAllMetadata(),
  }

  contents, err := client.Fetch(ctx, imageName, opts...)
  if err != nil {
    fmt.Println(err)
    }

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
