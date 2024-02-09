package goutils

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"os/exec"
	"sort"
	"strings"
)

var Version string = "0.2.0"

func Check(err error) {
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

func CreateHttpClientWrapper(proxy string) (client *http.Client) {
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

func HttpRequestWrapper(client *http.Client, method string, url string, data string, headers []string) (respLine []byte, respCode int, err error) {
	request, err := http.NewRequest(method, url, bytes.NewBuffer([]byte(data)))
	if err != nil {
		return nil, 0, err
	}

	for _, header := range headers {
		arr := strings.Split(header, ":")
		request.Header.Set(arr[0], arr[1])
	}

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

func ExecuteShellScript(scriptString string, args ...string) error {
	shellArgs := []string{"-s", "-"}
	shellArgs = append(shellArgs, args...)

	cmd := exec.Command("sh", shellArgs...)
	cmd.Stdin = strings.NewReader(scriptString)

	var stdBuffer bytes.Buffer
	mw := io.MultiWriter(os.Stdout, &stdBuffer)
	cmd.Stdout = mw
	cmd.Stderr = mw

	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

var lastPrintrStrLen int = 0

func Printr1(str string) {
	if str[len(str)-1:] == "\r" {
		str = str[:len(str)-1]
	}
	lastPrintrStrLen = len(str)
	fmt.Print(str)
}

func Printr2(str string) {
	fmt.Print("\r" + strings.Repeat(" ", lastPrintrStrLen) + "\r")
	fmt.Print(str)
}
