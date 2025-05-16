package data

import (
	"context"

	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

func TestData(t *testing.T) {
	// 设置环境变量为实际的值:
	// export DB_SOURCE="postgresql://postgres:postgres@192.168.3.121:5432/ecommence?sslmode=disable&timezone=Asia/Shanghai"
	// export REDIS_ADDRESS="192.168.3.132:6379"
	// export REDIS_USERNAME="default"
	// export REDIS_PASSWORD="msdnmm"
	// 在外部的终端工具执行. 在IDE内点击运行可能无法直接读取该环境变量

	// 初始化 Viper 以读取环境变量
	viper.AutomaticEnv()

	// 测试数据库连接
	t.Run("TestDatabaseConnection", func(t *testing.T) {
		// 读取数据库连接字符串
		databaseURL := viper.GetString("DB_SOURCE")
		if databaseURL == "" {
			t.Fatal("DB_SOURCE environment variable is not set")
		}

		// 连接数据库
		conn, err := pgxpool.New(context.Background(), databaseURL)
		if err != nil {
			t.Errorf("Unable to connect to database: %v", err)
			return
		}
		defer conn.Close()

		// 测试数据库 Ping
		err = conn.Ping(context.Background())
		if err != nil {
			t.Errorf("Unable to ping database: %v", err)
		} else {
			t.Log("Database connection successful")
		}
	})

	// 测试 Redis 缓存
	t.Run("TestRedisCache", func(t *testing.T) {
		// 读取 Redis 配置
		redisAddress := viper.GetString("REDIS_ADDRESS")
		redisUsername := viper.GetString("REDIS_USERNAME")
		redisPassword := viper.GetString("REDIS_PASSWORD")

		if redisAddress == "" || redisUsername == "" || redisPassword == "" {
			t.Fatal("Redis environment variables are not set")
		}

		// 创建 Redis 客户端
		client := redis.NewClient(&redis.Options{
			Addr:     redisAddress,
			Username: redisUsername,
			Password: redisPassword,
			DB:       0,
		})
		defer client.Close()

		// 测试 Redis 设置和获取
		ctx := context.Background()
		keys := []string{"example_key", "example_key2", "example_key3", "example_key4"}
		values := []string{"example_value2", "example_value3", "example_value31", "example_value34"}

		// 设置键值对
		for i, key := range keys {
			err := client.Set(ctx, key, values[i], 0).Err()
			if err != nil {
				t.Errorf("Failed to set Redis key %s: %v", key, err)
			}
		}

		// 获取键值对
		for i, key := range keys {
			val, err := client.Get(ctx, key).Result()
			if err != nil {
				t.Errorf("Failed to get Redis key %s: %v", key, err)
			} else if val != values[i] {
				t.Errorf("Value mismatch for key %s: expected %s, got %s", key, values[i], val)
			} else {
				t.Logf("Redis key %s: %s", key, val)
			}
		}
	})
}
