package cacis

const (
  MasterIP = "10.0.100.1"
  MasterPort = "27001"
  TargetDir = "./tmp/"
  ContainerdSock = "/var/snap/microk8s/common/run/containerd.sock"
  ContainerdNameSpace = "k8s.io"
  //containerdNameSpace = "cacis"
  DhcpConfPath = "config/dhcp.conf"
  NetplanConfTemplatePath = "config/netplan.conf.template"
  NetplanConfPath = "/etc/netplan/60-cacis.yaml"
  HostapdConfTemplatePath = "config/hostapd.conf.template"
  HostapdConfPath = "config/hostapd.conf"
)

var (
  NodeLabels = map[string]string {
    "testLabel1"           : "test1",
    "testLabel2"           : "test2",
  }
  ComponentsList = map[string]string {
    "cni.img"             : "docker.io/calico/cni:v3.19.1",
    "pause.img"           : "k8s.gcr.io/pause:3.1",
    "kube-controllers.img": "docker.io/calico/kube-controllers:v3.17.3",
    "pod2daemon.img"      : "docker.io/calico/pod2daemon-flexvol:v3.19.1",
    "node.img"            : "docker.io/calico/node:v3.19.1",
    "coredns.img"         : "docker.io/coredns/coredns:1.8.0",
    "metrics-server.img"  : "k8s.gcr.io/metrics-server/metrics-server:v0.5.0",
    "dashboard.img"       : "docker.io/kubernetesui/dashboard:v2.2.0",
    "hostpath-arm64.img"  : "docker.io/cdkbot/hostpath-provisioner-arm64:1.0.0",
    "registry-arm64.img"  : "docker.io/cdkbot/registry-arm64:2.6",
    "metrics-scraper.img" : "docker.io/kubernetesui/metrics-scraper:v1.0.6",
  }
  Microk8sSnap = "microk8s_cacis"
  Microk8sSnaps = []string {"microk8s_cacis.assert", "microk8s_cacis.snap"}
  Snapd = []string{"core_11420.assert", "core_11420.snap"}
)
