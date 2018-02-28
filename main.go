package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

func setEnv(env, value string) {
	fmt.Printf("%s=%s\n", env, value)
	os.Setenv(env, value)
}

func main() {

	for {
		cmd := []string{"get", "endpoints", "etcd-cluster", "-o", "jsonpath={.subsets[].addresses[].ip}"}
		etcdPeerBytes, err := exec.Command("/Users/thrawn/.devgun/bin/kubectl", cmd...).Output()
		if err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok {
				fmt.Fprintf(os.Stderr, "%s - %s", strings.Join(cmd, ","), exitErr.Stderr)
				time.Sleep(time.Second)
				continue
			}
			fmt.Fprintf(os.Stderr, "-- fatal error while running %s - %s", strings.Join(cmd, ","), err)
			os.Exit(1)
		}
		etcdPeerIP := strings.TrimRight(string(etcdPeerBytes), "\n")
		setEnv("ETCD_CLUSTER_IP", etcdPeerIP)
		break
	}
}
