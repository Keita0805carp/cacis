package master

import (
  "fmt"
  "log"
  "net"

  "github.com/keita0805carp/cacis/cacis"
)

func snapd() {
  //TODO snapd check
  //TODO snapd install
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

func sendMicrok8sSnap(conn net.Conn) {
  log.Print("\n[Debug] Start Send Snap files\n")
  s := []string{"microk8s_2347.assert", "microk8s_2347.snap", "core_11420.assert", "core_11420.snap"}

  for _, fileName := range s {
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
    cacis.ExecCmd("snap ack " + targetDir + "microk8s_2347.assert", false)
    cacis.ExecCmd("snap install " + targetDir + "microk8s_2347.snap" + " --classic", true)
    log.Printf("[Debug] Installed\n")
    log.Printf("[Debug] Waiting for ready\n")
    cacis.ExecCmd("microk8s status --wait-ready", false)
    log.Printf("[Debug] Install Completely\n")
  }
}

func enableMicrok8s() {
  cacis.ExecCmd("microk8s enable dns dashboard", false)
}

func getKubeconfig() {
  cacis.ExecCmd("microk8s config", true)
}

