package robot

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
)

// 企业微信robot

type qw struct {
	url string
	rtx []string
}

var gbRobotUrl = "" // 机器人地址
var gbNoticeRtx = make([]string, 0)

// GlobalInit 全局设置默认机器人地址 以及 通知rtx
func GlobalInit(robotUrl string, noticeRtx []string) {
	gbRobotUrl = robotUrl
	gbNoticeRtx = noticeRtx
}

func NewQw() *qw {
	return &qw{
		url: gbRobotUrl,
		rtx: gbNoticeRtx,
	}
}

func (q *qw) SetUrl(url string) *qw {
	q.url = url
	return q
}

func (q *qw) SetNoticeRtx(rtx []string) *qw {
	q.rtx = rtx
	return q
}

// SendText 企业微信发文本消息
func (q *qw) SendText(ctx context.Context, msg string) ([]byte, error) {
	params := map[string]interface{}{
		"msgtype": "text",
		"text":    map[string]string{"content": msg},
	}
	return post(ctx, q.url, params)
}

func post(ctx context.Context, url string, params map[string]interface{}) ([]byte, error) {
	paramData, _ := json.Marshal(params)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(paramData))
	req.Header.Set("content-Type", "application/json")
	if err != nil {
		return []byte{}, err
	}
	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()
	// 读取响应
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}
	return respBody, nil
}
