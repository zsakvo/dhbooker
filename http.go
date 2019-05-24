package main

import (
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// var quit = make(chan int)
// var downloadSuccess (chan string)
// var downloadFailed (chan string)
// var downloadChan (chan int) //下载索引

func getBody(res *http.Response) (string, error) {
	resBody, err := ioutil.ReadAll(res.Body)
	return string(resBody), err
}

func httpGet(url string, paramsMap map[string]string) (*http.Response, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := http.Client{
		Transport: tr,
		Timeout:   time.Duration(ping * int64(time.Millisecond)),
	}
	request, err := http.NewRequest("GET", url, nil)
	check(err)
	request.Header.Set("User-Agent", "dhbooker")
	params := request.URL.Query()
	if paramsMap != nil {
		for m, n := range paramsMap {
			params.Add(m, n)
		}
		request.URL.RawQuery = params.Encode()
	}
	return client.Do(request)
}

func httpPost(url string, content string) string {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := http.Client{
		Transport: tr,
		Timeout:   time.Duration(ping * int64(time.Millisecond)),
	}
	request, err := http.NewRequest("POST", url, strings.NewReader(content))
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Add("User-Agent", "dhbooker")
	response, err := client.Do(request)
	check(err)
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	check(err)
	return string(body)
}
