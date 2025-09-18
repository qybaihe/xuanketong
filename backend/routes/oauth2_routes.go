package routes

import (
	"xuan-ke-tong/controllers"

	"github.com/gin-gonic/gin"
)

func OAuth2Routes(r *gin.Engine) {
	oauth2 := r.Group("/api/v1/auth/oauth2")
	{
		oauth2.GET("/state", controllers.GenerateOAuth2State)
		oauth2.GET("/callback", controllers.OAuth2Callback)
	}
}

