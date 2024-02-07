package goutils

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"sort"
	"strings"
)

var Version string = "0.1.4"

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func SortUnique(arr []string) (ans []string) {
	var sortedArr []string
	sortedArr = append(sortedArr, arr...)
	sort.Strings(sortedArr)

	if len(sortedArr) == 0 {
		return ans
	}
	for i := 0; i < len(sortedArr)-1; i++ {
		if sortedArr[i] == sortedArr[i+1] {
			continue
		}
		ans = append(ans, sortedArr[i])
	}
	ans = append(ans, sortedArr[len(sortedArr)-1])
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
	file, _ := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
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

func HttpRequestWrapper(client *http.Client, method string, url string, data string, headers []string) (respLine string, respCode int) {
	request, _ := http.NewRequest(method, url, bytes.NewBuffer([]byte(data)))
	for _, header := range headers {
		arr := strings.Split(header, ":")
		request.Header.Set(arr[0], arr[1]) // There is an additional space after :
	}
	response, err := client.Do(request)
	check(err)
	respBytes, err := io.ReadAll(response.Body)
	check(err)
	response.Body.Close()
	respCode = response.StatusCode
	respLine = string(respBytes[:])
	return respLine, respCode
}
