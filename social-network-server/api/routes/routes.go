package routes

import (
	"social-network-server/api/handlers"
	middlewares "social-network-server/api/middleware"

	"github.com/gin-gonic/gin"
)

func InitRoutes(r *gin.RouterGroup) {
	r.POST("/signup", handlers.Signup)
	r.POST("/validate-email", handlers.ExistEmail)
	r.POST("/validate-username", handlers.ExistUsername)
	r.POST("/login", handlers.UserLogin)

	r.Use(middlewares.AuthMiddleware())

	r.POST("/create-post", handlers.CreateNewPost)
	r.POST("/follow", handlers.Follow)
	r.POST("/unfollow", handlers.Unfollow)
	r.POST("/feed", handlers.Feed)
	r.POST("/loggout", handlers.Logout)
	r.POST("/profile/:username", handlers.Profile)
	r.POST("/edit-profile", handlers.EditProfile)
	r.POST("create-message/:username", handlers.CreateNewMessage)
	r.POST("/chat/:username", handlers.Chat)
	r.GET("/websocket/:username", handlers.WebSocketHandler)
	r.POST("/chats", handlers.FeedChats)
	r.GET("/chats", handlers.WebSocketFeedChats)

}
