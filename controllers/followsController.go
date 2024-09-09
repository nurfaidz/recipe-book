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

func GetFollow(c *gin.Context) {
	db := database.GetDB()
	socialMedia := []models.Follows{}
	followId := c.Param("followId")

	if followId != "" {
		result := db.Where("id = ?", followId).Find(&socialMedia)
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

		c.JSON(http.StatusOK, socialMedia[0])
		return
	}

	err := db.Find(&socialMedia).Error

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, socialMedia)
}

func CreateFollow(c *gin.Context) {
	db := database.GetDB()
	userData := c.MustGet("userData").(jwt.MapClaims)
	contentType := helpers.GetContentType(c)

	Follow := models.Follows{}
	userID := uint(userData["id"].(float64))

	if contentType == appJSON {
		c.ShouldBindJSON(&Follow)
	} else {
		c.ShouldBind(&Follow)
	}

	Follow.FollowerID = userID

	err := db.Debug().Create(&Follow).Error

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})

		return
	}

	c.JSON(http.StatusCreated, Follow)
}

func DeleteFollow(c *gin.Context) {
	db := database.GetDB()
	followId, _ := strconv.Atoi(c.Param("followId"))
	Follow := models.Follows{}

	err := db.Where("id = ?", followId).Delete(&Follow).Error

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Follow with id %v has been deleted", followId),
	})
}

func HelloFollow(g *gin.Context) {
	g.JSON(http.StatusOK, "hello world")
}
