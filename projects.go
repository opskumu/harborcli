package harborcli

import (
	"fmt"

	"github.com/goharbor/harbor/src/common/models"
)

const (
	ProjectAPIPath = "api/projects"
)

type ProjectAPI struct {
	client *HarborClient
}

func (p *ProjectAPI) Create(project *models.ProjectRequest) error {
	err := p.client.authPing()
	if err != nil {
		return err
	}
	req, err := p.client.newRequest("POST", ProjectAPIPath, project)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")

	_, err = p.client.do(req, nil)
	return err
}

// Check if the project name user provided already exists
func (p *ProjectAPI) Check(name string) error {
	err := p.client.authPing()
	if err != nil {
		return err
	}
	path := fmt.Sprintf("%s?project_name=%s", ProjectAPIPath, name)
	req, err := p.client.newRequest("HEAD", path, nil)
	if err != nil {
		return err
	}

	_, err = p.client.do(req, nil)
	return err
}

// Return specific project detail infomation
func (p *ProjectAPI) Get(id int64) (*models.Project, error) {
	err := p.client.authPing()
	if err != nil {
		return nil, err
	}
	path := fmt.Sprintf("%s/%d", ProjectAPIPath, id)
	req, err := p.client.newRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	project := &models.Project{}
	_, err = p.client.do(req, project)

	return project, err
}

// Update project by id
func (p *ProjectAPI) Update(id int64, project *models.ProjectRequest) error {
	err := p.client.authPing()
	if err != nil {
		return err
	}
	path := fmt.Sprintf("%s/%d", ProjectAPIPath, id)
	req, err := p.client.newRequest("PUT", path, project)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")

	_, err = p.client.do(req, nil)
	return err
}

// Delete project by id
func (p *ProjectAPI) Delete(id int64) error {
	err := p.client.authPing()
	if err != nil {
		return err
	}
	path := fmt.Sprintf("%s/%d", ProjectAPIPath, id)
	req, err := p.client.newRequest("DELETE", path, nil)
	if err != nil {
		return err
	}

	_, err = p.client.do(req, nil)
	return err
}

// List projects
func (p *ProjectAPI) List(name string) ([]*models.Project, error) {
	err := p.client.authPing()
	if err != nil {
		return nil, err
	}
	path := fmt.Sprintf("%s?name=%s", ProjectAPIPath, name)
	req, err := p.client.newRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var projects []*models.Project
	_, err = p.client.do(req, &projects)

	return projects, err
}
