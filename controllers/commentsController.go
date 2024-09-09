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

func GetComment(c *gin.Context) {
	db := database.GetDB()
	comments := []models.Comments{}
	commentId := c.Param("commentId")

	if commentId != "" {
		result := db.Where("id = ?", commentId).Find(&comments)
		err := result.Error
		count := result.RowsAffected

		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error":   "Bad Request",
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

		c.JSON(http.StatusOK, comments[0])
		return
	}

	err := db.Find(&comments).Error

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, comments)
}

func CreateComment(c *gin.Context) {
	db := database.GetDB()
	userData := c.MustGet("userData").(jwt.MapClaims)
	contentType := helpers.GetContentType(c)

	Comment := models.Comments{}
	userID := uint(userData["id"].(float64))

	if contentType == appJSON {
		c.ShouldBindJSON(&Comment)
	} else {
		c.ShouldBind(&Comment)
	}

	Comment.UserID = userID

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

func UpdateComment(c *gin.Context) {
	db := database.GetDB()
	userData := c.MustGet("userData").(jwt.MapClaims)
	contentType := helpers.GetContentType(c)

	Comment := models.Comments{}

	commentId, _ := strconv.Atoi(c.Param("commentId"))
	userID := uint(userData["id"].(float64))

	if contentType == appJSON {
		c.ShouldBindJSON(&Comment)
	} else {
		c.ShouldBind(&Comment)
	}

	Comment.UserID = userID
	Comment.ID = uint(commentId)

	err := db.Model(&Comment).Where("id = ?", commentId).Updates(models.Comments{Message: Comment.Message, RecipeID: Comment.RecipeID}).Error

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, Comment)
}

func DeleteComment(c *gin.Context) {
	db := database.GetDB()
	commentId, _ := strconv.Atoi(c.Param("commentId"))
	Comment := models.Comments{}

	err := db.Debug().Where("id = ?", commentId).Delete(&Comment).Error

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Comment with ID %d has been deleted", commentId),
	})
}

func HelloComment(g *gin.Context) {
	g.JSON(http.StatusOK, "hello world")
}
