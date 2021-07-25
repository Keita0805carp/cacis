package master

import (
  "fmt"
  "io"
  "net"
  "os"
  "strconv"
)

const BUFFERSIZE = 1024
const MASTER = "10.0.100.1:27001"

func Main() {
  //exportImg()
  server()
}

func transforImg(){
}

func server() {
  server, err := net.Listen("tcp", MASTER)
  if err != nil {
    fmt.Println("Error listening: ", err)
  }
  fmt.Println("Server Started. Waiting for Connections...")
  for {
    connection, err := server.Accept()
    if err != nil {
      fmt.Println("Error", err)
      os.Exit(1)
    }
    fmt.Println("Client connected")
    go sendData(connection)
  }
}

func sendData(connection net.Conn) {
	fmt.Println("A client has connected")
	file, err := os.Open("test/hoge1.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println(err)
		return
	}
	fileSize := fillString(strconv.FormatInt(fileInfo.Size(), 10), 10)
	fileName := fillString(fileInfo.Name(), 64)
	fmt.Println("Sending filename and filesize")
	connection.Write([]byte(fileSize))
	connection.Write([]byte(fileName))
	sendBuffer := make([]byte, BUFFERSIZE)
	fmt.Println("Start sending file")

	for {
		_, err = file.Read(sendBuffer)
		if err == io.EOF {
			break
		}
		connection.Write(sendBuffer)
	}
	fmt.Println("File has been sent, closing connection")
	return
}

func fillString(retunString string, toLength int) string {
	for {
		lengtString := len(retunString)
		if lengtString < toLength {
			retunString = retunString + ":"
			continue
		}
		break
	}
	return retunString
}

func exportImg(){
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
  fmt.Println(images)
  fmt.Println(len(images))
  fmt.Println()
  fmt.Println(images["cni.img"])

}

func microk8s_enable() {
}

func notify() {
}
