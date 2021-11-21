# cacis
Container-based Adaptive Clustering IoT System

### Requirements

- Hardware
  - Edge devices (RaspberryPi/Jetson/...)
  - Wi-Fi Adapter
  - Bluetooth Adapter
- Software
  - Go (1.16.6)
  - Snap
  - Microk8s
  - hostapd
  - netplan
  - avahi-daemon (mDNS)
  - bluez, pi-bluetooth

### Setup
```
## Add boot parameter (for Raspberry Pi)
Add `cgroup_enable=memory cgroup_memory=1` to `/boot/firmware/cmdline.txt`

## Install Go
$ rm -rf /usr/local/go && tar -C /usr/local -xzf go1.16.linux-amd64.tar.gz
$ export PATH=$PATH:/usr/local/go/bin
$ go version

## Install Snapd
$ sudo apt install snapd
$ sudo snap instal core18

## Install Microk8s (Option)
$ sudo snap install microk8s --classic

## For Master
$ sudo apt install hostapd

## For Slave
$ sudo apt install netplan.io
```

### Default Add-on
- registry
- dashboard
- dns

### How to run Docker image

```
docker tag [local image] [localhost:32000/hoge:fuga]
docker push [localhost:32000/hoge:fuga]
```

Then, create kuberentes yaml and apply it.

```
kubectl apply -f [piyo.yaml]
```

