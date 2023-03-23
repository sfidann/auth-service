package main

import (
	"encoding/json"
	"log"

	"github.com/gofiber/fiber/v2"
	a "github.com/sfidann/auth-service"
)

func main() {
	//conn := a.GetRedisDB("C:/Users/Sevgi/Desktop/.env")

	app := fiber.New()
	api := app.Group("/api")

	/* api.Use(func(ctx *fiber.Ctx) error {
		a.Authentication(ctx, conn, "your-256-bit-secret")
		status := ctx.Response().StatusCode()
		if status == 200 {
			ctx.Next()
			return nil
		}
		return nil
	}) */
	/*
		var c *fiber.Ctx
		api.Use(a.Authentication(c, conn, "your-256-bit-secret")) */

	api.Use(a.AuthMiddleware)

	api.Get("/api/get", func(ctx *fiber.Ctx) error {
		data, _ := json.Marshal("It works‼")
		ctx.Response().SetStatusCode(200)
		ctx.Response().Header.Add("Content-Type", "application/json")
		ctx.Write(data)
		return nil
	})

	err := app.Listen(":3000")
	if err != nil {
		log.Fatal(err)
	}
}
