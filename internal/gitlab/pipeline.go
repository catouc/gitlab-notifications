package gitlab

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type Pipeline struct {
	ID         int         `json:"id"`
	Iid        int         `json:"iid"`
	ProjectID  int         `json:"project_id"`
	Name       string      `json:"name"`
	Status     string      `json:"status"`
	Ref        string      `json:"ref"`
	Sha        string      `json:"sha"`
	BeforeSha  string      `json:"before_sha"`
	Tag        bool        `json:"tag"`
	YamlErrors interface{} `json:"yaml_errors"`
	User       struct {
		Name      string `json:"name"`
		Username  string `json:"username"`
		ID        int    `json:"id"`
		State     string `json:"state"`
		AvatarURL string `json:"avatar_url"`
		WebURL    string `json:"web_url"`
	} `json:"user"`
	CreatedAt      time.Time   `json:"created_at"`
	UpdatedAt      time.Time   `json:"updated_at"`
	StartedAt      interface{} `json:"started_at"`
	FinishedAt     time.Time   `json:"finished_at"`
	CommittedAt    interface{} `json:"committed_at"`
	Duration       float64     `json:"duration"`
	QueuedDuration float64     `json:"queued_duration"`
	Coverage       string      `json:"coverage"`
	WebURL         string      `json:"web_url"`
}

func (c *Client) GetPipeline(ctx context.Context, pipelineURL string) (*Pipeline, error) {
	match := c.PipelineRegexp.FindStringSubmatch(pipelineURL)
	pipelineData := make(map[string]string)

	for i, value := range match {
		pipelineData[c.PipelineRegexp.SubexpNames()[i]] = value
	}

	projectID := url.PathEscape(pipelineData["projectID"])

	pipelineID, err := strconv.Atoi(pipelineData["pipelineID"])
	if err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf("projects/%s/pipelines/%d", projectID, pipelineID) 
	resp, err := c.callAPI(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	var returnData Pipeline
	err = c.transformHTTPBodyIntoType(resp.Body, &returnData)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return &returnData, nil
}
