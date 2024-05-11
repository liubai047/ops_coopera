package delayQueue

import "errors"

var RepeatError = errors.New("message key has repeat") // 消息key重复
var KeyError = errors.New("invalid key")               // 无效的消息key
var CtxDoneError = errors.New("context done quit")     // ctx退出
var EmptyQueue = errors.New("empty queue")             // 空队列
