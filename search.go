package harborcli

import (
	"fmt"

	"github.com/goharbor/harbor/src/common/models"
	"k8s.io/helm/cmd/helm/search"
)

type SearchAPI struct {
	client *HarborClient
}

type searchResult struct {
	Projects     []*models.Project        `json:"project"`
	Repositories []map[string]interface{} `json:"repository"`
	Chart        *[]*search.Result        `json:"chart,omitempty"`
}

// Search for projects, repositories and helm charts
func (s *SearchAPI) Search(name string) (*searchResult, error) {
	path := fmt.Sprintf("api/search?q=%s", name)
	req, err := s.client.newRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	results := &searchResult{}
	_, err = s.client.do(req, results)

	return results, err
}
