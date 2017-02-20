package main

import (
	"net/http"
	"net/url"
	"fmt"
	"io"
	"io/ioutil"
)

//ShortlyClient
type ShortlyClient interface {
	Shorten(url string) (string, error)
	Expand(url string) (string, error)
}

// HttpShortlyClient connects to the running server over http
type HttpShortlyClient struct {
	baseUrl string
	http    *http.Client
}

func MakeShortlyClient(baseUrl string) ShortlyClient {
	client := HttpShortlyClient{}
	client.baseUrl = baseUrl
	http.DefaultTransport.(*http.Transport).MaxIdleConnsPerHost = 500
	http.DefaultTransport.(*http.Transport).MaxIdleConns = 500
	client.http = &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Transport: http.DefaultTransport,
	}
	return &client
}

func (client *HttpShortlyClient) Shorten(urlString string) (string, error) {
	resp, err := client.http.PostForm(client.baseUrl+"/shorten", url.Values{"url": {urlString}})
	if err != nil {
		return "", err
	}
	defer close(resp)
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
	defer close(resp)
	if resp.StatusCode != 302 {
		return "", fmt.Errorf("Expected 302, got %d", resp.StatusCode)
	}
	location := resp.Header.Get("Location")
	return location, nil
}

func close(res *http.Response) {
	io.Copy(ioutil.Discard, res.Body)
	res.Body.Close()
}

