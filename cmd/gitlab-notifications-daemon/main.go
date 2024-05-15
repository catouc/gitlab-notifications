package main

import (
	"net"
	"os"
	"os/signal"
	"time"
	"strings"
	"syscall"

        flag "github.com/spf13/pflag"

	"github.com/catouc/gitlab-notifications/internal/gitlab"
	"github.com/catouc/gitlab-notifications/internal/notifier"
)

const (
	socketPath = "/tmp/gitlab-notifications.sock"
)

var (
	pollDuration = flag.DurationP("poll-duration", "p", 30 * time.Second, "Specify how often to poll GitLab for updates.")
	gitlabToken = flag.StringP("gitlab-token", "t", "", "Set the GitLab auth token to interact with the API, needs to have write permissions for auto merges")
	gitlabBaseURL = flag.StringP("gitlab-base-url", "b", "", "Set the base url for the gitlab instance, this is the url we will hit for all API calls.")
)

func main() {
	baseURL := decideStringFlagValue(gitlabBaseURL, "GITLAB_BASE_URL")
	if baseURL == "" {
		flag.PrintDefaults()
		return
	}

	token := decideStringFlagValue(gitlabToken, "GITLAB_TOKEN")
	if token == "" {
		flag.PrintDefaults()
		return
	}

	sock, err := net.Listen("unix", socketPath)
	if err != nil {
		panic(err)
	}
	
	gitlabClient := gitlab.New(baseURL, token, time.Second*5)

	// Cleanup the sockfile.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		os.Remove(socketPath)
		os.Exit(1)
	}()

	for {
		conn, err := sock.Accept()
		if err != nil {
			panic(err)
		}

		go func(conn net.Conn) {
		    defer conn.Close()
		    // Create a buffer for incoming data.
		    buf := make([]byte, 4096)

		    // Read data from the connection.
		    n, err := conn.Read(buf)
		    if err != nil {
			panic(err)
		    }

		    go notifier.NewTicker(
			    gitlabClient,
			    strings.TrimSpace(string(buf[:n])),
			    *pollDuration,
		    )
		}(conn)
	}
}

func decideStringFlagValue(flag *string, envVarName string) string {
	switch {
	case *flag == "" && os.Getenv(envVarName) == "":
		return ""
	case *flag == "" && os.Getenv(envVarName) != "":
		return os.Getenv(envVarName)
	case *flag != "" && os.Getenv(envVarName) == "":
		return *flag
	default:
		return ""
	}
}
