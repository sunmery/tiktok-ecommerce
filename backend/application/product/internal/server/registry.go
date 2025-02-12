package server

import (
	"backend/application/product/internal/conf"
	"github.com/go-kratos/kratos/contrib/registry/consul/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/hashicorp/consul/api"
)

// NewRegistrar 使用Consul作为注册中心
func NewRegistrar(c *conf.Consul) registry.Registrar {
	// 创建consul客户端
	configs := &api.Config{
		Address: c.RegistryCenter.Address,
		Scheme:  c.RegistryCenter.Scheme,
		Token:   c.RegistryCenter.AclToken,
	}
	// 从conf/conf.proto定义Consul配置与configs/config.yml配置文件中读取consul的配置
	// 创建consul客户端
	consulClient, err := api.NewClient(configs)
	if err != nil {
		log.Fatal(err)
	}
	// 创建consul注册中心
	r := consul.New(consulClient, consul.WithHealthCheck(c.RegistryCenter.HealthCheck))
	return r
}
