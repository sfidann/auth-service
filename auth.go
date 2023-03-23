package authservice

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	u "github.com/sfidann/auth-service/utils"
)

func Authentication(ctx *fiber.Ctx, conn *redis.Conn, tokenKey string) error {

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
		return []byte(tokenKey), nil
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
	err = CheckRedis(conn, claims.Userid, accessToken)
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
