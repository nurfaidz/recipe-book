package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"recipebook/database"
	"recipebook/helpers"
	"recipebook/models"
	"strconv"
)

func GetRecipe(c *gin.Context) {
	db := database.GetDB()
	recipes := []models.Recipes{}
	recipeId := c.Param("recipeId")

	if recipeId == "" {
		result := db.Where("id == ?", recipeId).Find(&recipes)
		err := result.Error
		count := result.RowsAffected

		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error":   "bad Request",
				"message": "Invalid Parameter",
			})

			return
		}

		if count < 1 {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error":   "Data not found",
				"message": "data doesn't exist",
			})

			return
		}

		c.JSON(http.StatusOK, recipes[0])
		return
	}

	err := db.Find(&recipes).Error

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, recipes)
}

func UpdateRecipe(c *gin.Context) {
	db := database.GetDB()
	userData := c.MustGet("UserData").(jwt.MapClaims)
	contentType := helpers.GetContentType(c)

	Recipe := models.Recipes{}

	recipeId, _ := strconv.Atoi(c.Param("recipeId"))
	userID := uint(userData["id"].(float64))

	if contentType == appJSON {
		c.ShouldBindJSON(&Recipe)
	} else {
		c.ShouldBind(&Recipe)
	}

	Recipe.UserID = userID
	Recipe.ID = uint(recipeId)

	err := db.Model(&Recipe).Where("id = ?", recipeId).Updates(models.Recipes{Title: Recipe.Title, Description: Recipe.Description, Ingredients: Recipe.Ingredients, Steps: Recipe.Steps, Category: Recipe.Category, Tags: Recipe.Tags, PictureUrl: Recipe.PictureUrl}).Error

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, Recipe)
}

func DeleteRecipe(c *gin.Context) {
	db := database.GetDB()
	recipeId, _ := strconv.Atoi(c.Param("recipeId"))
	Recipe := models.Recipes{}

	err := db.Debug().Where("id = ?", recipeId).Delete(&Recipe).Error

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Recipe with ID %d has been deleted", recipeId),
	})
}

func HelloRecipe(g *gin.Context) {
	g.JSON(http.StatusOK, "hello world")
}

func CreateRecipeComment(c *gin.Context) {
	db := database.GetDB()
	userData := c.MustGet("UserData").(jwt.MapClaims)
	contentType := helpers.GetContentType(c)

	Comment := models.Comments{}
	userID := uint(userData["id"].(float64))
	recipeId, _ := strconv.Atoi(c.Param("recipeId"))

	if contentType == appJSON {
		c.ShouldBindJSON(&Comment)
	} else {
		c.ShouldBind(&Comment)
	}

	Comment.UserID = userID
	Comment.RecipeID = uint(recipeId)

	err := db.Debug().Create(&Comment).Error

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})

		return
	}

	c.JSON(http.StatusCreated, Comment)
}

func CreateRecipeLike(c *gin.Context) {
	db := database.GetDB()
	userData := c.MustGet("UserData").(jwt.MapClaims)
	recipeId, _ := strconv.Atoi(c.Param("recipeId"))

	Like := models.Likes{}
	userID := uint(userData["id"].(float64))

	Like.UserID = userID
	Like.RecipeID = uint(recipeId)

	result := db.Debug().Where("recipe_id = ?", recipeId).Where("user_id = ?", userID).First(&Like)
	count := result.RowsAffected

	if count > 0 {
		err := db.Debug().Where("recipe_id = ?", recipeId).Where("user_id = ?", userID).Delete(&Like).Error

		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error":   "Bad Request",
				"message": err.Error(),
			})

			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Like with recipe_id %v and user_id %v has been deleted", recipeId, userID),
		})

		return
	}

	errCreate := db.Debug().Create(&Like).Error

	if errCreate != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": errCreate.Error(),
		})

		return
	}

	c.JSON(http.StatusCreated, Like)
}

func CreateRecipeFollow(c *gin.Context) {
	db := database.GetDB()
	userData := c.MustGet("UserData").(jwt.MapClaims)

	Recipe := models.Recipes{}
	Follow := models.Follows{}
	userID := uint(userData["id"].(float64))
	recipeId, _ := strconv.Atoi(c.Param("recipeId"))

	recipe := db.Debug().Where("id = ?", recipeId).First(&Recipe)
	count := recipe.RowsAffected

	if count < 1 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Recipe doesn't exist",
		})

		return
	}

	Follow.FollowedID = Recipe.UserID
	Follow.FollowerID = userID

	result := db.Where("followed_id == ?", Recipe.UserID).Where("follower_id == ?", userID).First(&Follow)
	countResult := result.RowsAffected

	if countResult > 0 {
		err := db.Where("followed_id == ?", Recipe.UserID).Where("follower_id == ?", userID).Delete(&Follow).Error

		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error":   "Bad Request",
				"message": err.Error(),
			})

			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Follow with followed_id %v and follower_id %v has been deleted", Recipe.UserID, userID),
		})

		return
	}

	errCreate := db.Debug().Create(&Follow).Error

	if errCreate != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": errCreate.Error(),
		})

		return
	}

	c.JSON(http.StatusCreated, Follow)
}

func CreateNewRecipe(c *gin.Context) {
	db := database.GetDB()
	userData := c.MustGet("UserData").(jwt.MapClaims)
	contentType := helpers.GetContentType(c)

	Recipe := models.Recipes{}
	userID := uint(userData["id"].(float64))

	if contentType == appJSON {
		c.ShouldBindJSON(&Recipe)
	} else {
		c.ShouldBind(&Recipe)
	}

	Recipe.UserID = userID

	err := db.Debug().Create(&Recipe).Error

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, Recipe)
}
