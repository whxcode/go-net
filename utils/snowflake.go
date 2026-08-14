package utils

import (
	"strconv"

	"github.com/bwmarrin/snowflake"
)

var snowflakeNode *snowflake.Node

func InitSnowflake(nodeID int64) error {
	node, err := snowflake.NewNode(nodeID)
	if err != nil {
		return err
	}
	snowflakeNode = node
	return nil
}

func GenerateSnowflakeID() string {
	return strconv.FormatInt(snowflakeNode.Generate().Int64(), 10)
}
