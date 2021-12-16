package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	masterURL      string
	kubeconfig     string
	addr           = flag.String("listen-address", ":8080", "The address to listen on for HTTP requests.")
	OpenPorts      *prometheus.GaugeVec
	startPortRange = os.Getenv("startPortRange")
	endPortRange   = os.Getenv("endPortRange")
)

type Status struct {
	Status string `json:"status"`
}

func main() {

	// Build Kubernetes config
	kubeCfg, err := clientcmd.BuildConfigFromFlags(masterURL, kubeconfig)
	if err != nil {
		log.Fatalf("Error building kubeconfig: %s", err.Error())
	}

	log.Info("Built config from flags...")

	// Create Kubernetes client
	kubeClient, err := kubernetes.NewForConfig(kubeCfg)
	if err != nil {
		log.Fatalf("Error building watcher clientset: %s", err.Error())
	}

	// Get port range from ebv variables to int

	spr, err := strconv.Atoi(startPortRange)
	if err != nil {
		log.Errorf("Could not get port start range")
	}
	epr, err := strconv.Atoi(endPortRange)
	if err != nil {
		log.Errorf("Could not get port end range")
	}

	// Prometheus Http handler
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.HandleFunc("/health", healthStatus)
		log.Fatal(http.ListenAndServe(*addr, nil))
	}()

	// Prometheus metrics

	OpenPorts = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "portscanexporter_open_ports",
		Help: "Number of ports that are open.",
	}, []string{"Name", "IP", "Namespace", "protocol", "OpenPortNumber"})

	prometheus.MustRegister(OpenPorts)
	// Runtime loop

	for {

		// Fetch pods in cluster
		podList, err := kubeClient.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			log.Errorf("Error fetching pods: %s", err.Error())
		} else {
			for _, lst := range podList.Items {
				if !lst.Spec.HostNetwork {
					for i := spr; i <= epr; i++ {

						address := lst.Status.PodIP + ":" + strconv.Itoa(i)
						conn, err := net.DialTimeout("tcp", address, 5*time.Second)

						if err == nil {
							conn.Close()
							OpenPorts.WithLabelValues(lst.Name, lst.Status.PodIP, lst.Namespace, "tcp", fmt.Sprint(i)).Set(float64(1))
						}
						conn, err = net.DialTimeout("udp", address, 5*time.Second)

						if err == nil {
							conn.Close()
							OpenPorts.WithLabelValues(lst.Name, lst.Status.PodIP, lst.Namespace, "udp", fmt.Sprint(i)).Set(float64(1))
						}

					}

				}
			}

		}
		time.Sleep(time.Second * 600)

	}

}

func healthStatus(w http.ResponseWriter, req *http.Request) {

	status := Status{Status: "ok"}
	json.NewEncoder(w).Encode(status)
}

func init() {
	flag.StringVar(&kubeconfig, "kubeconfig", "", "Path to a kubeconfig. Only required if out-of-cluster.")
	flag.StringVar(&masterURL, "master", "", "The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster.")
	flag.Parse()
}
