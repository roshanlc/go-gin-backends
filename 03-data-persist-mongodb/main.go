// @title          Recipes API
// @version         1.0.0
//
// @description     This is a basic recipes API in go (Gin). A CRUD demonstration.

// @contact.name   Roshan Lamichhane

// @host      localhost:9000
// @BasePath  /

// @accept application/json
// @produce pplication/json

package main

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/roshanlc/go-gin-backends/03-data-persist-mongodb/docs"
	"github.com/roshanlc/go-gin-backends/03-data-persist-mongodb/models"
	"github.com/rs/xid"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// data storage
var recipes []models.Recipe

func init() {
	// initialize the Recipe slice
	recipes = make([]models.Recipe, 0)
}

// NewRecipeHandler godoc
// @Summary      Create a new recipe
// @Description  Create a new recipe
// @Accept       json
// @Produce      json
// @Success      201  {object}  models.Recipe
// @Failure      400
// @Router       /recipes [post]
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

// ListRecipesHandler godoc
// @Summary      List all the recipes
// @Description  List all the recipes
// @Accept       json
// @Produce      json
// @Success      200  {array}  models.Recipe
// @Failure      500
// @Router       /recipes [get]
// handler for GET Method to "/recipes" endpoint
func ListRecipesHandler(c *gin.Context) {
	c.JSON(http.StatusOK, recipes)
}

// UpdateRecipeHandler godoc
// @Summary      Update a recipe
// @Description  Update a recipe
// @Accept       json
// @Produce      json
// @Success      200  {object}  models.Recipe
// @Failure      404
// @Router       /recipes/id [put]
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

// DeleteRecipeHandler
// @Summary      Delete a recipe
// @Description  Delete a recipe
// @Accept       json
// @Produce      json
// @Success      200
// @Failure      404
// @Router       /recipes/id [delete]
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

	// enable global cors
	router.Use(cors.New(cors.Config{AllowAllOrigins: true}))

	// routes for recipes CRUD
	router.POST("/recipes", NewRecipeHandler)
	router.GET("/recipes", ListRecipesHandler)
	router.PUT("/recipes/:id", UpdateRecipeHandler)
	router.DELETE("/recipes/:id", DeleteRecipeHandler)
	router.GET("/recipes/search", SearchHandler)

	// route for swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// run the server
	router.Run(":9000")
}
