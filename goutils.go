package goutils

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
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

// [TODO] Fix: doesn't return error if shell returns with code 2 (error)
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

func ReadBinFile(filename string) ([]byte, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	stats, err := file.Stat()
	if err != nil {
		return nil, err
	}

	var size int64 = stats.Size()
	bytes := make([]byte, size)

	bufr := bufio.NewReader(file)
	_, err = bufr.Read(bytes)
	if err != nil {
		return nil, err
	}

	return bytes, nil
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
