package cos

import (
	"errors"
	"net"
	"reflect"
	"strings"
)

// GetLocalIP 获取本机ip地址
func GetLocalIP() (ips []string, err error) {
	ips = make([]string, 0)
	faces, err := net.Interfaces()
	if err != nil {
		return
	}
	for _, i := range faces {
		adds, err := i.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range adds {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
				if ip.IsGlobalUnicast() {
					ips = append(ips, ip.String())
				}
			}
		}
	}
	return
}

// GetNetIP 通过网络，获取本机ip地址
func GetNetIP() (ip string, err error) {
	conn, err := net.Dial("udp", "google.com:80")
	if err != nil {
		return "", errors.New("获取主机ip失败")
	}
	defer conn.Close()
	return strings.Split(conn.LocalAddr().String(), ":")[0], nil
}

// Call 调用方法
func Call(method interface{}, params ...interface{}) ([]reflect.Value, error) {
	if reflect.TypeOf(method).Kind() != reflect.Func {
		return nil, errors.New("the name of input not func")
	}
	f := reflect.ValueOf(method)
	if len(params) != f.Type().NumIn() {
		return nil, errors.New("the number of input params not match")
	}
	var mp = make([]reflect.Value, 0)
	for _, v := range params {
		mp = append(mp, reflect.ValueOf(v))
	}
	return f.Call(mp), nil
}
