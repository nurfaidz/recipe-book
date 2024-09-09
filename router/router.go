package router

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"recipebook/controllers"
	"recipebook/middlewares"
)

func StartServer() *gin.Engine {
	r := gin.Default()

	userRouter := r.Group("/users")
	{
		userRouter.POST("register", controllers.UserRegister)

		userRouter.POST("login", controllers.UserLogin)
	}

	recipeRouter := r.Group("/recipes")
	{
		recipeRouter.Use(middlewares.Authentication())
		recipeRouter.GET("/", controllers.GetRecipe)
		recipeRouter.GET("/:recipeId", controllers.GetRecipe)
		recipeRouter.POST("/", controllers.CreateNewRecipe)
		recipeRouter.PUT("/:recipeId", middlewares.RecipeAuthorization(), controllers.UpdateRecipe)
		recipeRouter.DELETE("/:recipeId", middlewares.RecipeAuthorization(), controllers.DeleteRecipe)
		recipeRouter.POST("/:recipeId/comments", controllers.CreateRecipeComment)
		recipeRouter.POST("/:recipeId/likes", controllers.CreateRecipeLike)
		recipeRouter.POST("/:recipeId/follows", controllers.CreateRecipeFollow)
	}

	likeRouter := r.Group("/likes")
	{
		likeRouter.Use(middlewares.Authentication())
		likeRouter.GET("/", controllers.GetLike)
		likeRouter.GET("/:likeId", controllers.GetLike)
		likeRouter.POST("/", controllers.CreateLike)
		likeRouter.DELETE("/:likeId", middlewares.LikeAuthorization(), controllers.DeleteLike)
	}

	commentRouter := r.Group("/comments")
	{
		commentRouter.Use(middlewares.Authentication())
		commentRouter.GET("/", controllers.GetComment)
		commentRouter.GET("/:commentId", controllers.GetComment)
		commentRouter.POST("/", controllers.CreateComment)
		commentRouter.PUT("/:commentId", middlewares.CommentAuthorization(), controllers.UpdateComment)
		commentRouter.DELETE("/:commentId", middlewares.CommentAuthorization(), controllers.DeleteComment)
	}

	followRouter := r.Group("follows")
	{
		followRouter.Use(middlewares.Authentication())
		followRouter.GET("/", controllers.GetFollow)
		followRouter.GET("/:followId", controllers.GetFollow)
		followRouter.POST("/", controllers.CreateFollow)
		followRouter.DELETE("/:followId", middlewares.FollowAuthorization(), controllers.DeleteFollow)
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
