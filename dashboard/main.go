package main

import (
	"gitlab.com/pro-tracker/dashboard/app/handler/email"
	"gitlab.com/pro-tracker/dashboard/app/handler/html"

	"github.com/gin-gonic/gin"
)

var router *gin.Engine

func main() {
	router = gin.Default()
	initializeRoutes()
	router.Run(":8000")

}

func initializeRoutes() {
	//load HTML
	router.LoadHTMLGlob("frontend/*")
	router.GET("/", html.IndexPage)
	//CRUD operations in email entity
	router.POST("/email", email.SendEmail)
	router.GET("/email", email.ListEmail)
}
