package gitlab

import (
	"context"
	"encoding/json"
	"io"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"time"
)

const gitlabAPIBasePath = "/api/v4"

// GET /projects/:id/merge_requests/:merge_request_iid/approvals
type Client struct {
	BaseURL    string
	Token      string
	HTTPClient *http.Client

	MRRegexp *regexp.Regexp
	JobRegexp *regexp.Regexp
	PipelineRegexp *regexp.Regexp
}

func New(baseURL, token string, timeout time.Duration) *Client {
	httpClient := http.DefaultClient
	httpClient.Timeout = timeout

	c := Client{
		BaseURL:    baseURL,
		Token:      token,
		HTTPClient: httpClient,
		MRRegexp: regexp.MustCompile("^" + baseURL + "/(?P<projectID>.+)/-/merge_requests/(?P<mrID>[1-9][0-9]*)$"),
		JobRegexp: regexp.MustCompile("^" + baseURL + "/(?P<projectID>.+)/-/jobs/(?P<jobID>[1-9][0-9]*)$"),
		PipelineRegexp: regexp.MustCompile("^" + baseURL + "/(?P<projectID>.+)/-/pipelines/(?P<pipelineID>[1-9][0-9]*)$"),
	}

	return &c
}

// gets in a URL of either job, pipeline or MR
// if MR -> check for approvals left
// if pipeline -> check for finished or not
// if job -> check if finished or not
// TODO: make this an enum, actually just rewrite this stuff in
// Rust to have proper typing
func (c *Client) IdentifyGitlabResourceFromURL(url string) string {
	switch {
	case c.MRRegexp.Match([]byte(url)):
		return "mr"
	case c.JobRegexp.Match([]byte(url)):
		return "job"
	case c.PipelineRegexp.Match([]byte(url)):
		return "pipeline"
	default:
		return "invalid"
	}	
}

func (c *Client) callAPI(ctx context.Context, method string, endpoint string, body io.Reader) (*http.Response, error) {
	requestURL, err := url.JoinPath(c.BaseURL, gitlabAPIBasePath, endpoint)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, method, requestURL, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("PRIVATE-TOKEN", c.Token)

	return c.HTTPClient.Do(req)
}

func (c *Client) transformHTTPBodyIntoType(body io.Reader, target interface{}) error {
	bodyBytes, err := io.ReadAll(body)
	if err != nil {
		return fmt.Errorf("failed to read bytes from HTTP response body: %s", err)
	}

	err = json.Unmarshal(bodyBytes, &target)
	if err != nil {
		return fmt.Errorf("failed to unmarshal HTTP body bytes: %s", err)
	}

	return nil
}
