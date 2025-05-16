package main

import (
	"fmt"
	"go-jwt-api/config"
	"go-jwt-api/driver"
	"go-jwt-api/handler"
	"go-jwt-api/middleware"
	"net/http"
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
	http.HandleFunc("/login", handler.Login)
	http.HandleFunc("/register", handler.Register)
	http.HandleFunc("/user", handler.GetUser)
	http.HandleFunc("/product", handler.GetOneProduct)
	http.HandleFunc("/listproduct", middleware.AuthMiddleware(handler.GetListProduct))
	http.HandleFunc("/addproduct", handler.AddNewProduct)

	fmt.Println("Server running [:8000]")
	http.ListenAndServe(":8000", nil)
	// user, _ := userRepo.FindUserByEmail("nguyenanh@gmail.com")
	// fmt.Println(user)
}
