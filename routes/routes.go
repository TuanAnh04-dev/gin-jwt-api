package routes

import (
	"go-jwt-api/handler"
	"go-jwt-api/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	// r.POST("/login", handler.Login)
	// r.POST("/register", handler.Register)
	// r.GET("/user", handler.GetUser)
	// r.GET("/product", handler.GetOneProduct)
	// r.GET("/listproduct", handler.GetListProduct)
	// r.POST("/addproduct", handler.AddNewProduct)

	r.POST("/login", handler.Login)
	r.POST("/register", handler.Register)
	r.GET("/user", middleware.AuthMiddleware(), handler.GetUser)
	r.GET("/product", middleware.AuthMiddleware(), handler.GetOneProduct)
	r.GET("/listproduct", middleware.AuthMiddleware(), handler.GetListProduct)
	r.POST("/addproduct", middleware.AuthMiddleware(), handler.AddNewProduct)
}
