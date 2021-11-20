package slave

import (
  "os"
  "log"
  "net"
  "time"
  "bufio"

  "github.com/keita0805carp/cacis/cacis"
)

func installSnapd() {
  log.Printf("[Debug] Check Snap\n")
  if cacis.IsCommandAvailable("snap") {
    /// if debian(raspberry pi os)
      //recieve zip
      //unzip
    return
  }

  log.Printf("[Debug] Start RECIEVE snapd and install\n")
  // Socket
  conn, err := net.Dial("tcp", masterIP+":"+masterPort)
  cacis.Error(err)
  defer conn.Close()

  log.Printf("[Debug] Request snapd\n")
  cLayer := cacis.RequestSnapd()
  packet := cLayer.Marshal()
  //fmt.Println(packet)
  conn.Write(packet)
  log.Printf("Requested\n")

  recieveFile(conn, "snapd.zip")

  log.Printf("[Debug] End RECIEVE COMPONENT IMAGES\n")
  cacis.ExecCmd("dpkg -i ./*.deb", false)
  //reboot
  cacis.ExecCmd("snap install core", false)
}

func recieveMicrok8s(listen net.Listener) {
  log.Printf("[Debug] Start RECIEVE SNAP FILES\n")
  // Socket
  conn2master, err := net.Dial("tcp", masterIP+":"+masterPort)
  cacis.Error(err)
  log.Printf("[Debug] Request Snap files\n")
  cLayer := cacis.RequestMicrok8sSnap()
  packet := cLayer.Marshal()
  //fmt.Println(packet)
  conn2master.Write(packet)
  log.Printf("Requested\n")
  conn2master.Close()

  conn2slave, err := listen.Accept()
  cacis.Error(err)

  for _, fileName := range cacis.Microk8sSnaps {
    recieveFile(conn2slave, fileName)
  }
  log.Printf("[Debug] End RECIEVE SNAP FILES\n")
  conn2slave.Close()
}

func installMicrok8s() {
  log.Printf("Install microk8s via snap\n")
  log.Printf("Installing...\n")
  cacis.ExecCmd("snap ack " + targetDir + cacis.Microk8sSnaps[0], false)
  cacis.ExecCmd("snap install " + targetDir + cacis.Microk8sSnaps[1] + " --classic", true)
  log.Printf("[Debug] Installed\n")
}

func WaitReadyMicrok8s() {
  log.Printf("[Debug] Waiting for ready\n")
  _, err := cacis.ExecCmd("microk8s status --wait-ready --timeout 15", false)
  if err != nil {
    log.Printf("[Error] Microk8s is NOT Ready in time\n")
    log.Printf("[Error] Please try restart or reinstall microk8s.\n")
    cacis.Error(err)
  }
  log.Printf("[Debug] Microk8s is Ready\n")
}

func IsClustered() bool {
  var lines []string
  file, err := os.Open("/var/snap/microk8s/current/var/kubernetes/backend/cluster.yaml")
  cacis.Error(err)

  scanner := bufio.NewScanner(file)

  for scanner.Scan() {
    line := scanner.Text()
    lines = append(lines, line)
  }
  file.Close()

  if len(lines) > 3 {
    return true
  }
  return false
}

func clustering(listen net.Listener) {
  log.Printf("[Debug] Start CLUSTERING\n")
  if IsClustered() {
    log.Printf("[Warn] This node join to cluster already.\n")
    log.Printf("[Warn] 'cacis slave --leave' to leave cluster\n")
    return
  }
  // Socket
  conn2master, err := net.Dial("tcp", masterIP+":"+masterPort)
  cacis.Error(err)
  log.Printf("[Debug] Request Clustering\n")
  cLayer := cacis.RequestClustering()
  packet := cLayer.Marshal()
  //fmt.Println(packet)
  conn2master.Write(packet)
  log.Printf("Requested\n")
  conn2master.Close()

  conn2slave, err := listen.Accept()
  cacis.Error(err)

  packet = make([]byte, cacis.CacisLayerSize)
  packetLength, err := conn2slave.Read(packet)
  cacis.Error(err)
  log.Printf("[Debug] Read Packet HEADER. len: %d\n", packetLength)
  cLayer = cacis.Unmarshal(packet)
  log.Printf("Debug: Read Packet PAYLOAD\n")
  cLayer.Payload = loadPayload(conn2slave, cLayer.Length)

  log.Printf("[Debug] Clustering...\n")
  cacis.ExecCmd(string(cLayer.Payload), true)

  log.Printf("[Debug] End CLUSTERING\n")
  conn2slave.Close()
}

func Unclustering() {
  if !IsClustered() {
    log.Printf("[Error]: Single Node Cluster\n")
    return 
  }
  log.Printf("[Debug] Start UNCLUSTERING\n")
  log.Printf("[Debug] Leaving...\n")
  go cacis.ExecCmd("microk8s leave", false)
  time.Sleep(time.Second * 20)
  // Socket
  conn, err := net.Dial("tcp", masterIP+":"+masterPort)
  cacis.Error(err)
  defer conn.Close()

  hostname, err := os.Hostname()
  cacis.Error(err)

  log.Printf("[Debug] Request Unclustering\n")
  cLayer := cacis.RequestUnclustering([]byte(hostname))
  packet := cLayer.Marshal()
  //fmt.Println(packet)
  conn.Write(packet)
  log.Printf("[Debug] End UNCLUSTERING\n")
}

func labelNode() {
  log.Printf("[Debug] Node Labeling...\n")
  hostname, err := os.Hostname()
  cacis.Error(err)
  cmd := "microk8s kubectl label nodes " + hostname
  for key, value := range cacis.NodeLabels {
    cmd += " " + key + "=" + value
  }
  WaitReadyMicrok8s()
  time.Sleep(time.Second * 3)
  cacis.ExecCmd(cmd, false)
  log.Printf("[Debug] Node Labeled\n")
}

func RemoveMicrok8s() {
  //cacis.ExecCmd("microk8s stop", true)
  //cacis.ExecCmd("microk8s reset --destroy-storage", true)
  cacis.ExecCmd("snap remove --purge microk8s", true)
}

