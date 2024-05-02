package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	urllib "net/url"
	"strings"
)

func setHeader(req *http.Request, mobile bool, referer string, cookies ...*http.Cookie) {
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}
	if referer != "" {
		req.Header.Add("Referer", referer)
	}
	if mobile {
		req.Header.Add("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 14_2 like Mac OS X) AppleWebKit/604.1.28 (KHTML, like Gecko) CriOS/111.0.5563.8 Mobile/14E5239e Safari/602.1")
	} else {
		req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3")
	}
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Add("Accept-Language", "en-US,en;q=0.9")
}
func RequestPostForm(url string, body map[string]string, mobile bool) ([]byte, []*http.Cookie, error) {
	fmt.Println("RequestPostForm:", url)
	// fmt.Println("Start---POST Form---")
	// for k, v := range form {
	// 	fmt.Println(k, " : ", v)
	// }
	// fmt.Println("End---POST Form---")
	formData := urllib.Values{}
	for k, v := range body {
		formData.Add(k, v)
	}
	// 创建一个新的请求
	req, err := http.NewRequest("POST", url, strings.NewReader(formData.Encode()))
	if err != nil {
		return nil, nil, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	// 添加自定义的header
	setHeader(req, mobile, "")
	// 执行请求
	client := &http.Client{}
	resp, err := client.Do(req)
	cookies := resp.Cookies()
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()
	// 处理响应
	res, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}
	return res, cookies, nil
}

func RequestPost(url string, body interface{}, mobile bool) ([]byte, error) {
	fmt.Println("RequestPost:", url)
	// 创建一个新的请求
	json, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	reader := strings.NewReader(string(json))
	req, err := http.NewRequest("POST", url, reader)
	if err != nil {
		return nil, err
	}
	// 添加自定义的header
	setHeader(req, mobile, "")
	// 执行请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	// 处理响应
	res, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return res, nil
}
func RequestGet(url string, mobile bool) ([]byte, error) {
	return RequestGetReferer(url, mobile, "")
}
func RequestGetReferer(url string, mobile bool, referer string, cookies ...*http.Cookie) ([]byte, error) {
	fmt.Println("RequestGet:", url)
	// 创建一个新的请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	// 添加自定义的header
	setHeader(req, mobile, referer, cookies...)
	// 创建一个HTTP客户端并执行请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	// 处理响应
	res, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return res, nil
}
