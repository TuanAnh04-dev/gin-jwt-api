package main

import (
	"go-jwt-api/config"
	"go-jwt-api/driver"
	"go-jwt-api/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	driver.ConnectMongoDB(config.DB_USER, config.DB_PASS)

	// userRepo := repoimpl.NewUserRepo(mongo.Client.Database(config.DB_NAME))

	// user := model.User{
	// 	Email:       "nguyenanh@gmail.com",
	// 	Password:    "123456",
	// 	DisplayName: "Nguyen Anh",
	// }

	// err := userRepo.Insert(user)
	// if err == nil {
	// 	fmt.Println("Chèn thành công!")
	// }

	r := gin.Default()
	routes.RegisterRoutes(r)
	r.Use(gin.Logger())

	r.Run(":8000")
	// fmt.Println("Server running [:8000]")
	// http.ListenAndServe(":8000", nil)
	// user, _ := userRepo.FindUserByEmail("nguyenanh@gmail.com")
	// fmt.Println(user)
}
