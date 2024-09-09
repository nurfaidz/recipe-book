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

func GetLike(c *gin.Context) {
	db := database.GetDB()
	likes := []models.Likes{}
	likeId := c.Param("likeId")

	if likeId == "" {
		result := db.Where("id == ?", likeId).Find(&likes)
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

			c.JSON(http.StatusOK, likes[0])
			return
		}
	}

	err := db.Find(&likes).Error

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, likes)
}

func CreateLike(c *gin.Context) {
	db := database.GetDB()
	userData := c.MustGet("userData").(jwt.MapClaims)
	contentType := helpers.GetContentType(c)

	Like := models.Likes{}
	userID := uint(userData["id"].(float64))

	if contentType == appJSON {
		c.ShouldBindJSON(&Like)
	} else {
		c.ShouldBind(&Like)
	}

	Like.UserID = userID

	err := db.Debug().Create(&Like).Error

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})

		return
	}

	c.JSON(http.StatusCreated, Like)
}

func DeleteLike(c *gin.Context) {
	db := database.GetDB()
	likeId, _ := strconv.Atoi(c.Param("likeId"))
	Like := models.Likes{}

	err := db.Where("id = ?", likeId).Delete(&Like).Error

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Like with id %v has been deleted", likeId),
	})
}

func HelloLike(g *gin.Context) {
	g.JSON(http.StatusOK, "hello world")
}
