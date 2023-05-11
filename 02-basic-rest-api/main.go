package main

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/roshanlc/go-gin-backends/02-basic-rest-api/models"
	"github.com/rs/xid"
)

// data storage
var recipes []models.Recipe

func init() {
	// initialize the Recipe slice
	recipes = make([]models.Recipe, 0)
}

// handler for POST Method to "/recipes" endpoint
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

// handler for GET Method to "/recipes" endpoint
func ListRecipesHandler(c *gin.Context) {
	c.JSON(http.StatusOK, recipes)
}

// handler for PUT Method to "/recipes/:id" endpoint
func UpdateRecipeHandler(c *gin.Context) {
	id := c.Param("id")
	var recipe models.Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "malformed JSON: " + err.Error(),
		})
		return
	}
	index := -1
	for i := 0; i < len(recipes); i++ {
		if recipes[i].ID == id {
			index = i
		}
	}

	if index == -1 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Recipe not found",
		})
		return
	}

	recipes[index] = recipe
	// assign the old id
	recipe.ID = id
	// assign the current time as it was updated
	recipe.PublishedAt = time.Now()
	c.JSON(http.StatusOK, recipe)
}

// handler for DELETE Method to "/recipes/:id" endpoint
func DeleteRecipeHandler(c *gin.Context) {
	id := c.Param("id")

	index := -1
	for i := 0; i < len(recipes); i++ {
		if recipes[i].ID == id {
			index = i
		}
	}

	if index == -1 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Recipe not found",
		})
		return
	}

	recipes = append(recipes[:index], recipes[index+1:]...)

	c.JSON(http.StatusOK, gin.H{
		"message": "Recipe has been deleted",
	})
}

// handler for GET method to "/recipes/search?tag=" endpoint
func SearchHandler(c *gin.Context) {
	tag := c.Query("tag")

	recipesList := make([]models.Recipe, 0)

	for i := 0; i < len(recipes); i++ {
		// break on first find
		found := false
		for _, t := range recipes[i].Tags {
			if strings.EqualFold(t, tag) {
				found = true
			}
		}

		if found {
			recipesList = append(recipesList, recipes[i])
		}
	}

	c.JSON(http.StatusOK, recipesList)
}

func main() {
	router := gin.Default()
	router.POST("/recipes", NewRecipeHandler)
	router.GET("/recipes", ListRecipesHandler)
	router.PUT("/recipes/:id", UpdateRecipeHandler)
	router.DELETE("/recipes/:id", DeleteRecipeHandler)
	router.GET("/recipes/search", SearchHandler)
	router.Run(":9000")
}
