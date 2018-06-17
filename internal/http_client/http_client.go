package world_cup_http_client

import (
	"io/ioutil"
	"net/http"
	"net/url"
)

var (
	defaultBaseURL = "https://world-cup-json.herokuapp.com"
)

type HTTPClient struct {
	BaseURL   *url.URL
	ClientImp *http.Client
}

type Response struct {
	Body []byte
}

func New(baseURL string) (*HTTPClient, error) {
	resolvedBaseURL := defaultBaseURL
	if baseURL != "" {
		resolvedBaseURL = baseURL
	}
	parsedURL, err := url.Parse(resolvedBaseURL)
	if err != nil {
		return nil, err
	}
	return &HTTPClient{
		BaseURL:   parsedURL,
		ClientImp: http.DefaultClient, // Maybe support different client in the future
	}, nil
}

func (httpClient *HTTPClient) Get(path string, queryParams map[string]string) (*Response, error) {
	url := httpClient.prepareURL(path, queryParams)
	resp, err := httpClient.ClientImp.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	return &Response{Body: body}, nil
}

func (httpClient *HTTPClient) prepareURL(path string, queryParams map[string]string) string {
	url := *httpClient.BaseURL // dereference for copying
	url.Path = path
	values := url.Query()
	for key, value := range queryParams {
		values.Add(key, value)
	}
	url.RawQuery = values.Encode()
	return url.String()
}
