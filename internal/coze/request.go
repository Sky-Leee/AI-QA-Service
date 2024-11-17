package coze

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func sendRequest(method, api, token, host string, body any) ([]byte, error) {
	// 将结构体转换为 JSON 数据
	jsonData, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	// 创建 HTTP 请求
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, method, api, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	// 设置请求头
	token = fmt.Sprintf("Bearer %v", token)
	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Host", host)
	req.Header.Set("Connection", "keep-alive")

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 获取响应结果
	if resp.Status != "200 OK" {
		return nil, fmt.Errorf("response from %v is %v", api, resp.Status)
	}
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return respBody, nil
}

func SendCozeRequest(args *RequestBody, reply *CozeResponseBody) error {
	respBody, err := sendRequest(args.Method, args.URL, args.Token, args.Host, args.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(respBody, reply)
}
