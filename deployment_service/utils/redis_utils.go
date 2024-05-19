package utils

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client
var buildQueue = "build-queue"

func init() {

	godotenv.Load(".env")
	log.Default().Println("Init Redis")
	redisURI, _ := os.LookupEnv("REDIS_URI")
	opt, _ := redis.ParseURL(redisURI)
	rdb := redis.NewClient(opt)
	redisClient = rdb
	_, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}
	log.Default().Println("Connected to Redis")
}

func PullFromRedisBuildQueue() (string, error) {
	val, err := redisClient.RPop(context.TODO(), buildQueue).Result()
	if err != nil {
		log.Default().Println(err)
		return "", err
	}
	return val, nil
}
