package connection

import (
  "os"
  "log"
  "time"
  "bufio"
  "strconv"
  "strings"

  "github.com/keita0805carp/cacis/cacis"
)

func UnstableWifiEvent(cancel chan struct{}) {
  cntUnstable := 0
  for {
    if !isWifiStable() {
      cntUnstable++
    } else {
      cntUnstable = 0
    }
    if cntUnstable > 5 {
      log.Println("Unstable WiFi Connection")
      close(cancel)
      return
    }
    time.Sleep(time.Second * 5)
  }
}

func isWifiStable() bool {
  rssi := getRSSI()
  if (rssi > -70) {
    return true
  } else {
    return false
  }
}

func getRSSI() int {
  var lines []string
  file, err := os.Open("/proc/net/wireless")
  cacis.Error(err)

  scanner := bufio.NewScanner(file)

  for scanner.Scan() {
    line := scanner.Text()
    lines = append(lines, line)
  }
  file.Close()

  info := strings.Fields(lines[2])
  rssi, err :=  strconv.Atoi(info[3][:len(info[3])-1])
  cacis.Error(err)
  //log.Printf("RSSI: %d\n", signal)
  return rssi
}
