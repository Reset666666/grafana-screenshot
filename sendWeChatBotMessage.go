package main

import (
	"bytes"
	"fmt"
	"grafana-screenshot/logs"
	"io"
	"net/http"
)

// 企业微信群机器人推送（独立版本）
// key: 企业微信机器人 key
// msg: 支持 markdown
func SendWeChatBotMessage(key string, msg string) error {
	if key == "" {
		return fmt.Errorf("企业微信机器人 key 未配置")
	}

	url := "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=" + key

	payload := fmt.Sprintf(`{
		"msgtype": "markdown",
		"markdown": {
			"content": "%s"
		}
	}`, msg)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer([]byte(payload)))
	if err != nil {
		return fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	logs.InfoLogger.Println("🤖 企业微信返回:", string(body))
	return nil
}
