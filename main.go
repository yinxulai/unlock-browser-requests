package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
)

type request struct {
	Url    string              `json:"url"`
	Method string              `json:"method"`
	Body   *[]byte             `json:"body"`
	Header map[string][]string `json:"header"`
}

type response struct {
	Status int                 `json:"status"`
	Body   *[]byte             `json:"body"`
	Header map[string][]string `json:"header"`
}

func requestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		return
	}

	var requestParams request
	err := json.NewDecoder(r.Body).Decode(&requestParams)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	requestParamUrl, err := url.Parse(requestParams.Url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	targetRequest := &http.Request{}
	targetRequest.URL = requestParamUrl
	targetRequest.Method = requestParams.Method
	targetRequest.Header = requestParams.Header
	if requestParams.Body != nil {
		readerCloser := io.NopCloser(bytes.NewReader(*requestParams.Body))
		targetRequest.Body = readerCloser
	}

	client := &http.Client{}
	targetResponse, err := client.Do(targetRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	var responseResult response
	responseResult.Header = targetResponse.Header
	responseResult.Status = targetResponse.StatusCode

	bytes, err := io.ReadAll(targetResponse.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	responseResult.Body = &bytes

	responseResultBytes, err := json.Marshal(responseResult)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	if _, err := w.Write(responseResultBytes); err != nil {
		log.Println("Failed to write response:", err)
	}
}

func main() {
	// 注册请求处理函数
	http.HandleFunc("/", requestHandler)

	// 启动HTTP服务，监听在端口8080
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		log.Fatal("Failed to start the server:", err)
	}
}
