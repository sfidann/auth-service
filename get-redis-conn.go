package authservice

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

// var redis_addr, redis_pass string
// var redis_db int

func GetRedisDB(envFile string) *redis.Conn {
	err := godotenv.Load(envFile)
	if err != nil {
		log.Fatal(err)
	}

	redis_addr := os.Getenv("RDB_ADDR")
	redis_pass := os.Getenv("RDB_PASS")
	redis_db, _ := strconv.Atoi(os.Getenv("RDB_INT"))

	rdb := redis.NewClient(&redis.Options{
		Addr:     redis_addr,
		Password: redis_pass,
		DB:       redis_db,
	})

	rconn := rdb.Conn()
	return rconn
}

func CloseRedisDB(rconn *redis.Conn) {
	err := rconn.Close()
	if err != nil {
		log.Fatal(err)
	}
}
