package authservice

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	u "github.com/sfidann/auth-service/utils"
)

type Auth struct {
	signingKey  []byte
	redisClient *redis.Client
}

func LoadEnv(envFilePath string) {
	err := godotenv.Load(envFilePath)
	if err != nil {
		log.Fatalf("Error loading env file: %v", err)
	}
}

func NewAuth(signingKey string, redisAddr string, redisPassword string) *Auth {
	return &Auth{
		signingKey: []byte("your-256-bit-secret"),
		redisClient: redis.NewClient(&redis.Options{
			Addr:     "3.74.235.128:6379",
			Password: "123456",
			DB:       0,
		}),
	}
}

func GetRedis() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "3.74.235.128:6379",
		Password: "123456",
		DB:       0,
	})

	return rdb

}

func CheckRDB(uuid, accessToken string) error {
	splitted := strings.Split(accessToken, ".") //the token normally comes in format `Bearer {token-body}`
	payload := splitted[1]

	redisclient := GetRedis()
	defer redisclient.Close()

	ctx := context.Background()

	_, err := redisclient.HGet(ctx, uuid, payload).Result()

	return err
}

func AuthMiddleware(ctx *fiber.Ctx) error {
	//check endpoint
	notAuth := []string{"/api/post"}
	requestPath := ctx.Path()
	for _, value := range notAuth {
		if value == requestPath {
			ctx.Next()
			return nil
		}
	}

	//get token from Authorization header info
	authHeader := string(ctx.Request().Header.Peek("Authorization"))
	if !strings.HasPrefix(authHeader, "Bearer ") {
		fmt.Print("Invalid/Malformed auth token")
		return nil
	}
	//parse token and check validity
	accessToken := strings.TrimPrefix(authHeader, "Bearer ")
	claims := &u.AccessToken{}
	tkn, err := jwt.ParseWithClaims(accessToken, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte("your-256-bit-secret"), nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			resp := u.Response{Success: false, Message: "Parse failed."}
			response, _ := json.Marshal(resp)
			ctx.Response().SetStatusCode(401)
			ctx.Response().Header.Add("Content-Type", "application/json")
			ctx.Write(response)
			return nil
		}
		resp := u.Response{Success: false, Message: "Malformed access token."}
		response, _ := json.Marshal(resp)
		ctx.Response().SetStatusCode(401)
		ctx.Response().Header.Add("Content-Type", "application/json")
		ctx.Write(response)
		return nil
	}
	if !tkn.Valid {
		resp := u.Response{Success: false, Message: "Malformed access token."}
		response, _ := json.Marshal(resp)
		ctx.Response().SetStatusCode(401)
		ctx.Response().Header.Add("Content-Type", "application/json")
		ctx.Write(response)
		return nil
	}
	//todo check redis
	err = CheckRDB(claims.Userid, accessToken)
	if err == redis.Nil {
		resp := u.Response{Success: false, Message: "Failed to authenticate from Redis."}
		response, _ := json.Marshal(resp)
		ctx.Response().SetStatusCode(401)
		ctx.Response().Header.Add("Content-Type", "application/json")
		ctx.Write(response)
		return nil
	}

	ctx.Response().SetStatusCode(200)
	ctx.Next()
	return nil
}
