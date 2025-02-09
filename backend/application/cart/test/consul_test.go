package test

import (
	"fmt"
	"github.com/gin-gonic/gin"
	consul "github.com/hashicorp/consul/api"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"testing"
)

var (
	client *consul.Client
	err    error
	config *consul.Config
	TAGS   = []string{"test"}
)

const ID = "test"
const ServerName = "tiktok-e-commence-test"
const HOST = "159.75.231.54"
const PORT = 8500

func init() {
	config = &consul.Config{
		Address: fmt.Sprintf("%s:%d", HOST, PORT), // consul注册地址, 默认是localhost:8500
		Scheme:  "http",
		// PathPrefix: "",
		// Datacenter: "",
		// Transport:  nil,
		// HttpClient: nil,
		// HttpAuth:   nil,
		// WaitTime:   0,
		// Token:      "",
		// TokenFile:  "",
		// Namespace:  "",
		// Partition:  "",
		// TLSConfig:  consul.TLSConfig{},
	}
	client, err = consul.NewClient(config)
	if err != nil {
		log.Fatal("Consul客户端创建失败:", err)
	}
}

func TestConsulClient(t *testing.T) {
	// 创建Consul客户端

	// 注册服务
	err = registerService(client)
	if err != nil {
		log.Fatal("服务注册失败:", err)
	}

	// 发现服务
	discoveredService, err := discoverService(client, ServerName)
	if err != nil {
		log.Fatal("服务发现失败:", err)
	}

	// 调用发现的服务
	callService(discoveredService)

	// 等待信号以进行优雅关闭
	waitForSignal()
}

func registerService(client *consul.Client) error {
	// 微服务名称采用三段式的命名规则，中间使用中横线分隔，即xxxx-xxxx-xxxx形式
	// 一级服务名为组织名称，如hope，二级服务名为应用或项目的名称，如madp，三级服务名为功能模块的名称，如auth。
	// 整体为hope-madp-auth，使用英文拼写，单词间不要使用空格和_。请全部使用小写字母

	// 创建服务注册信息
	reg := &consul.AgentServiceRegistration{
		ID:      ID,
		Name:    ServerName,
		Tags:    TAGS,
		Address: HOST,
		Port:    PORT,
		// Check: &consul.AgentServiceCheck{
		// 	HTTP:     "http://192.168.2.181:8080/health",
		// 	Interval: "10s",
		// },
	}

	// 注册服务
	err := client.Agent().ServiceRegister(reg)
	if err != nil {
		return err
	}

	fmt.Println("服务注册成功")
	return nil
}

// 配置中心
func TestConsulConfigStore(t *testing.T) {
	key := fmt.Sprintf("%s-%s-%s,", ServerName, ID, "key1")
	value := []byte("value1")

	// 向配置中心添加的配置
	p := &consul.KVPair{
		Key:         key,
		CreateIndex: 0,
		ModifyIndex: 0,
		LockIndex:   0,
		Flags:       0,
		Value:       value,
		Session:     "",
		Namespace:   "",
		Partition:   "",
	}

	// 向配置中心添加一个Key和Value
	pair, err := client.KV().Put(p, nil)
	if err != nil {
		log.Fatal("put error:", err)
	}
	fmt.Println("put result:", pair) // 返回添加完成时的时间, 单位ms

	// 向配置中心查询一个Key
	pair1, _, err1 := client.KV().Get(key, nil)
	if err1 != nil {
		log.Fatal("get error:", err)
	}
	fmt.Println("get result:", pair1)

	// 删除
	pair2, err2 := client.KV().Delete(key, nil)
	if err2 != nil {
		log.Fatal("delete error:", err)
	}
	fmt.Println("delete result:", pair2) // 返回删除的时间, 单位ms
}

// 健康检查
func TestServer(t *testing.T) {
	s := gin.Default()
	s.GET("/livez", func(c *gin.Context) {
		c.Status(200)
	})

	// 创建健康检查规则
	healthCheck := consul.AgentServiceCheck{
		CheckID:                        "4",
		Name:                           ServerName,
		Interval:                       "10s",                         // 检查间隔
		Timeout:                        "2s",                          // 超时时间
		HTTP:                           "http://localhost:4000/livez", // 健康检查的URL
		Method:                         "GET",
		DeregisterCriticalServiceAfter: "1m", // 注销服务的时间
	}

	// 创建一个新的服务实例
	registration := consul.AgentServiceRegistration{
		ID:      "4",
		Name:    ServerName,
		Tags:    TAGS,
		Port:    4000,
		Address: "127.0.0.1",
		Check:   &healthCheck,
	}
	// 注册服务
	registrationErr := client.Agent().ServiceRegister(&registration)
	if registrationErr != nil {
		panic(registrationErr)
	}

	err3 := s.Run(":4000")
	if err3 != nil {
		panic(err3)
	}
}

func TestConsulLivez(t *testing.T) {
	// 执行健康检查
	checks, _, err := client.Health().Checks(ServerName, nil)
	if err != nil {
		panic(err)
	}

	// 打印健康检查结果
	for _, check := range checks {
		fmt.Printf("CheckID: %s, Status: %s\n", check.CheckID, check.Status)
	}
}

func discoverService(client *consul.Client, serviceName string) (*consul.AgentService, error) {
	// 使用Consul进行服务发现
	passingOnly := true
	services, _, err := client.Health().Service(serviceName, "", passingOnly, nil)
	if err != nil {
		return nil, err
	}

	if len(services) == 0 {
		return nil, fmt.Errorf("未找到服务: %s", serviceName)
	}

	// 返回发现的服务
	return services[0].Service, nil
}
func callService(service *consul.AgentService) {
	// 调用发现的服务
	fmt.Printf("调用服务: %s\n", service.Service)
	resp, err := http.Get(fmt.Sprintf("http://%s:%d", service.Address, service.Port))
	if err != nil {
		log.Fatal("服务调用失败:", err)
	}
	defer resp.Body.Close()
	fmt.Println("服务调用成功")
}
func waitForSignal() {
	// 等待信号以进行优雅关闭
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
	fmt.Println("接收到信号，开始关闭服务...")
	// 执行清理操作
	fmt.Println("服务已关闭")
}
