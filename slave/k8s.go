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

func recieveMicrok8s(listen net.Listener) {
  log.Println("[Debug] Start RECIEVE SNAP FILES")
  // Socket
  conn2master, err := net.Dial("tcp", masterIP+":"+masterPort)
  cacis.Error(err)
  log.Printf("[Debug] Request Snap files\n")
  cLayer := cacis.RequestMicrok8sSnap()
  packet := cLayer.Marshal()
  //fmt.Println(packet)
  conn2master.Write(packet)
  log.Printf("Requested\n\n")
  conn2master.Close()

  conn2slave, err := listen.Accept()
  cacis.Error(err)

  for _, fileName := range cacis.Microk8sSnaps {
    recieveFile(conn2slave, fileName)
  }
  log.Println("[Debug] End RECIEVE SNAP FILES")
  conn2slave.Close()
}

func installMicrok8s() {
  log.Printf("Install microk8s via snap\n")
  log.Printf("Installing...")
  cacis.ExecCmd("snap ack " + targetDir + cacis.Microk8sSnaps[0], false)
  cacis.ExecCmd("snap install " + targetDir + cacis.Microk8sSnaps[1] + " --classic", true)
  log.Printf("[Debug] Installed\n")
}

func WaitReadyMicrok8s() {
  log.Printf("[Debug] Waiting for ready\n")
  cacis.ExecCmd("microk8s status --wait-ready", false)
  log.Printf("[Debug] Microk8s is Ready\n")
}

func RemoveMicrok8s() {
  cacis.ExecCmd("microk8s stop", true)
  cacis.ExecCmd("microk8s reset --destroy-storage", true)
  cacis.ExecCmd("snap remove --purge microk8s", true)
}

