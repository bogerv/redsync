package goredis

import (
	"strings"
	"time"

	"github.com/go-redis/redis"
	redsyncredis "github.com/go-redsync/redsync/v3/redis"
)

type ClusterPool struct {
	delegate *redis.ClusterClient
}

func (self *ClusterPool) Get() redsyncredis.Conn {
	return &ClusterConn{self.delegate}
}

func NewClusterPool(delegate *redis.ClusterClient) *ClusterPool {
	return &ClusterPool{delegate}
}

type ClusterConn struct {
	delegate *redis.ClusterClient
}

func (self *ClusterConn) Get(name string) (string, error) {
	value, err := self.delegate.Get(name).Result()
	err = noErrNil(err)
	return value, err
}

func (self *ClusterConn) Set(name string, value string) (bool, error) {
	reply, err := self.delegate.Set(name, value, 0).Result()
	return err == nil && reply == "OK", nil
}

func (self *ClusterConn) SetNX(name string, value string, expiry time.Duration) (bool, error) {
	return self.delegate.SetNX(name, value, expiry).Result()
}

func (self *ClusterConn) PTTL(name string) (time.Duration, error) {
	return self.delegate.PTTL(name).Result()
}

func (self *ClusterConn) Eval(script *redsyncredis.Script, keysAndArgs ...interface{}) (interface{}, error) {
	var keys []string
	var args []interface{}

	if script.KeyCount > 0 {

		keys = []string{}

		for i := 0; i < script.KeyCount; i++ {
			keys = append(keys, keysAndArgs[i].(string))
		}

		args = keysAndArgs[script.KeyCount:]

	} else {
		keys = []string{}
		args = keysAndArgs
	}

	v, err := self.delegate.EvalSha(script.Hash, keys, args...).Result()
	if err != nil && strings.HasPrefix(err.Error(), "NOSCRIPT ") {
		v, err = self.delegate.Eval(script.Src, keys, args...).Result()
	}
	err = noErrNil(err)
	return v, err
}

func (self *ClusterConn) Close() error {
	// Not needed for this library
	return nil
}
