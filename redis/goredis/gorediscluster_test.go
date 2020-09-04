package goredis

import "github.com/go-redsync/redsync/v3/redis"

var _ (redis.Conn) = (*ClusterConn)(nil)

var _ (redis.Pool) = (*ClusterPool)(nil)
