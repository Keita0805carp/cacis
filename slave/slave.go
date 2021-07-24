package slave

import (
  "fmt"
  "io"
  "net"
  "os"
  "strconv"
  "strings"
)

const BUFFERSIZE = 1024

func Main() {
  //fmt.Println("This is slave.main()")
  client()
}

func client() {
  connection, err := net.Dial("tcp", "10.0.100.1:27001")
  if err != nil {
    panic(err)
  }
  defer connection.Close()
  fmt.Println("Connected to server, start recieving the file name and file size")
  bufferFileName := make([]byte, 64)
  bufferFileSize := make([]byte, 10)

  connection.Read(bufferFileSize)
  fileSize, _ := strconv.ParseInt(strings.Trim(string(bufferFileSize), ":"), 10, 64)

  connection.Read(bufferFileName)
  fileName := strings.Trim(string(bufferFileName), ":")

  newFile, err := os.Create(fileName)
  if err != nil {
    panic(err)
  }
  defer newFile.Close()

  var receivedBytes int64

  for {
    if (fileSize - receivedBytes) < BUFFERSIZE {
      io.CopyN(newFile, connection, (fileSize -receivedBytes))
      connection.Read(make([]byte, (receivedBytes+BUFFERSIZE)-fileSize))
      break
    }
    io.CopyN(newFile, connection, BUFFERSIZE)
    receivedBytes += BUFFERSIZE
  }
  fmt.Println("Received file Completely")
}



func install_microk8s() {
}

func recieve_data() {
}

func import_img() {
}

func notify() {
}
