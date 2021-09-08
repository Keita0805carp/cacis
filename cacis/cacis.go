package cacis

import (
  "fmt"
  "log"
  "strings"
  "os/exec"
  "encoding/json"
  "encoding/binary"
)

type CacisLayer struct {
  Type    uint8
  Length  uint64
  Payload []byte
}

const CacisLayerSize = 1 + 8

func (c *CacisLayer) Marshal() []byte {
  buf := make([]byte, CacisLayerSize)
  buf[0]   = byte(c.Type)
  binary.BigEndian.PutUint64(buf[1:], c.Length)
  buf      = append(buf, c.Payload...)
  return buf
}

func Unmarshal(buf []byte) CacisLayer {
  var c CacisLayer
  c.Type   = uint8(buf[0])
  c.Length = binary.BigEndian.Uint64(buf[1:])
  c.Payload = buf[9:]
  return c
}

func NewCacisPacket(cacisType uint8, l uint64, p []byte) CacisLayer {
  return CacisLayer{
    Type:    cacisType,
    Length:  l,
    Payload: []byte(p),
  }
}

// Request
func RequestComponentsList() CacisLayer {
  return NewCacisPacket(10, 0, nil)
}

func RequestImage() CacisLayer {
  return NewCacisPacket(20, 0, nil)
}

func RequestMicrok8sSnap() CacisLayer {
  return NewCacisPacket(30, 0, nil)
}

func RequestSnapd() CacisLayer {
  return NewCacisPacket(40, 0, nil)
}

func RequestClustering() CacisLayer {
  return NewCacisPacket(50, 0, nil)
}

func RequestUnclustering(p []byte) CacisLayer {
  return NewCacisPacket(60, uint64(len(p)), p)
}

// Send
func SendComponentsList(list map[string]string) CacisLayer {
  p, err := json.Marshal(list)
  Error(err)
  return NewCacisPacket(11, uint64(len(p)), p)
}

func SendImage(p []byte) CacisLayer {
  return NewCacisPacket(21, uint64(len(p)), p)
}

func SendMicrok8sSnap(p []byte) CacisLayer {
  return NewCacisPacket(31, uint64(len(p)), p)
}

func SendSnapd(p []byte) CacisLayer {
  return NewCacisPacket(41, uint64(len(p)), p)
}

func SendClusterInfo(p []byte) CacisLayer {
  return NewCacisPacket(51, uint64(len(p)), p)
}


func Error(err error) {
  if err != nil {
    log.Fatal(err)
  }
}

func ExecCmd(cmd string, log bool) ([]byte, error) {
  slice := strings.Split(cmd, " ")
  stdout, err := exec.Command(slice[0], slice[1:]...).Output()
  if log {
    fmt.Println(string(stdout))
    Error(err)
    return stdout, err
  } else {
    return nil, nil
  }
  //fmt.Printf("exec: %s\noutput:\n%s", cmd, stdout)
}

func ExecCmd(cmd string, log bool) ([]byte, error) {
  slice := strings.Split(cmd, " ")
  stdout, err := exec.Command(slice[0], slice[1:]...).Output()
  if log {
    fmt.Println(string(stdout))
    Error(err)
    return stdout, err
  } else {
    return nil, nil
  }
  //fmt.Printf("exec: %s\noutput:\n%s", cmd, stdout)
}

