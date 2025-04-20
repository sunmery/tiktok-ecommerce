package id

import (
	"log"

	"backend/constants"

	"github.com/bwmarrin/snowflake"
)

var node *snowflake.Node

// 每次调用 SnowflakeID() 都会重新创建一个新的 snowflake.Node，导致节点号（nodeID）可能重复。
// Snowflake 节点应作为单例（Singleton）初始化一次，而不是在每次生成 ID 时重新创建。
func init() {
	var err error
	node, err = snowflake.NewNode(constants.CheckoutServiceNode)
	if err != nil {
		log.Panicf("Failed to initialize Snowflake node: %v", err)
	}
}

func SnowflakeID() int64 {
	return node.Generate().Int64()
}
