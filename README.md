# portscan-prometheus-exporter

This repository contains:
1) Create a port-scan-exporter service written in Go.
The Exporter scans each pod and collect metrics about what ports are open.
A sample prometheus metric exposed is :

```
portscanexporter_open_ports{IP="",Name="nginx",Namespace="kube-system",OpenPortNumber="8079",protocol="udp"} 1
portscanexporter_open_ports{IP="",Name="nginx",Namespace="kube-system",OpenPortNumber="8080",protocol="tcp"} 1
```
2) Helm chart to deploy the exporter.
To deploy the chart specify the Kubernetes namespace in the Makefile and run the following :
```
make deploy
```
3) A Dockerfile to package the code to docker container


4) A make file that has following targets:

- image (For Building Docker Image)
- push (For Pushing to Container Registry)
- scan (Perform a Vulnerability Scan On Image Using Trivy)
- deploy (To Deploy the Helm Chart)


