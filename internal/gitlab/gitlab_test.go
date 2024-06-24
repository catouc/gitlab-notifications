package gitlab

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestIdentifyGitlabResourceFromURL(t *testing.T) {
	testData := []struct{
		Name string
		In string
		Out string
	}{
		{
			"HappyPipelineSingleDigit",
			"https://gitlab.com/project/-/pipelines/1",
			"pipeline",
		},
		{
			"HappyPipelineMultipleDigit",
			"https://gitlab.com/project/-/pipelines/111",
			"pipeline",
		},
		{
			"HappyMergeRequestSingleDigit",
			"https://gitlab.com/project/-/merge_requests/1",
			"mr",
		},
		{
			"HappyMergeRequestsMultipleDigit",
			"https://gitlab.com/project/-/merge_requests/111",
			"mr",
		},
		{
			"HappyJobsSingleDigit",
			"https://gitlab.com/project/-/jobs/1",
			"job",
		},
		{
			"HappyJobsMultipleDigit",
			"https://gitlab.com/project/-/jobs/111",
			"job",
		},
	}

	for _, td := range testData {
		t.Run(td.Name, func(t *testing.T) {
			client := New("https://gitlab.com", "", time.Second)
			result := client.IdentifyGitlabResourceFromURL(td.In)
			assert.Equal(t, td.Out, result)
		})
	}
}
