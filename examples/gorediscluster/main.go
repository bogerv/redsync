package main

import (
	goredislib "github.com/go-redis/redis"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis"
	"github.com/go-redsync/redsync/v4/redis/goredis"
)

func main() {
	client := goredislib.NewClusterClient(&goredislib.ClusterOptions{
		Addrs: []string{
			"101.32.1.100:6379",
			"101.32.1.101:6379",
			"101.32.1.10:6379",
			"101.32.1.100:6380",
			"101.32.1.101:6380",
			"101.32.1.10:6380",
		},
	})

	pool := goredis.NewClusterPool(client)

	rs := redsync.New([]redis.Pool{pool}...)

	mutex := rs.NewMutex("test-redsync")
	err := mutex.Lock()

	if err != nil {
		panic(err)
	}

	mutex.Unlock()
}
