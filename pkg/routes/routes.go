package routes

import (
	"social-network-go/pkg/handlers"
	"social-network-go/pkg/views"

	"github.com/gin-gonic/gin"
)

func InitRoutes(r *gin.RouterGroup) {
	r.GET("/signup", views.SignupView)
	r.GET("/login", views.LoginView)
	r.GET("/home", views.HomeView)
	r.GET("/profile")
	r.POST("/signup", handlers.Signup)
	r.POST("/validate-email", handlers.ExistEmail)
	r.POST("/validate-username", handlers.ExistUsername)
	r.POST("/login", handlers.UserLogin)
	r.POST("/create-post", handlers.CreateNewPost)
	r.POST("/follow", handlers.Follow)
	r.POST("/unfollow", handlers.Unfollow)
	r.POST("/feed", handlers.Feed)
	r.POST("/loggout", handlers.Logout)
	r.GET("/:username", handlers.RenderProfileTemplate)
	r.POST("/profile", handlers.Profile)
	r.POST("/profile/:username", handlers.AnotherUserProfile)
	r.GET("/edit-profile", views.EditProfileView)
	r.POST("/edit-profile", handlers.EditProfile)
}



