package master

import (
  "os"
  "fmt"
  "net"
  "log"
  "time"
  "regexp"
  "strings"

  "github.com/keita0805carp/cacis/cacis"
)

func installSnapd() {
  log.Printf("[Debug] Install snap via apt\n")
  if cacis.IsCommandAvailable("snap") {
    log.Printf("[Debug] Already installed\n")
    return
  }
  cacis.ExecCmd("apt install snapd", false)
}

func sendSnapd(conn net.Conn) {
  log.Print("[Debug] Start Send Snapd\n")
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
  cacis.ExecCmd("snap download microk8s --channel=latest/stable --basename=" + cacis.Microk8sSnap + " --target-directory=" + targetDir, false)
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
  log.Printf("[Debug] Install microk8s via snap\n")
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
  return config, err
}

func ExportKubeconfig(path string) (error) {
  config, err := GetKubeconfig()
  file, err := os.Create(path)
  cacis.Error(err)
  _, err = file.WriteString(config)
  return err
}

func clustering(conn net.Conn) {
  log.Printf("[Debug] Start Clustering\n")
  stdout, err := cacis.ExecCmd("microk8s add-node", false)
  cacis.Error(err)
  //fmt.Println(string(output))
  regex := regexp.MustCompile("microk8s join " + masterIP + ".*")
  joinCmd := regex.FindAllStringSubmatch(stdout, 1)[0][0]

  /// Send CLuster Info
  cLayer := cacis.SendClusterInfo([]byte(joinCmd))
  packet := cLayer.Marshal()
  //fmt.Println(cLayer)
  conn.Write(packet)
  log.Printf("[Debug] End Clustering\n")
}

func unclustering(conn net.Conn, cLayer cacis.CacisLayer) {
  log.Printf("[Debug] Start Unclustering\n")

  buf := make([]byte, cLayer.Length)
  packetLength, err := conn.Read(buf)
  cacis.Error(err)
  log.Printf("[Debug] Read Packet PAYLOAD. len: %d\n", packetLength)
  //fmt.Println(string(buf))

  cacis.ExecCmd("microk8s remove-node " + string(buf), true)
  log.Printf("[Debug] End Unclustering\n")
}

func removeNotReadyNode() {
  nodes := make(map[string]int)
  for {
    time.Sleep(time.Second * 30)
    nodes = getNodeStatus(nodes)
    for k, v := range nodes {
      if k != "master" && v > 5 {
        log.Printf("Node '%s' is unstable. Force remove...\n", k)
        stdout, err := cacis.ExecCmd("microk8s remove-node " + k, false)
        if stdout != "" || err != nil {
          cacis.ExecCmd("microk8s remove-node " + k + " --force", false)
        }
        delete(nodes, k)
      }
    }
  }
}

func getNodeStatus(nodes map[string]int) map[string]int {
  stdout, err := cacis.ExecCmd("microk8s kubectl get nodes", false)
  if err != nil {
    log.Printf("[Error] Failed to get nodes.\n")
    return nodes
  }
  lines := strings.Split(stdout, "\n")
  lines = lines[1:len(lines)-1]
  
  for _, line := range lines {
    info := strings.Fields(line)
    node := info[0]
    status := info[1]
    if status == "Ready" {
      nodes[node] = 0
    } else {
      nodes[node] += 1
    }
  }
  //fmt.Println(nodes)
  return nodes
}

