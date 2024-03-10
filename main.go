package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

const (
	// Use in Response
	OverwriteRequestUrl           = "overwrite-request-url"      // example: overwrite-request-url: https://example.com
	ExposeResponseHeader          = "expose-response-header"     // example: expose-response-header: set-cookie
	OverwriteRequestHeaderPrefix  = "overwrite-request-header-"  // example: overwrite-request-header:cookie: example-cookie=example
	OverwriteResponseHeaderPrefix = "overwrite-response-header-" // example: overwrite-response-header:Access-Control-Allow-Origin: *

	// Use in Response
	ExposedResponseHeader = "exposed-header-" // 用于在 response 中暴露 ExposeResponseHeader 指定的 header
)

type AgentOptions struct {
	OverwriteRequestURL     string
	ExposeResponseHeader    []string
	OverwriteRequestHeader  map[string][]string
	OverwriteResponseHeader map[string][]string
}

func ParseOptionsFormHeader(header *http.Header) *AgentOptions {
	var agentOptions AgentOptions

	var supportedOptionPrefixList = []string{
		OverwriteRequestUrl,
		ExposeResponseHeader,
		OverwriteRequestHeaderPrefix,
		OverwriteResponseHeaderPrefix,
	}

	for headerKey, headerValues := range *header {
		lowerHeaderKey := strings.ToLower(headerKey)
		for _, supportedOptionPrefix := range supportedOptionPrefixList {
			if strings.HasPrefix(lowerHeaderKey, supportedOptionPrefix) {
				// OverwriteRequestURL
				if supportedOptionPrefix == OverwriteRequestUrl {
					agentOptions.OverwriteRequestURL = headerValues[0]
				}

				// ExposeResponseHeader
				if supportedOptionPrefix == ExposeResponseHeader {
					agentOptions.ExposeResponseHeader = headerValues
				}

				// OverwriteRequestHeader
				if supportedOptionPrefix == OverwriteRequestHeaderPrefix {
					if agentOptions.OverwriteRequestHeader == nil {
						agentOptions.OverwriteRequestHeader = make(map[string][]string)
					}

					overwriteKey, _ := strings.CutPrefix(lowerHeaderKey, OverwriteRequestHeaderPrefix)
					agentOptions.OverwriteRequestHeader[overwriteKey] = headerValues
				}

				// OverwriteResponseHeader
				if supportedOptionPrefix == OverwriteResponseHeaderPrefix {
					if agentOptions.OverwriteResponseHeader == nil {
						agentOptions.OverwriteResponseHeader = make(map[string][]string)
					}

					overwriteKey, _ := strings.CutPrefix(lowerHeaderKey, OverwriteResponseHeaderPrefix)
					agentOptions.OverwriteResponseHeader[overwriteKey] = headerValues
				}

				// 清理请求的 header 便于后面使用
				header.Del(lowerHeaderKey)
			}
		}
	}

	return &agentOptions
}

func AutoProxy(w http.ResponseWriter, request *http.Request) {
	// 设置允许跨域访问的源
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")

	// 跨域预检请求
	if request.Method == "OPTIONS" {
		return
	}

	agentOptions := ParseOptionsFormHeader(&request.Header)
	log.Printf("Received a request: %v", agentOptions)

	// 处理 OverwriteRequestURL
	if agentOptions.OverwriteRequestURL == "" {
		err := fmt.Errorf("OverwriteRequestURL cannot be empty")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("Received a request: %s %s, options: %v", request.Method, request.URL.String(), agentOptions)
	agentRequest, err := http.NewRequestWithContext(request.Context(), request.Method, agentOptions.OverwriteRequestURL, request.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Print(err)
		return
	}

	// 转发全部的 header
	for requestHeaderKey, requestHeaderValues := range request.Header {
		for _, requestHeaderValue := range requestHeaderValues {
			agentRequest.Header.Add(requestHeaderKey, requestHeaderValue)
		}
	}

	// 处理 OverwriteRequestHeader
	for overwriteKey, overwriteValues := range agentOptions.OverwriteRequestHeader {
		agentRequest.Header.Del(overwriteKey)
		for _, overwriteValue := range overwriteValues {
			agentRequest.Header.Add(overwriteKey, overwriteValue)
		}
	}

	// 发送请求
	client := &http.Client{}
	agentResponse, err := client.Do(agentRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Print(err)
		return
	}

	// 处理 ExposeResponseHeader
	for _, exposeKey := range agentOptions.ExposeResponseHeader {
		exposedKey := ExposedResponseHeader + exposeKey
		exposedValues := agentResponse.Header.Values(exposeKey)
		for _, exposedValue := range exposedValues {
			agentResponse.Header.Add(exposedKey, exposedValue)
		}
	}

	// 处理 OverwriteResponseHeader
	for overwriteKey, overwriteValues := range agentOptions.OverwriteResponseHeader {
		agentResponse.Header.Del(overwriteKey)
		for _, overwriteValue := range overwriteValues {
			agentResponse.Header.Add(overwriteKey, overwriteValue)
		}
	}

	// 读取响应体内容
	body, err := ioutil.ReadAll(agentResponse.Body)
	if err != nil {
		log.Printf("Failed to read response body:", err)
		return
	}

	// 设置响应头
	for key, values := range agentResponse.Header {
		for _, value := range values {
			// 将响应头写入到 ResponseWriter
			// 这里假设你已经有一个 http.ResponseWriter 对象，命名为 "w"
			w.Header().Add(key, value)
		}
	}

	// 设置响应状态码
	w.WriteHeader(agentResponse.StatusCode)

	// 将响应体写入到 ResponseWriter
	if _, err = w.Write(body); err != nil {
		log.Printf("Failed to read response body:", err)
	}
}

func main() {
	// 注册请求处理函数
	http.HandleFunc("/", AutoProxy)

	// 启动 HTTP 服务，监听在端口808
	log.Println("start the server on port 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Failed to start the server:", err)
	}
}
