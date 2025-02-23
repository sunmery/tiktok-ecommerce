package pkg

import (
	"log"

	"backend/application/assistant/internal/conf"

	"github.com/bwmarrin/snowflake"
)

func SnowflakeID() int64 {
	// 节点号必须唯一：每个服务实例（或分布式节点）必须分配一个全局唯一的节点号
	node, err := snowflake.NewNode(conf.Server.Snowflake.Node)
	if err != nil {
		log.Panic(err)
	}

	// Generate a snowflake ID.
	id := node.Generate().Int64()

	return id
}
