package database

import (
	log "github.com/Sirupsen/logrus"
	"github.com/garyburd/redigo/redis"
)

var c redis.Conn

func Init() (err error) {
	c, err = redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		log.Fatal(err)
	}
	return
}

func Get(section string, key string, values interface{})  {
	v, err := redis.Values(c.Do("HGETALL", section + "-" + key))
	if err != nil {
		log.Fatal(err)
	}
	scanStruct(v, values)
	return
}

func Set(section string, key string, value interface{})  {
	if _, err := c.Do("HMSET", redis.Args{section + "-" + key}.AddFlat(value)...); err != nil {
		log.Fatal(err)
	}
}

func List(section string, key string, values interface{}) {
	v, err := redis.Values(c.Do("SMEMBERS", section + "-" + key))
	redis.ScanSlice(v, values)
	if err != nil {
		log.Fatal(err)
	}
	return
}

func Push(section string, key string, value string) {
	if _, err := c.Do("SADD", section + "-" + key, value); err != nil {
		log.Fatal(err)
	}
}

func Count(section string, key string) (value int)  {
	value, err := redis.Int(c.Do("SCARD", section + "-" + key))
	if err != nil {
		log.Fatal(err)
	}
	return
}

func Remove(section string, key string, value string)  {
	if _, err := c.Do("SREM", section + "-" + key, value); err != nil {
		log.Fatal(err)
	}
}

func scanStruct(src []interface{}, dest interface{}) error {
	err := redis.ScanStruct(src, dest)
	return err
}

func SetIncr(section string) (value int)  {
	value, err := redis.Int(c.Do("INCR", section))
	if err != nil {
		log.Fatal(err)
	}
	return
}
