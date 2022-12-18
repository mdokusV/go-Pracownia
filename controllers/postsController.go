package controllers

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/mdokusV/go-Pracownia/initializers"
	"github.com/mdokusV/go-Pracownia/models"
)

func PostsCreate(c *gin.Context) {
	//Get data from request
	var body struct {
		Title string
		Body  string
	}
	c.Bind(&body)
	//Create POST
	Post := models.Post{Title: body.Title, Body: body.Body}

	result := initializers.DB.Create(&Post)
	if result.Error != nil {
		c.Status(400)
		fmt.Println(result.Error)
		return
	}
	//Return it

	c.JSON(200, gin.H{
		"message": Post,
	})
}
