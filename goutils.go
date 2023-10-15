package goutils

import (
	"bufio"
	"crypto/tls"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"
)

var Version string = "0.1.3"

func Unique(arr []string) (ans []string) {
	if len(arr) == 0 {
		return ans
	}
	for i := 0; i < len(arr)-1; i++ {
		if arr[i] == arr[i+1] {
			continue
		}
		ans = append(ans, arr[i])
	}
	ans = append(ans, arr[len(arr)-1])
	return ans
}

func ReadFile(filename string) (lines []string) {
	file, _ := os.ReadFile(filename)

	datareader := bufio.NewScanner(strings.NewReader(string(file)))
	for datareader.Scan() {
		lines = append(lines, datareader.Text())
	}
	return lines
}

func WriteFile(filename string, lines []string) {
	file, _ := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 755)
	datawriter := bufio.NewWriter(file)
	for _, data := range lines {
		_, _ = datawriter.WriteString(data + "\n")
	}
	datawriter.Flush()
	file.Close()
}

func CreateHttpClient(proxy string) (client *http.Client) {
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

func HttpRequest(client *http.Client, method string, url string, data io.Reader, headers map[string]string) []byte {
	request, _ := http.NewRequest(method, url, data)
	for key, value := range headers {
		request.Header.Set(key, value)
	}
	response, _ := client.Do(request)
	body, _ := io.ReadAll(response.Body)
	response.Body.Close()
	return body
}
