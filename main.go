package main

import (
	"github.com/gin-gonic/gin"
	"github.com/mdokusV/go-Pracownia/controllers"
	"github.com/mdokusV/go-Pracownia/initializers"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToSQLite()
}

func main() {
	r := gin.Default()
	r.POST("/posts", controllers.PostsCreate)
	r.Run() // listen and serve on 0.0.0.0:3000
}
