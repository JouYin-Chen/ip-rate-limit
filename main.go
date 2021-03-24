package main

import (
	"ip-rate-limit/middleware"
	"ip-rate-limit/api"
	"fmt"
	"github.com/gin-gonic/gin"
	
)

func main() {

	app := gin.New()

	fmt.Println("Hello, GOLang!")

	userRouter := app.Group("/")

	userRouter.Use(middleware.RateLimitMiddleware())
	userRouter.GET("/user", user.GetUser)
	userRouter.POST("/user", user.CreateUser)


	app.Run(":3000") // listen and serve on 0.0.0.0:3000
}
