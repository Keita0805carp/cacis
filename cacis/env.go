package cacis

const (
  MasterIP = "10.0.100.1"
  MasterPort = "27001"
  TargetDir = "./tmp/"
  ContainerdSock = "/var/snap/microk8s/common/run/containerd.sock"
  ContainerdNameSpace = "k8s.io"
  //containerdNameSpace = "cacis"

)

var (
  ComponentsList = map[string]string {
    "cni.img"             : "docker.io/calico/cni:v3.13.2",
    "pause.img"           : "k8s.gcr.io/pause:3.1",
    "kube-controllers.img": "docker.io/calico/kube-controllers:v3.13.2",
    "pod2daemon.img"      : "docker.io/calico/pod2daemon-flexvol:v3.13.2",
    "node.img"            : "docker.io/calico/node:v3.13.2",
    "coredns.img"         : "docker.io/coredns/coredns:1.8.0",
    "metrics-server.img"  : "k8s.gcr.io/metrics-server-arm64:v0.3.6",
    "dashboard.img"       : "docker.io/kubernetesui/dashboard:v2.0.0",
  }
)
