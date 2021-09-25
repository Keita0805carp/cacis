package master

import (
  "fmt"
  "log"
  "net"
  "os"

  "github.com/keita0805carp/cacis/cacis"
)

func installSnapd() {
  log.Printf("Install snap via apt\n")
  if cacis.IsCommandAvailable("snap") {
    log.Printf("[Debug] Already installed\n")
    return
  }
  cacis.ExecCmd("apt install snapd", false)
}

func sendSnapd(conn net.Conn) {
  log.Print("\n[Debug] Start Send Snapd\n")
  s := []string{"snapd.zip"}

  for _, fileName := range s {
    fileBuf := readFileByte(fileName)

    /// Send Image
    fmt.Printf("\r[Debug] Sending Snapd %s ...", fileName)
    cLayer := cacis.SendSnapd(fileBuf)
    packet := cLayer.Marshal()
    //fmt.Println(cLayer)
    conn.Write(packet)
    fmt.Printf("\r[Debug] Send Snapd %s Completely\n", fileName)
  }
  log.Print("\n[Debug] End Send Snapd\n")
}


func downloadMicrok8s() {
  log.Printf("[Debug] Download microk8s via snap\n")
  log.Printf("[Debug] Downloading...\n")
  cacis.ExecCmd("snap download microk8s --target-directory=" + targetDir, false)
  log.Printf("[Debug] Download Completely\n")
}

func sendMicrok8s(conn net.Conn) {
  log.Print("\n[Debug] Start Send Snap files\n")

  for _, fileName := range cacis.Microk8sSnaps {
    fileBuf := readFileByte(fileName)

    /// Send Image
    fmt.Printf("\r[Debug] Sending Snap files %s ...", fileName)
    cLayer := cacis.SendMicrok8sSnap(fileBuf)
    packet := cLayer.Marshal()
    //fmt.Println(cLayer)
    conn.Write(packet)
    fmt.Printf("\r[Debug] Send Snap files %s Completely\n", fileName)
  }
  log.Print("\n[Debug] End Send Snap files\n")
}

func installMicrok8s() {
  log.Printf("Install microk8s via snap\n")
  if cacis.IsCommandAvailable("microk8s") {
    log.Printf("[Debug] Already installed\n")
    return 
  } else {
    log.Printf("Install microk8s via snap\n")
    log.Printf("Installing...")
    cacis.ExecCmd("snap ack " + targetDir + cacis.Microk8sSnaps[0], false)
    cacis.ExecCmd("snap install " + targetDir + cacis.Microk8sSnaps[1] + " --classic", true)
    log.Printf("[Debug] Installed\n")
    log.Printf("[Debug] Waiting for ready\n")
    cacis.ExecCmd("microk8s status --wait-ready", false)
    log.Printf("[Debug] Install Completely\n")
  }
}

func enableMicrok8s() {
  log.Println("enable add-on...")
  cacis.ExecCmd("microk8s enable registry dns dashboard", false)
  log.Println("enabled add-on")
}

func GetKubeconfig() (string, error) {
  config, err := cacis.ExecCmd("microk8s config", false)
  return string(config), err
}

func ExportKubeconfig(path string) (error) {
  config, err := GetKubeconfig()
  file, err := os.Create(path)
  cacis.Error(err)
  _, err = file.WriteString(config)
  return err
}

