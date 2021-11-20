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
      log.Printf("Unstable WiFi Connection\n")
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
    log.Printf("[Info] RSSI: %d\n", rssi)
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

  rssi := -999
  if len(lines) > 2 {
    info := strings.Fields(lines[2])
    parse := strings.Replace(info[3], ".", "", 1)
    rssi, err =  strconv.Atoi(parse)
    cacis.Error(err)
  }
  //log.Printf("[Info] RSSI: %d\n", rssi)
  return rssi
}
