package main

import (
	"fmt"
	"net"
	"os"
	"syscall"

	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

const (
	namespace    = "kube-system"
	endpointName = "kube-dns"
)

func main() {
	// If 'kubernetes' dns name resolves, no need to modify /etc/resolv.conf
	addrs, err := net.LookupHost("kubernetes")
	if err == nil {
		if len(addrs) != 0 {
			fmt.Printf("kube-dns accessible from host; nothing todo...\n")
			os.Exit(0)
		}
	}
	fmt.Printf("kude-dns NOT accessible from host - modifying /etc/resolv.conf\n")

	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// creates the client
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	var endpoint *v1.Endpoints
	endpoint, err = client.CoreV1().Endpoints(namespace).Get(endpointName, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		fmt.Fprintln(os.Stderr, "kube-dns endpoint not found")
	} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
		fmt.Printf("error getting endpoint %v\n", statusError.ErrStatus.Message)
	} else if err != nil {
		panic(err.Error())
	}
	clusterIP := endpoint.Subsets[0].Addresses[0].IP
	fmt.Printf("kube-dns endpoint %s\n", clusterIP)

	// Enter PID 1 namespace and replace the /etc/resolv.conf file
	err = syscall.Exec("/usr/bin/nsenter1", []string{"nsenter1", "/bin/sh", "-c",
		fmt.Sprintf(`cp /etc/resolv.conf /tmp/resolv.conf; printf "nameserver %s\n`+
			`nameserver 192.168.65.1\n`+
			`search default.svc.cluster.local svc.cluster.local cluster.local\n`+
			`options ndots:5\n" > /etc/resolv.conf`, clusterIP)}, os.Environ())
	if err != nil {
		panic(err.Error())
	}
}
