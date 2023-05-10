package main

import "github.com/gin-gonic/gin"

// handler for "/" endpoint
func indexHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "hello, world!",
	})
}

// handler for the "/:name" endpoint
func nameHandler(c *gin.Context) {
	name := c.Params.ByName("name")
	c.JSON(200, gin.H{
		"message": "hello, " + name,
	})
}
func main() {
	// default gin Engine
	router := gin.Default()

	// registering a GET endpoint
	router.GET("/", indexHandler)

	router.GET("/:name", nameHandler)

	// run the router
	// .i.e start the server at specified port
	router.Run(":9000")
}
