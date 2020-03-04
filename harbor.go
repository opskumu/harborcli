package harborcli

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

const (
	AuthPingPath = "api/users/current"
	HealthPath   = "api/health"
	LoginPath    = "c/login"
)

type HarborError struct {
	Code   int
	method string
	url    *url.URL
}

func (e HarborError) Error() string {
	return fmt.Sprintf("%s %s returned %d", e.method, e.url, e.Code)
}

type LoginForm struct {
	Username string
	Password string
}

type HarborClient struct {
	BaseURL *url.URL
	Auth    LoginForm
	Client  *http.Client

	Project    *ProjectAPI
	Search     *SearchAPI
	Repository *RepositoryAPI
}

func NewHarborClient(harborURL string, auth LoginForm) (*HarborClient, error) {
	baseURL, err := url.Parse(harborURL)
	if err != nil {
		return nil, err
	}

	cookieJar, _ := cookiejar.New(nil)
	client := &HarborClient{
		BaseURL: baseURL,
		Auth:    auth,
		Client: &http.Client{
			Jar: cookieJar,
		},
	}

	client.Project = &ProjectAPI{
		client: client,
	}
	client.Search = &SearchAPI{
		client: client,
	}
	client.Repository = &RepositoryAPI{
		client: client,
	}

	return client, err
}

func (client *HarborClient) newRequest(method, urlStr string, body interface{}) (*http.Request, error) {
	rel, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	urlPath := client.BaseURL.ResolveReference(rel)

	buf := &bytes.Buffer{}
	if body != nil {
		err = json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, urlPath.String(), buf)
	if err != nil {
		return nil, err
	}

	req.Header.Add("User-Agent", "Mozilla/5.0 Gecko/20100101 Firefox/50.0")

	return req, nil
}

func (client *HarborClient) do(req *http.Request, v interface{}) (*http.Response, error) {
	// add X-Xsrftoken header
	for _, v := range client.Client.Jar.Cookies(client.BaseURL) {
		if v.Name == "_xsrf" {
			b64, _ := base64.StdEncoding.DecodeString(strings.Split(v.Value, "|")[0])
			req.Header.Add("X-Xsrftoken", string(b64))
			break
		}
	}

	resp, err := client.Client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	if status := resp.StatusCode; status < 200 || status > 299 {
		return resp, HarborError{Code: status, method: req.Method, url: req.URL}
	}

	if v != nil {
		err = json.NewDecoder(resp.Body).Decode(v)
	}

	return resp, err
}

func (client *HarborClient) healthAPIReq() error {
	req, _ := client.newRequest("GET", HealthPath, nil)
	_, err := client.do(req, nil)

	return err
}

func (client *HarborClient) Login() error {
	// get cookie with _xsrf
	err := client.healthAPIReq()
	if err != nil {
		return err
	}

	rel, err := url.Parse(LoginPath)
	if err != nil {
		return err
	}
	form := url.Values{
		"principal": {client.Auth.Username},
		"password":  {client.Auth.Password},
	}
	urlPath := client.BaseURL.ResolveReference(rel)
	body := bytes.NewBufferString(form.Encode())
	req, err := http.NewRequest("POST", urlPath.String(), body)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	_, err = client.do(req, nil)
	return err
}

func (client *HarborClient) authPing() error {
	req, err := client.newRequest("GET", AuthPingPath, nil)
	if err != nil {
		return err
	}

	resp, _ := client.do(req, nil)
	// If not auth, try login
	if statusCode := resp.StatusCode; statusCode == 401 {
		err := client.Login()
		if err != nil {
			return err
		}
	}

	return err
}
