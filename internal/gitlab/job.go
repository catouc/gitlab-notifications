package gitlab

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type Job struct {
	Commit struct {
		AuthorEmail string    `json:"author_email"`
		AuthorName  string    `json:"author_name"`
		CreatedAt   time.Time `json:"created_at"`
		ID          string    `json:"id"`
		Message     string    `json:"message"`
		ShortID     string    `json:"short_id"`
		Title       string    `json:"title"`
	} `json:"commit"`
	Coverage          interface{} `json:"coverage"`
	Archived          bool        `json:"archived"`
	AllowFailure      bool        `json:"allow_failure"`
	CreatedAt         time.Time   `json:"created_at"`
	StartedAt         time.Time   `json:"started_at"`
	FinishedAt        time.Time   `json:"finished_at"`
	ErasedAt          interface{} `json:"erased_at"`
	Duration          float64     `json:"duration"`
	QueuedDuration    float64     `json:"queued_duration"`
	ArtifactsExpireAt time.Time   `json:"artifacts_expire_at"`
	TagList           []string    `json:"tag_list"`
	ID                int         `json:"id"`
	Name              string      `json:"name"`
	Pipeline          struct {
		ID        int    `json:"id"`
		ProjectID int    `json:"project_id"`
		Ref       string `json:"ref"`
		Sha       string `json:"sha"`
		Status    string `json:"status"`
	} `json:"pipeline"`
	Ref       string        `json:"ref"`
	Artifacts []interface{} `json:"artifacts"`
	Runner    interface{}   `json:"runner"`
	Stage     string        `json:"stage"`
	Status    string        `json:"status"`
	Tag       bool          `json:"tag"`
	WebURL    string        `json:"web_url"`
	Project   struct {
		CiJobTokenScopeEnabled bool `json:"ci_job_token_scope_enabled"`
	} `json:"project"`
	User struct {
		ID           int         `json:"id"`
		Name         string      `json:"name"`
		Username     string      `json:"username"`
		State        string      `json:"state"`
		AvatarURL    string      `json:"avatar_url"`
		WebURL       string      `json:"web_url"`
		CreatedAt    time.Time   `json:"created_at"`
		Bio          interface{} `json:"bio"`
		Location     interface{} `json:"location"`
		PublicEmail  string      `json:"public_email"`
		Skype        string      `json:"skype"`
		Linkedin     string      `json:"linkedin"`
		Twitter      string      `json:"twitter"`
		WebsiteURL   string      `json:"website_url"`
		Organization string      `json:"organization"`
	} `json:"user"`
}

func (c *Client) GetJob(ctx context.Context, jobURL string) (*Job, error) {
	match := c.JobRegexp.FindStringSubmatch(jobURL)
	jobData := make(map[string]string)

	for i, value := range match {
		jobData[c.JobRegexp.SubexpNames()[i]] = value
	}

	projectID := url.PathEscape(jobData["projectID"])

	mrID, err := strconv.Atoi(jobData["jobID"])
	if err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf("projects/%s/jobs/%d", projectID, mrID)
	resp, err := c.callAPI(ctx, http.MethodGet, endpoint, nil)

	var returnData Job
	err = c.transformHTTPBodyIntoType(resp.Body, &returnData)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return &returnData, nil
}
