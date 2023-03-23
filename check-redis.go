package authservice

import (
	"context"
	"strings"

	"github.com/redis/go-redis/v9"
)

func CheckRedis(conn *redis.Conn, uuid, accessToken string) error {
	//get payload from refresh token
	splitted := strings.Split(accessToken, ".") //the token normally comes in format `Bearer {token-body}`
	payload := splitted[1]
	//connection to redis
	defer CloseRedisDB(conn)
	ctx := context.Background()

	//check hash key-field
	_, err := conn.HGet(ctx, uuid, payload).Result()

	return err
}
