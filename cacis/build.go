package cacis

import (
  "log"
)

const (
  registry = "localhost:32000"
)

func BuildFromDockerfile(dockerfilePath, imageTag string) {
  tag := registry + "/" + imageTag
  log.Println("Build with tag")
  ExecCmd("docker build -t " + tag + " " + dockerfilePath, true)
  log.Println("Push")
  ExecCmd("docker push " + tag, true)
}

func PullAndPush(imageRef string) {
  tag := registry + "/" + imageRef
  log.Println("Docker pull")
  ExecCmd("docker pull " + imageRef, true)
  log.Println("Docker tag")
  ExecCmd("docker tag " + imageRef + " " + tag, false)
  log.Println("Push")
  ExecCmd("docker push " + tag, true)
}

