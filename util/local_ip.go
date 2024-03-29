package util

import (
	"context"
	"cope/logs"
	"fmt"
	"net"
	"sync/atomic"
)

var ipAuto atomic.Value

//GetLocalIP 获取本地ip
func GetLocalIP() (ip string, err error) {

	ip, ok := ipAuto.Load().(string)
	if ok && len(ip) > 0 {
		return
	}

	netInterfaces, err := net.Interfaces()
	if err != nil {
		return
	}

	for _, interfaces := range netInterfaces {
		if (interfaces.Flags & net.FlagUp) != 0 {
			addressAll, _ := interfaces.Addrs()
			for _, address := range addressAll {
				// 检查ip地址判断是否回环地址
				if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					if ipnet.IP.To4() != nil {
						ip = ipnet.IP.String()
						ipAuto.Store(ip)
						logs.Info(context.TODO(), "local ip:%v", ip)
						return
					}
				}
			}
		}
	}

	err = fmt.Errorf("get local ip failed")
	return
}
