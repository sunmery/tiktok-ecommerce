package test

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"
	"testing"
)

func TestCitus(t *testing.T) {
	// 设置环境变量为实际的值: export DB_SOURCE="postgresql://postgres:postgres@192.168.3.121:5432/ecommence?sslmode=disable&timezone=Asia/Shanghai"
	// 在外部的终端工具执行. 在IDE内点击运行可能无法直接读取该环境变量

	viper.AutomaticEnv() // 自动读取环境变量
	// 使用 Viper 读取环境变量
	databaseURL := viper.GetString("DB_SOURCE")
	if databaseURL == "" {
		t.Fatal("DB_SOURCE environment variable is not set")
	}

	conn, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		t.Error(fmt.Sprintf("Unable to connect to database: %v", err))
	}
	err = conn.Ping(context.Background())
	if err != nil {
		t.Errorf("Unable to ping database: %v", err)
	}
}
