package cmd

import (
	"errors"
	"fmt"
	"io"
	"net"
	"net/url"
	"strings"
	"time"
)

// InterfaceSend tcp客户端发送请求，并获取数据返回
func InterfaceSend(addr, cmd string, timeout time.Duration) (string, error) {
	// 创建请求链接
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return "", err
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return "", err
	}
	defer conn.Close()
	now := time.Now()
	// 发起请求
	err = conn.SetDeadline(now.Add(timeout))
	if err != nil {
		return "", err
	}
	_, err = conn.Write([]byte(cmd + "\r\n"))
	if err != nil && err.Error() != "EOF" {
		return "", err
	}
	// 接收响应
	_ = conn.SetReadDeadline(now.Add(timeout))
	b, err := readConn(conn, timeout, -1)
	if err != nil && err.Error() != "EOF" {
		return "", err
	}
	resStr := strings.Trim(string(b), "\r\n")
	return resStr, nil
}

// ParseUrlParams 解析url参数为map返回
// splitNum int 控制接口的解码数
// decode bool 控制接口是否需要解码
func ParseUrlParams(str string, decode bool, splitNum int) map[string]interface{} {
	strArr := strings.Split(str, "&")
	if splitNum > 0 {
		strArr = strings.SplitN(str, "&", 2)
	}
	data := map[string]interface{}{}
	for _, row := range strArr {
		values := strings.Split(row, "=")
		valueLen := len(values)
		if valueLen >= 2 {
			multStr := strings.Join(values[1:valueLen], "") // 兼容 k="v=v" 这种
			val := multStr
			if decode {
				tmp, err := url.QueryUnescape(multStr)
				if err == nil {
					val = tmp
				}
			}
			data[values[0]] = val
		} else if valueLen == 1 {
			data[values[0]] = ""
		}
	}
	return data
}

// CheckHostPort 检测端口是否可访问
func CheckHostPort(address string) []string {
	var adds []string
	addr := strings.Split(address, ";")
	for _, v := range addr {
		conn, err := net.DialTimeout("tcp", v, time.Millisecond*10)
		if err != nil {
			continue
		}
		conn.Close()
		adds = append(adds, v)
	}
	return adds
}

// readConn 读取tcp连接响应
// @param receiveTimeOut 接收超时时间
// @param length 数据长度，当<0时，表示数据长度不定,；当>0时，表示读取指定长度字段
func readConn(c net.Conn, receiveTimeOut time.Duration, length int) ([]byte, error) {
	receiveTime := time.Now().Add(receiveTimeOut)
	const defaultReadBufferSize int = 128
	var (
		err        error  // Reading error.
		size       int    // Reading size.
		index      int    // Received size.
		buffer     []byte // Buffer object.
		bufferWait bool   // Whether buffer reading timeout set.
	)
	if length > 0 {
		buffer = make([]byte, length)
	} else {
		buffer = make([]byte, defaultReadBufferSize)
	}
	for {
		if length < 0 && index > 0 {
			bufferWait = true
			if err = c.SetReadDeadline(time.Now().Add(time.Millisecond * 10)); err != nil {
				err = fmt.Errorf("SetReadDeadline for connection failed:%s", err.Error())
				return nil, err
			}
		}
		size, err = c.Read(buffer[index:])
		if size > 0 {
			index += size
			if length > 0 {
				// It reads til `length` size if `length` is specified.
				if index == length {
					break
				}
			} else {
				if index >= defaultReadBufferSize {
					// If it exceeds the buffer size, it then automatically increases its buffer size.
					buffer = append(buffer, make([]byte, defaultReadBufferSize)...)
				} else {
					// It returns immediately if received size is lesser than buffer size.
					if !bufferWait {
						break
					}
				}
			}
		}
		if err != nil {
			// Connection closed.
			if err == io.EOF {
				break
			}
			// Re-set the timeout when reading data.
			if bufferWait && isTimeout(err) {
				if err = c.SetReadDeadline(receiveTime); err != nil {
					err = fmt.Errorf("SetReadDeadline for connection failed:%s", err.Error())
					return nil, err
				}
				err = nil
				break
			}
			break
		}
		// Just read once from buffer.
		if length == 0 {
			break
		}
	}
	return buffer[:index], err
}

// isTimeout checks whether given `err` is a timeout error.
func isTimeout(err error) bool {
	if err == nil {
		return false
	}
	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		return true
	}
	return false
}
