package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/roshanlc/go-gin-backends/02-basic-rest-api/models"
	"github.com/rs/xid"
)

var recipes []models.Recipe

func init() {
	// initialize the Recipe slice
	recipes = make([]models.Recipe, 0)
}

func NewRecipeHandler(c *gin.Context) {
	var recipe models.Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "malformed JSON: " + err.Error(),
		})
		return
	}

	recipe.ID = xid.New().String()
	recipe.PublishedAt = time.Now()
	// Add it to the global in-memory slice
	recipes = append(recipes, recipe)

	// return 201-Created status code
	c.JSON(http.StatusCreated, recipe)

}

func ListRecipesHandler(c *gin.Context) {
	c.JSON(http.StatusOK, recipes)
}

func main() {
	router := gin.Default()
	router.POST("/recipes", NewRecipeHandler)
	router.GET("/", ListRecipesHandler)
	router.Run(":9000")
}
