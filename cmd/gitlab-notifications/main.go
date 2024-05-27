package main

import (
	"fmt"
	"net"
	"os"
)

const socketPath = "/tmp/gitlab-notifications.sock"

func main() {
	url := os.Args[0]

	sock, err := net.Dial("unix", socketPath)	
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to establish connection to /tmp/gitlab-notifications.sock %s, is the daemon running?", err)
		os.Exit(1)
	}

	_, err = sock.Write([]byte(url))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to write data to socket: %s", err)
		os.Exit(1)
	}
}
