package gitlab

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type MergeRequestApprovals struct {
	ID                int        `json:"id"`
	Iid               int        `json:"iid"`
	ProjectID         int        `json:"project_id"`
	Title             string     `json:"title"`
	Description       string     `json:"description"`
	State             string     `json:"state"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	MergeStatus       string     `json:"merge_status"`
	ApprovalsRequired int        `json:"approvals_required"`
	ApprovalsLeft     int        `json:"approvals_left"`
	ApprovedBy        []Approver `json:"approved_by"`
}

type Approver struct {
	User User `json:"user"`
}

type User struct {
	Name      string `json:"name"`
	Username  string `json:"username"`
	ID        int    `json:"id"`
	State     string `json:"state"`
	AvatarURL string `json:"avatar_url"`
	WebURL    string `json:"web_url"`
}

func (c *Client) GetMergeRequestApprovals(ctx context.Context, mrURL string) (*MergeRequestApprovals, error) {
	match := c.MRRegexp.FindStringSubmatch(mrURL)
	mrData := make(map[string]string)

	for i, value := range match {
		mrData[c.MRRegexp.SubexpNames()[i]] = value
	}

	projectID := url.PathEscape(mrData["projectID"])

	mrID, err := strconv.Atoi(mrData["mrID"])
	if err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf("projects/%s/merge_requests/%d/approvals", projectID, mrID)
	resp, err := c.callAPI(ctx, http.MethodGet, endpoint, nil)

	var returnData MergeRequestApprovals
	err = c.transformHTTPBodyIntoType(resp.Body, &returnData)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return &returnData, nil
}
