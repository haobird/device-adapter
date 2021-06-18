package faceguard

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func between(str string, start string, end string) string {
	// Get substring between two strings.
	s := strings.Index(str, start)
	if s == -1 {
		return ""
	}
	s += len(start)
	e := strings.Index(str[s:], end)
	if e == -1 {
		return ""
	}
	return str[s : e+s]
}

//Request 发起 Http请求
func Request(url string, method string, data []byte, headers map[string]string) (result string, err error) {
	client := &http.Client{Timeout: 10 * time.Second}

	// jsonStr, _ := json.Marshal(data)
	req, err := http.NewRequest(method, url, bytes.NewBuffer(data))
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	if err != nil {
		return "", err
	}
	defer req.Body.Close()

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	bodystr := string(body)
	return bodystr, nil

	// if resp.StatusCode == 200 {
	// 	body, err_ := ioutil.ReadAll(resp.Body)
	// 	if err_ != nil {
	// 		return "", err_
	// 	}
	// 	bodystr := string(body)
	// 	return bodystr, nil
	// }
	// return "", err

}
