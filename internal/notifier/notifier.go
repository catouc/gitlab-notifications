package notifier

import (
	"context"

	"github.com/catouc/gitlab-notifications/internal/gitlab"
)


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

