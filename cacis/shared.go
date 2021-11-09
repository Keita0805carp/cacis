package cacis

import (
  "fmt"
  "log"
  "sort"
  "os"
  "os/exec"
  "strings"
)

func Error(err error) {
  if err != nil {
    log.Fatal(err)
  }
}

func CreateTempDir() {
  log.Printf("[Debug] Create Temporary Directory")
  os.Mkdir(TargetDir, 0775)
}

func RemoveTempDir() {
  log.Printf("[Debug] Remove Temporary Directory")
  err := os.RemoveAll(TargetDir)
  Error(err)
}

func IsCommandAvailable(command string) bool {
  slice := strings.Split(command, " ")
  _, err := exec.LookPath(slice[0])
  if err != nil {
    log.Printf("[Debug] command check: '%s' %s\n", slice[0], "Fail")
    return false
  }
  log.Printf("[Debug] command check: '%s' %s\n", slice[0], "Success")
  return true
}

func ExecCmd(command string, log bool) ([]byte, error) {
  slice := strings.Split(command, " ")
  stdout, err := exec.Command(slice[0], slice[1:]...).Output()
  if log {
    fmt.Println(string(stdout))
    Error(err)
  }
  return stdout, err
}

func SortKeys(m map[string]string) []string {
  ///sort
  sorted := make([]string, len(m))
  index := 0
  for key := range m {
        sorted[index] = key
        index++
    }
    sort.Strings(sorted)
  /*
  for _, exportFile := range exportFileNameSort {
    fmt.Printf("%-20s : %s\n", exportFile, componentsList[exportFile])
  }
  */
  return sorted
}

