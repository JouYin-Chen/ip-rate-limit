package helper

import (
	"os"
	"os/signal"
	"syscall"
	"time"
	"fmt"
	"github.com/garyburd/redigo/redis"
)

type Response struct {
	Value []interface {}
}

var (
	Pool *redis.Pool
)

func init() {
	redisHost := os.Getenv("REDIS_HOST")
	redisPost := os.Getenv("REDIS_PORT")

	if redisHost == "" {
		redisHost = "localhost:6379"
	} else {
		redisHost = fmt.Sprintf("%v:%v", redisHost, redisPost)
	}
	
	fmt.Println("redisHost", redisHost)
	Pool = newPool(redisHost)
	cleanupHook()
}

func newPool(server string) *redis.Pool {
	return &redis.Pool{

		MaxIdle: 3,
		MaxActive: 10,

		IdleTimeout: 240 * time.Second,

		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			return c, err
		},

		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

func cleanupHook() {

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	signal.Notify(c, syscall.SIGKILL)
	go func() {
		<-c
		Pool.Close()
		os.Exit(0)
	}()
}

func SetHashFieldValue(key string, args interface{}) (error){
	conn := Pool.Get()
	defer conn.Close()

	if _, err := conn.Do("HMSET", redis.Args{}.Add(key).AddFlat(args)...); err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func FindValeByHashField(key string, field string) ([]byte, error){
	var i []byte
	conn := Pool.Get()
	defer conn.Close()

	res, err := conn.Do("hget", key, field)
	// fmt.Println(reflect.TypeOf(res))
	if err != nil {
		fmt.Println("hget failed", err)
		return i, err
	} else if res == nil {
		return i, nil
	} else {
		fmt.Printf("hget value: %s\n", res.([]byte))
		return res.([]byte), nil
	}
}

func SetValueExpByKey(key string, value string) {
	conn := Pool.Get()
	defer conn.Close()

	_, err := conn.Do("SETNX", key, value)
	if err != nil{
    fmt.Println("SET failed ",err.Error())
	}
}

func SetExpireKey(key string, ttl int) {
	conn := Pool.Get()
	defer conn.Close()

	n, _ := conn.Do("EXPIRE", key, ttl)
	if n == int64(1){
    fmt.Println("SET EXPIRE success")
	}
}
func FindValueByKey(key string) ([]byte, error){
	var i []byte
	conn := Pool.Get()
	defer conn.Close()

	res, err := conn.Do("get", key)

	if err != nil{
		fmt.Println("FindValueByKey get failed", err)
		return i, err
	} else if res == nil{
		return i, nil
	} else {
		fmt.Printf("get value: %s\n", res.([]byte))
		return res.([]byte), nil
	}
}

func IncrNumberByKey(key string) {
	conn := Pool.Get()
	defer conn.Close()

	_, err := conn.Do("incr", key)
	// fmt.Println(reflect.TypeOf(res))
	if err != nil {
		fmt.Println("incr failed", err)
	} 
}