package main

import (
	"net/http"
	"net/url"
	"fmt"
)

//ShortlyClient
type ShortlyClient interface {
	Name() string
	Shorten(url string) (string, error)
	Expand(url string) (string, error)
}

// HttpShortlyClient connects to the running server over http
type HttpShortlyClient struct {
	name    string
	baseUrl string
	http    *http.Client
}

func MakeShortlyClient(name string, baseUrl string) ShortlyClient {
	client := HttpShortlyClient{}
	client.name = name
	client.baseUrl = baseUrl
	client.http = &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	return &client
}

//Name return this client's name.
func (client *HttpShortlyClient) Name() string {
	return client.name
}

func (client *HttpShortlyClient) Shorten(urlString string) (string, error) {
	resp, err := client.http.PostForm(client.baseUrl+"/shorten", url.Values{"url": {urlString}})
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 302 {
		return "", fmt.Errorf("Expected 302, got %d", resp.StatusCode)
	}
	location := resp.Header.Get("Location")
	slug := location[len("/status/"):]
	return slug, nil
}

func (client *HttpShortlyClient) Expand(slug string) (string, error) {
	resp, err := client.http.Get(client.baseUrl + "/" + slug)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 302 {
		return "", fmt.Errorf("Expected 302, got %d", resp.StatusCode)
	}
	location := resp.Header.Get("Location")
	return location, nil
}
