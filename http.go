package goutils

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"math/rand"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

func CreateClient(proxy string) (client *http.Client) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	if proxy != "" {
		var proxyUrl, _ = url.Parse(proxy)
		tr.Proxy = http.ProxyURL(proxyUrl)
	}
	jar, _ := cookiejar.New(nil)
	client = &http.Client{
		Transport: tr,
		Jar:       jar,
	}
	return client
}

func Request(client *http.Client, method string, url string, data string, headers []string) (respLine []byte, respCode int, err error) {
	request, err := http.NewRequest(method, url, bytes.NewBuffer([]byte(data)))
	if err != nil {
		return nil, 0, err
	}

	for _, header := range headers {
		arr := strings.Split(header, ":")
		request.Header.Set(arr[0], arr[1])
	}

	request.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:122.0) Gecko/20100101 Firefox/122.0")
	request.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
	request.Header.Set("Accept-Language", "ru-RU,ru;q=0.8,en-US;q=0.5,en;q=0.3")
	request.Header.Set("Accept-Encoding", "gzip, deflate, br")

	response, err := client.Do(request)
	if err != nil {
		return nil, 0, err
	}

	respBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, 0, err
	}

	respCode = response.StatusCode
	return respBytes, respCode, err
}

func RequestMultipart(client *http.Client, method string, url string, data []Form, headers []string) (respLine []byte, respCode int, err error) {
	body := new(bytes.Buffer)
	mw := multipart.NewWriter(body)
	boundary := "-----------------------------"

	for i := 0; i < 28; i++ {
		boundary = boundary + fmt.Sprint(rand.Intn(10))
	}
	mw.SetBoundary(boundary)

	for _, form := range data {
		form.CreateForm(mw)
	}

	mw.Close()

	request, err := http.NewRequest(method, url, body)

	if err != nil {
		return nil, 0, err
	}

	for _, header := range headers {
		arr := strings.Split(header, ":")
		request.Header.Set(arr[0], arr[1])
	}
	request.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:122.0) Gecko/20100101 Firefox/122.0")
	request.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
	request.Header.Set("Accept-Language", "ru-RU,ru;q=0.8,en-US;q=0.5,en;q=0.3")
	request.Header.Set("Accept-Encoding", "gzip, deflate, br")
	request.Header.Set("Content-Type", "multipart/form-data; boundary="+boundary)

	response, err := client.Do(request)
	if err != nil {
		return nil, 0, err
	}

	respBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, 0, err
	}

	respCode = response.StatusCode
	return respBytes, respCode, err
}
