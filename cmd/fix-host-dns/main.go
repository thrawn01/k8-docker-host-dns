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
		fmt.Sprintf(`cat /etc/resolv.conf | grep fix-host-dns; [ $? -eq 0 ] || `+
			`(cp /etc/resolv.conf /etc/resolv.bak;`+
			`rm /etc/resolv.conf;`+ // resolv.conf might be a link to systemd managed version (SEE NOTE)
			`printf "# Modified by fix-host-dns\n` +
			`nameserver %s\n`+
			"`cat /etc/resolv.bak | tail -1`\n"+
			`search default.svc.cluster.local svc.cluster.local cluster.local\n`+
			`options ndots:5\n" > /etc/resolv.conf)`, clusterIP)}, os.Environ())
	if err != nil {
		panic(err.Error())
	}

	// NOTE: systemd-resolved does not support nameserver fall through as defined in man /etc/resolv.conf
	// and it probably never will. As such we have to remove the link to the systemd managed
	// resolv.conf so our containers will use our properly formed /etc/resolv.conf file
	// See https://github.com/systemd/systemd/issues/5755

}
