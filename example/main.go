package main

import (
	"net/http"

	"github.com/finejian/paginator"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.LoadHTMLFiles("example.html", "example.html")
	router.GET("/example", func(c *gin.Context) {
		p := paginator.Custom(&paginator.Config{PageSize: 10, Current: 6, LinkedCount: 4}, 63)
		p.Request(c.Request)

		c.HTML(http.StatusOK, "example.html", gin.H{
			"paginator": p,
		})
	})
	router.Run(":8080")
}
