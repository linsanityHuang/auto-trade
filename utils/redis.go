package utils

import (
	"auto-trade/global"
	"fmt"
	"log"
	"strconv"

	"github.com/go-redis/redis"
)

const (
	key string = "pending"
)

func RPush(ids ...int64) error {

	for _, id := range ids {
		if err := global.RedisClient.RPush(key, fmt.Sprintf("%d", id)).Err(); err != nil {
			log.Printf("rpush err: %v\n", err)
			panic(err)
		}
	}

	return nil
}

func LPop() (int64, error) {
	r, err := global.RedisClient.LPop(key).Result()
	if err != redis.Nil && err != nil {
		// log.Printf("r: %s, err: %T, %#v\n", r, err, err)
		// panic(err)
		log.Fatalf("lpop failed: %v\n", err)
	}

	if r == "" {
		return 0, nil
	}

	return strconv.ParseInt(r, 10, 64)
}
