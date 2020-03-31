package etcd_test

import (
	"context"
	"cope/registry"
	"fmt"
	"testing"
	"time"
)

func TestEtcd(t *testing.T) {
	registryInit, err := registry.InitRegistry(
		context.TODO(),
		"etcd",
		registry.WithAddrs([]string{"127.0.0.1:2379"}),
		registry.WithHeartBeat(5),
		registry.WithTimeout(time.Second),
		registry.WithRegistryPath("/ibinarytree/cope/"),
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	service := &registry.Service{
		Name: "comment_service",
	}

	service.Nodes = append(service.Nodes, &registry.Nodes{
		IP:   "127.0.0.1",
		Port: 8801,
	},
		&registry.Nodes{
			IP:   "127.0.0.2",
			Port: 8801,
		},
	)

	err = registryInit.Register(context.TODO(), service)
	if err != nil {
		fmt.Println(err)
		return
	}
	go func() {
		time.Sleep(time.Second * 15)
		service.Nodes = append(service.Nodes, &registry.Nodes{
			IP:   "127.0.0.3",
			Port: 8801,
		},
			&registry.Nodes{
				IP:   "127.0.0.4",
				Port: 8801,
			},
		)
		registryInit.Register(context.TODO(), service)
	}()

	for {
		service, err := registryInit.GetService(context.TODO(), "comment_service")
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("%#v \n", service)

		// for _, node := range service.Nodes {
		// 	fmt.Printf("%#v \n", node)
		// }

		// fmt.Println(1111111111)
		time.Sleep(time.Second * 5)
	}

}
