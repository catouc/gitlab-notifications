package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gen2brain/beeep"
        flag "github.com/spf13/pflag" 

	"github.com/catouc/gitlab-notifications/internal/gitlab"
	"github.com/catouc/gitlab-notifications/internal/notifier"
)

var (
	pollDuration = flag.DurationP("poll-duration", "p", time.Minute, "Specify how often this notifications code polls GitLab for an update")
)

func main() {
	flag.Parse()
	url := flag.Arg(0)

	baseURL := os.Getenv("GITLAB_BASE_URL")
	token := os.Getenv("GITLAB_TOKEN")
	gitlabClient := gitlab.New(baseURL, token, time.Second*5)

	// TODO: just move this stuff into a struct but eh parsing the
	// url again every minute is going to be fine...
	var notifyFunc func(c *gitlab.Client, url string) (bool, error)
	switch gitlabClient.IdentifyGitlabResourceFromURL(url) {
	case "mr":
		notifyFunc = notifier.MRNotifyFunc
	case "job":
		notifyFunc = notifier.JobNotifyFunc
	case "pipeline":
		notifyFunc = notifier.PipelineNotifyFunc
	default:
		fmt.Fprintln(os.Stderr, "could not find a valid resource from the passed url")
	}

	notificationTicker := time.NewTicker(*pollDuration)
	for {
		select {
			case <- notificationTicker.C:
				notify, err := notifyFunc(gitlabClient, url)
				if err != nil {
					fmt.Fprintf(os.Stderr, "error getting notification status: %s", err)
					os.Exit(1)
				}

				if !notify {
					continue	
				}

				err = beeep.Notify("Your GitLab thing is ready", url, "")
				if err != nil {
					fmt.Fprintf(os.Stderr, "could not notify you of the completed resource: %s", err)
					os.Exit(1)
				}
				os.Exit(0)
		}
	}
}
