package notifier

import (
	"context"
	"fmt"
	"time"
	"os"

	"github.com/gen2brain/beeep"

	"github.com/catouc/gitlab-notifications/internal/gitlab"
)


func NewTicker(gClient *gitlab.Client, url string, pollDur time.Duration) {
	var notifyFunc func(c *gitlab.Client, url string) (bool, error)
	gitlabResourceIdentifier := gClient.IdentifyGitlabResourceFromURL(url)
	switch gitlabResourceIdentifier {
	case "mr":
		notifyFunc = MRNotifyFunc
	case "job":
		notifyFunc = JobNotifyFunc
	case "pipeline":
		notifyFunc = PipelineNotifyFunc
	default:
		fmt.Fprintln(os.Stderr, "could not find a valid resource from the passed url")
		return
	}

	notificationTicker := time.NewTicker(pollDur)
	for {
		select {
		case <- notificationTicker.C:
			notify, err := notifyFunc(gClient, url)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error getting notification status: %s", err)
				return
			}

			if !notify {
				continue	
			}

			err = beeep.Notify(fmt.Sprintf("Your GitLab %s thing is ready", gitlabResourceIdentifier), url, "")
			if err != nil {
				fmt.Fprintf(os.Stderr, "could not notify you of the completed resource: %s", err)
			}
			return
		}
	}

}

func MRNotifyFunc(c *gitlab.Client, url string) (bool, error) {
	mr, err := c.GetMergeRequestApprovals(context.TODO(), url)
	if err != nil {
		return false, err
	}

	if mr.ApprovalsLeft != 0 {
		return false, nil
	}

	return true, nil
}

func JobNotifyFunc(c *gitlab.Client, url string) (bool, error) {
	job, err := c.GetJob(context.TODO(), url)
	if err != nil {
		return false, err
	}

	// assuming the same status enum as for pipelines
	switch job.Status {
	case "success":
		return true, nil
	case "failed":
		return true, nil
	case "canceled":
		return true, nil
	case "skipped":
		return true, nil
	default:
		return false, nil
	}
}

func PipelineNotifyFunc(c *gitlab.Client, url string) (bool, error) {
	pipeline, err := c.GetPipeline(context.TODO(), url)
	if err != nil {
		return false, err
	}

	// The status of pipelines, one of:
	// created, waiting_for_resource, preparing, pending, running,
	// success, failed, canceled, skipped, manual, scheduled
	switch pipeline.Status {
	case "success":
		return true, nil
	case "failed":
		return true, nil
	case "canceled":
		return true, nil
	case "skipped":
		return true, nil
	default:
		return false, nil
	}
}

