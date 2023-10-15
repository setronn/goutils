package utils

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

func unique(arr []string) (ans []string) {
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

func readFile(filename string) (lines []string) {
	file, _ := os.ReadFile(filename)

	sc := bufio.NewScanner(strings.NewReader(string(file)))
	for sc.Scan() {
		lines = append(lines, sc.Text())
	}
	return lines
}

func writeFile(filename string, lines []string) {
	file, _ := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	datawriter := bufio.NewWriter(file)
	for _, data := range lines {
		_, _ = datawriter.WriteString(data + "\n")
	}
	datawriter.Flush()
	file.Close()
}

func create_http_client(proxy string) (client *http.Client) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	if proxy != "" {
		var proxyUrl, _ = url.Parse("http://127.0.0.1:8080")
		tr.Proxy = http.ProxyURL(proxyUrl)
	}

	jar, _ := cookiejar.New(nil)
	client = &http.Client{
		Transport: tr,
		Jar:       jar,
	}

	return client
}

func http_request_wrapper(client *http.Client, method string, url string, data io.Reader, headers map[string]string) []byte {
	request, _ := http.NewRequest(method, url, data)
	for key, value := range headers {
		request.Header.Set(key, value)
	}
	response, _ := client.Do(request)
	body, _ := io.ReadAll(response.Body)
	response.Body.Close()
	return body
}
