package cmd

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"time"

	"git.woa.com/kf_cdms/go-public/helper/carr"
	"github.com/ghodss/yaml"
	"github.com/spf13/cast"
)

type confGet interface {
	Get(ctx context.Context, name string) (string, error)
}

type cmd struct {
	ctx       context.Context
	name      string
	ipPort    string
	timeOut   time.Duration
	input     []interface{}
	command   string
	errMsg    error
	confProxy confGet
}

type cmdResult struct {
	value   string // tcp读取的原始响应
	result  int
	content map[string]interface{}
	errMsg  error
}

// NewCommand 初始化连接
// ctx上下文信息，目前仅仅用于获取调用人rtx和task_id
// conf 配置获取抽象接口
// name interface配置的名称
// params 传入参数（用于绑定input，当不需要input时，这个值直接穿nil）
func NewCommand(ctx context.Context, conf confGet, name string, params map[string]string) *cmd {
	cm := new(cmd)
	cm.ctx = ctx
	// 1.获取当前接口的配置
	interStrVal, err := conf.Get(ctx, name)
	if err != nil {
		cm.errMsg = err
		return cm
	}
	var interVal map[string]interface{}
	err = yaml.Unmarshal([]byte(interStrVal), &interVal)
	if err != nil {
		cm.errMsg = err
		return cm
	}
	// 2.组装请求
	cm.name = name
	cm.ipPort = cast.ToString(interVal["ip_port"])
	cm.timeOut = time.Second * time.Duration(cast.ToInt(interVal["timeout"]))
	cm.input = cast.ToSlice(interVal["input"])
	cm.makeCommand(cast.ToString(interVal["cmd"]), params)
	// 拼接额外参数
	var (
		fromUserId = params["__from_userId"]
		fromType   = params["__from_type"]
		taskId     = params["__task_id"]
		ticket     = params["__ticket"]
	)
	cm.command += fmt.Sprintf("&fromtype=%s&__task_id=%s&fromuserid=%s&__ticket=%s", fromType, taskId, fromUserId, ticket)
	return cm
}

// 生成command串
func (s *cmd) makeCommand(command string, params map[string]string) {
	s.command = command
	// 出现错误的，直接返回空
	if s.errMsg != nil {
		return
	}
	// 无额外参数时，也直接返回空
	if len(s.input) <= 0 {
		return
	}
	// 检测传入参数和额外参数的匹配
	// 首先检查必传参数是否都传入
	tmpVal := ""
	for _, v := range s.input {
		mv := cast.ToStringMap(v)
		keyName := cast.ToString(mv["key"])
		inputName := cast.ToString(mv["input"])
		// 当参数为必传参数时，若外部没给到该参数，直接报错
		if cast.ToBool(mv["detection"]) {
			if val, ok := params[keyName]; !ok || val == "" {
				s.errMsg = errors.New(keyName + "不能为空")
				return
			}
		}
		// 正常拼接参数，只接受Input中定义的参数，额外参数直接丢弃
		if parValue, ok := params[keyName]; !ok {
			tmpVal = tmpVal + "&" + inputName + "="
		} else {
			tmpVal = tmpVal + "&" + inputName + "=" + parValue
		}
	}
	s.command += tmpVal
	return
}

// Send 发送
// param ip []string 根据提供的ip列表随机选择ip进行请求（不会check）
func (s *cmd) Send(ips ...string) (*cmdResult, error) {
	cRes := &cmdResult{}
	if s.errMsg != nil {
		return cRes, s.errMsg
	}
	var hosts []string
	if len(ips) > 0 { // 指定了ip列表，则不进行check检测
		hosts = ips
	} else { // 检查host
		hosts = CheckHostPort(s.ipPort)
		if len(hosts) == 0 {
			return cRes, errors.New("check connection fail")
		}
	}
	// 随机取host请求
	host := carr.NewArr(hosts).RandSlice()
	// 记录接口请求时间
	res, err := InterfaceSend(host, s.command, s.timeOut)
	if err != nil {
		return cRes, err
	}
	cRes.value = res
	return cRes, nil
}

// UrlDecode url解码
func (r *cmdResult) UrlDecode() *cmdResult {
	// QueryEscape会把空格转成+，PathEscape会把空格转成%20 PathUnescape不会把+号解析为空格
	tmp, err := url.QueryUnescape(r.value)
	if err == nil {
		r.value = tmp
	}
	return r
}

// StrValue 返回原始返回串
func (r *cmdResult) StrValue() string {
	return r.value
}

// MapValue 格式化结果为map
func (r *cmdResult) MapValue(decode bool, splitNum int) map[string]interface{} {
	return ParseUrlParams(r.value, decode, splitNum)
}
