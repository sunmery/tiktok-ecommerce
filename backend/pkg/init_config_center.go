package pkg

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/go-kratos/kratos/contrib/config/consul/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/hashicorp/consul/api"
)

type ConfigCenter struct {
	// 配置中心地址
	Addr string
	// 微服务对应的配置文件路径
	Path  string
	Token string
}

func InitConsul(config ConfigCenter) config.Source {
	flag.Parse()

	// 如果环境变量存在，覆盖默认值
	if envConfigCenter := os.Getenv("config_center"); envConfigCenter != "" {
		config.Addr = envConfigCenter
	}
	if envConfigPath := os.Getenv("config_path"); envConfigPath != "" {
		config.Path = envConfigPath
	}
	if envConfigCenterToken := os.Getenv("config_center_token"); envConfigCenterToken != "" {
		config.Token = envConfigCenterToken
	}

	// debug
	log.Debugf("config_path:%v", config.Path)
	log.Debugf("config_center:%v", config.Addr)
	log.Debugf("config_center_token:%v", config.Token)

	consulClient, err := api.NewClient(&api.Config{
		Address:  config.Addr,
		Scheme:   "http",
		WaitTime: time.Second * 15,
		Token:    config.Token,
	})
	if err != nil {
		log.Fatal(fmt.Errorf("create consul client failed:%w", err))
	}
	cs, err := consul.New(consulClient, consul.WithPath(config.Path))
	if err != nil {
		log.Fatal(fmt.Errorf("create kratos config Source failed:%w", err))
	}
	return cs
}
