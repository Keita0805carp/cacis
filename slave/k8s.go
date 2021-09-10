package slave

import (
  "log"
  "net"

  "github.com/keita0805carp/cacis/cacis"
)

func installSnapd() {
  log.Println("[Debug] Check Snap")
  if cacis.IsCommandAvailable("snap") {
    /// if debian(raspberry pi os)
      //recieve zip
      //unzip
    return
  }

  log.Println("[Debug] Start RECIEVE snapd and install")
  // Socket
  conn, err := net.Dial("tcp", masterIP+":"+masterPort)
  cacis.Error(err)
  defer conn.Close()

  log.Println("[Debug] Request snapd")
  cLayer := cacis.RequestSnapd()
  packet := cLayer.Marshal()
  //fmt.Println(packet)
  conn.Write(packet)
  log.Printf("Requested\n\n")

  recieveFile(conn, "snapd.zip")

  log.Println("[Debug] End RECIEVE COMPONENT IMAGES")
  cacis.ExecCmd("dpkg -i ./*.deb", false)
  //reboot
  cacis.ExecCmd("snap install core", false)
}

func recieveMicrok8s() {
  log.Println("[Debug] Start RECIEVE SNAP FILES")
  // Socket
  conn, err := net.Dial("tcp", masterIP+":"+masterPort)
  cacis.Error(err)
  defer conn.Close()

  log.Println("[Debug] Request Snap files")
  cLayer := cacis.RequestMicrok8sSnap()
  packet := cLayer.Marshal()
  //fmt.Println(packet)
  conn.Write(packet)
  log.Printf("Requested\n\n")

  for _, fileName := range cacis.Microk8sSnaps {
    recieveFile(conn, fileName)
  }
  log.Println("[Debug] End RECIEVE SNAP FILES")
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

func removeMicrok8s() {
  cacis.ExecCmd("microk8s stop", false)
  cacis.ExecCmd("microk8s reset --destroy-storage", false)
  cacis.ExecCmd("snap remove --purge microk8s", false)
  cacis.ExecCmd("apt purge snap", false)
}

